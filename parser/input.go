package parser

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type RawSlideOpts struct {
	SubTitle string `yaml:"subTitle" json:"subTitle"`
	NextText string `yaml:"nextText" json:"nextText"`
}

type RawQuestionOpts struct {
	Answers []string
	// String or RawSpecialAnswer
	SpecialAnswers map[string]yaml.Node `yaml:"specialAnswers"`
}

type RawSection struct {
	Type  SectionType
	Title string
	ID    string `yaml:"id" json:"id"`
	Next  string

	Slide    RawSlideOpts
	Question RawQuestionOpts
}

func (p *parser) checkAndCacheInp() {
	for i, s := range p.inp {
		if s.ID != "" {
			if s.ID == FINAL_SECTION_NAME {
				p.fatalSection(i, "ID '%v' is reserved", s.ID)
			}
			if p.takenIDs[s.ID] {
				p.fatalSection(i, "ID '%v' is taken", s.ID)
			}

			p.takenIDs[s.ID] = true
		}

		p.checkStr(s.Title, "title", i)

		switch s.Type {
		case ST_SLIDE:
			p.checkStr(s.Slide.NextText, "slide.nextText", i)
		case ST_QUEST:
			if len(s.Question.Answers)+len(s.Question.SpecialAnswers) == 0 {
				p.fatalSection(i, "No answers provided")
			}
			specialMap := map[string]*AnswerRedirect{}

			if len(s.Question.SpecialAnswers) != 0 {
				for ans, node := range s.Question.SpecialAnswers {
					fin := &AnswerRedirect{
						CorrectMode: CM_GOOD,
					}

					if err := node.Decode(&fin.Next); err != nil {
						if err := node.Decode(&fin); err != nil {
							p.fatalSection(i, "unable to decode special answer '%v': %v", ans, err)
						}

						if fin.CorrectMode < CM_GOOD || fin.CorrectMode > CM_UNKNOWN {
							p.fatalSection(i, "unknown correct mode '%v'", fin.CorrectMode)
						}
					}

					specialMap[strings.ToLower(ans)] = fin
				}
			}

			p.rawSectionSpecialAnswers[i] = specialMap
		default:
			p.fatalSection(i, "Unknown section type '%v'", s.Type)
		}
	}

	for i, s := range p.inp {
		if s.Next != "" {
			p.validateSectionID(s.Next, i)
		}
		if s.Type == ST_QUEST {
			for _, resp := range p.rawSectionSpecialAnswers[i] {
				p.validateSectionID(resp.Next, i)
			}
		}
	}
}
