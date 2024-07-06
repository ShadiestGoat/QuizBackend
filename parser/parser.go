package parser

import "strings"

type parser struct {
	inp []*RawSection
	rawSectionSpecialAnswers map[int]map[string]*AnswerRedirect
	takenIDs map[string]bool

	prodIDToSec map[string]*Section
	prodSections []*Section
}


func (p *parser) prepInp() {
	for i, s := range p.inp {
		if s.ID != "" {
			s.ID = p.mkID(i)
		}

		if i != 0 {
			if p.inp[i - 1].Next == "" {
				p.inp[i - 1].Next = s.ID
			}
		}
	}

	if p.inp[len(p.inp) - 1].Next == "" {
		p.inp[len(p.inp) - 1].Next = FINAL_SECTION_NAME
	}
}

func (p *parser) parseSection(i int) *Section {
	s := p.inp[i]

	o := &Section{
		ID:       s.ID,
		Type:     s.Type,
		Title:    s.Title,
	}

	if s.Type == ST_QUEST {
		o.Question = &QuestionOpts{
			Answers: map[string]*AnswerRedirect{},
		}

		for _, ans := range s.Question.Answers {
			o.Question.Answers[strings.ToLower(ans)] = &AnswerRedirect{
				Next:        s.Next,
				CorrectMode: CM_GOOD,
			}
		}

		for ans, red := range p.rawSectionSpecialAnswers[i] {
			o.Question.Answers[ans] = red
		}
	} else {
		o.Slide = &SlideOpts{
			Next:     s.Next,
			SubTitle: s.Slide.SubTitle,
			NextText: s.Slide.NextText,
		}
	}

	return o
}


func (p *parser) parse() {
	p.checkAndCacheInp()
	p.prepInp()

	for i := range p.inp {
		s := p.parseSection(i)

		p.prodIDToSec[s.ID] = s
		p.prodSections = append(p.prodSections, s)
	}
}
