package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shadiestgoat/log"

	"gopkg.in/yaml.v3"
)

const FINAL_SECTION_NAME = "$COMPLETION"

type SectionType string

const (
	ST_SLIDE SectionType = "slide"
	ST_QUEST SectionType = "question"
)

type CorrectMode int

const (
	CM_GOOD CorrectMode = iota
	CM_BAD
	CM_UNKNOWN
)

type RawSection struct {
	Type  SectionType
	Title string
	ID    string `yaml:"id" json:"id"`
	Next  string

	Slide    RawSlideOpts
	Question RawQuestionOpts
}

type RawSlideOpts struct {
	SubTitle string
	NextText string
}

type RawQuestionOpts struct {
	Answers        []string
	SpecialAnswers map[string]*yaml.Node
}

type RawSpecialAnswer struct {
	Next        string
	CorrectMode CorrectMode
}

type Section struct {
	Type  SectionType
	Title string
	ID    string `json:"id"`

	Slide    *SlideOpts    `json:"slide,omitempty"`
	Question *QuestionOpts `json:"-"`
}

type QuestionOpts struct {
	Answers map[string]Redirect `json:"-"`
}

type Redirect struct {
	Next string
	CorrectMode CorrectMode
}

type SlideOpts struct {
	Next     string   `json:"-"`
	SubTitle string   `json:"subTitle,omitempty"`
	NextText string   `json:"nextText"`
}

func mkID(i int, taken map[string]bool) string {
	base := strconv.Itoa(i)
	suffix := ""
	suffixI := 0

	for taken[base+suffix] {
		suffix = "-" + strconv.Itoa(suffixI)
		suffixI++
	}

	taken[base+suffix] = true

	return base + suffix
}

func checkStr(inp string, k string, logBase string) {
	if inp == "" {
		log.Fatal("%v empty string '%v'", logBase, k)
	}
}

func parse(inp []*RawSection) {
	takenIDs := map[string]bool{}

	for i, s := range inp {
		if s.ID == "" {
			continue
		}
		logBase := fmt.Sprintf("Failed to parse section #%d:", i+1)

		if takenIDs[s.ID] {
			log.Fatal(logBase+" ID '%v' is taken", i+1, s.ID)
		}
		takenIDs[s.ID] = true

		checkStr(s.Title, "title", logBase)

		switch s.Type {
		case ST_SLIDE:
			checkStr(s.Slide.NextText, "slide.nextText", logBase)
		case ST_QUEST:
			if len(s.Question.Answers) == 0 {
				log.Fatal(logBase+" No answers provided", s.Type)
			}
		default:
			log.Fatal(logBase+" Unknown section type '%v'", s.Type)
		}

		checkStr(s.Title, "title", logBase)
	}

	rawSectionByID := map[string]*RawSection{}

	for i, s := range inp {
		if s.ID != "" {
			s.ID = mkID(i, takenIDs)
		}

		if i != 0 {
			if inp[i - 1].Next == "" {
				inp[i - 1].Next = s.ID
			}
		}

		rawSectionByID[s.ID] = s
	}

	if inp[len(inp) - 1].Next == "" {
		inp[len(inp) - 1].Next = FINAL_SECTION_NAME
	}

	cache := map[string]*Section{}
	o := []*Section{}

	for _, s := range inp {
		o = append(o, parseSection(s, cache, rawSectionByID))
	}
}

func validateSectionID(id string, allSections map[string]*RawSection) {
	if id == FINAL_SECTION_NAME {
		return
	}

	if _, ok := allSections[id]; !ok {
		log.Fatal("Failed to parse section: id '%v' isn't real", id)
	}
}

func parseSection(s *RawSection, cache map[string]*Section, rawSectionsByID map[string]*RawSection) *Section {
	if cache[s.ID] != nil {
		return cache[s.ID]
	}

	o := &Section{
		ID:       s.ID,
		Type:     s.Type,
		Title:    s.Title,
	}

	validateSectionID(s.Next, rawSectionsByID)

	if s.Type == ST_QUEST {
		o.Question = &QuestionOpts{
			Answers: map[string]Redirect{},
		}

		for _, ans := range s.Question.Answers {
			o.Question.Answers[strings.ToLower(ans)] = Redirect{
				Next:        s.Next,
				CorrectMode: CM_GOOD,
			}
		}

		for ans, n := range s.Question.SpecialAnswers {
			red := Redirect{
				CorrectMode: CM_GOOD,
			}

			if err := n.Decode(&red.Next); err != nil {
				log.FatalIfErr(n.Decode(&red), "trying to decode special answer '%v'", ans)
			}

			if red.CorrectMode < CM_GOOD || red.CorrectMode > CM_UNKNOWN {
				log.Fatal("Unknown correct mode '%v'", red.CorrectMode)
			}

			checkStr(red.Next, "specialAnswers." + ans + ".next", "Failed to decode special answers")
			validateSectionID(red.Next, rawSectionsByID)

			o.Question.Answers[strings.ToLower(ans)] = red
		}
	} else {
		o.Slide = &SlideOpts{
			Next:     s.Next,
			SubTitle: s.Slide.SubTitle,
			NextText: s.Slide.NextText,
		}
	}

	cache[s.ID] = o

	return o
}