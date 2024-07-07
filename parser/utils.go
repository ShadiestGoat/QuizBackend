package parser

import (
	"strconv"

	"github.com/shadiestgoat/log"
)

func (p *parser) fatalSection(sectionI int, err string, args ...any) {
	log.Fatal("Failed to parsed section #%d: "+err, append([]any{sectionI}, args...)...)
}

func (p *parser) mkID(i int) string {
	base := strconv.Itoa(i)
	suffix := ""
	suffixI := 0

	for p.takenIDs[base+suffix] {
		suffix = "-" + strconv.Itoa(suffixI)
		suffixI++
	}

	p.takenIDs[base+suffix] = true

	return base + suffix
}

func (p *parser) validateSectionID(id string, sectionI int) {
	if id == FINAL_SECTION_NAME {
		return
	}

	if !p.takenIDs[id] {
		p.fatalSection(sectionI, "next id '%v' is not real", id)
	}
}

func (p *parser) checkStr(inp string, k string, sectionI int) {
	if inp == "" {
		p.fatalSection(sectionI, "string in key '%v' cannot be empty", k)
	}
}
