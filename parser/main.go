package parser

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

type Section struct {
	Type     SectionType `json:"type"`
	Title    string      `json:"title"`
	ID       string      `json:"id"`

	Slide    *SlideOpts    `json:"slide,omitempty"`
	Question *QuestionOpts `json:"-"`
}

type QuestionOpts struct {
	Answers map[string]*AnswerRedirect `json:"-"`
}

type AnswerRedirect struct {
	Next        string
	CorrectMode CorrectMode
}

type SlideOpts struct {
	Next     string `json:"-"`
	SubTitle string `json:"subTitle,omitempty"`
	NextText string `json:"nextText"`
}

type SectionState struct {
	Sections  []*Section
	SectionID map[string]*Section
}

func Parse(inp []*RawSection) *SectionState {
	p := &parser{
		inp:                      inp,
		rawSectionSpecialAnswers: map[int]map[string]*AnswerRedirect{},
		takenIDs:                 map[string]bool{},
		prodIDToSec:              map[string]*Section{},
		prodSections:             []*Section{},
	}

	p.parse()

	return &SectionState{
		Sections:  p.prodSections,
		SectionID: p.prodIDToSec,
	}
}
