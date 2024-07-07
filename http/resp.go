package http

import (
	"encoding/json"
	"net/http"

	"github.com/shadiestgoat/log"
	"whotfislucy.com/parser"
)

func writeErr(w http.ResponseWriter, status int, err string) {
	writeJson(w, status, &RespErr{Error: err})
}

func writeJson(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)
	log.ErrorIfErr(json.NewEncoder(w).Encode(body), "writing json body")
}

func nextPage(w http.ResponseWriter, s *parser.Section, ans string, sections *parser.SectionState, key string) {
	red := &parser.AnswerRedirect{}
	status := 200

	if s.Type == parser.ST_SLIDE {
		red = &parser.AnswerRedirect{
			Next:        s.Slide.Next,
			CorrectMode: parser.CM_UNKNOWN,
		}
	} else {
		next, ok := s.Question.Answers[ans]
		if ok {
			red = next
		} else {
			status = 422
			red.CorrectMode = parser.CM_BAD
		}
	}

	resp := &RespSection{
		TransitionType: red.CorrectMode,
	}

	if red.Next == parser.FINAL_SECTION_NAME {
		resp.Next = parser.GenerateFinale(key)
	} else if red.Next != "" {
		resp.Next = sections.SectionID[red.Next]
	}

	writeJson(w, status, resp)
}

type Req struct {
	CurrentSec string `json:"currentSection"`
	Answer     string `json:"answer"`
}

type RespErr struct {
	Error string `json:"error"`
}

type RespSection struct {
	Next           *parser.Section    `json:"next,omitempty"`
	TransitionType parser.CorrectMode `json:"transitionMode"`
}
