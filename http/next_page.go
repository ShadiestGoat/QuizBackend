package http

import (
	"strings"

	"whotfislucy.com/parser"
)

type ReqNextSec struct {
	CurrentSec string `json:"currentSection"`
	Answer     string `json:"answer"`
	Key        string `json:"key"`
}

type RespNextSec struct {
	Next           *parser.Section    `json:"next,omitempty"`
	TransitionType parser.CorrectMode `json:"transitionMode"`
}

func postNextSec(b *ReqNextSec, info *parser.SectionState) (*RespNextSec, error) {
	if b.CurrentSec == "" {
		return &RespNextSec{info.Sections[0], parser.CM_UNKNOWN}, nil
	}

	if b.CurrentSec == parser.FINAL_SECTION_NAME {
		return nil, ErrNoNext
	}

	s, ok := info.SectionID[b.CurrentSec]
	if !ok {
		return nil, ErrSectionUnknown
	}

	red := &parser.AnswerRedirect{}

	if s.Type == parser.ST_SLIDE {
		red = &parser.AnswerRedirect{
			Next:        s.Slide.Next,
			CorrectMode: parser.CM_UNKNOWN,
		}
	} else {
		next, ok := s.Question.Answers[strings.ToLower(strings.TrimSpace(b.CurrentSec))]
		if ok {
			red = next
		} else {
			red.CorrectMode = parser.CM_BAD
		}
	}

	resp := &RespNextSec{
		TransitionType: red.CorrectMode,
	}

	if red.Next == parser.FINAL_SECTION_NAME {
		finale := parser.GenerateFinale(b.Key)
		if finale == nil {
			return nil, ErrBadSecret
		}

		resp.Next = finale
	} else if red.Next != "" {
		resp.Next = info.SectionID[red.Next]
	}

	return resp, nil
}
