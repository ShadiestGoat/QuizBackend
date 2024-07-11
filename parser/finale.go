package parser

import (
	"regexp"
	"strings"

	"github.com/shadiestgoat/log"
	"whotfislucy.com/encryption"
)

const DEFAULT_FINALE_NAME = "default_finale"
var regMultiline = regexp.MustCompile(`\n{2,}`)

var finaleCache = map[string]*FinaleCache{}

func ParseFinale(key string, fileContent string) {
	fileContent = regMultiline.ReplaceAllString(fileContent, "\n")

	faq := [][2]string{}
	curFaq := [2]string{}

	essay := ""

	lastH := ""
	headingContentStarted := false

	getStrPtr := func () *string {
		if lastH == "faq" {
			return &curFaq[1]
		} else if lastH == "essay" {
			return &essay
		}

		return nil
	}

	headingCB := func() {
		lastPtr := getStrPtr()

		if lastPtr != nil {
			*lastPtr = strings.TrimSpace(*lastPtr)
		}

		if lastH == "faq" && curFaq[0] != "" && headingContentStarted {
			faq = append(faq, [2]string{curFaq[0], strings.TrimSpace(curFaq[1])})
		}

		headingContentStarted = false
		curFaq = [2]string{}
	}

	for _, l := range strings.Split(fileContent, "\n") {
		if strings.HasPrefix(l, "# ") {
			headingCB()

			lastH = strings.ToLower(strings.TrimSpace(l[2:]))
			continue
		}
		if lastH == "faq" && strings.HasPrefix(l, "## ") {
			headingCB()

			curFaq[0] = strings.ToLower(strings.TrimSpace(l[3:]))
			continue
		}

		if l == "" {
			if !headingContentStarted {
				continue
			}
		} else {
			headingContentStarted = true
		}

		if lastH == "faq" && curFaq[0] == "" {
			log.Warn("Result file '%v' has unknown content under a FAQ header", key)
			continue
		}

		ptr := getStrPtr()

		if ptr != nil {
			*ptr += l + "\n"
		}
	}

	headingCB()

	if len(faq) == 0 && len(essay) == 0 {
		return
	}

	finaleCache[key] = &FinaleCache{
		FAQ:   faq,
		Essay: strings.TrimSpace(essay),
	}
}

func GetFinale(secret string) *FinaleCache {
	dec := encryption.Decrypt(secret)
	if dec == "" {
		return nil
	}

	return finaleCache[dec]
}

func GenerateFinale(key string) *Section {
	f := finaleCache[key]

	if f == nil {
		f = finaleCache[DEFAULT_FINALE_NAME]
		key = DEFAULT_FINALE_NAME

		if f == nil {
			log.Error("No default section name!")
			return nil
		}
	}

	return &Section{
		Type:     ST_FINAL,
		ID:       FINAL_SECTION_NAME,
		Finale:   &FinaleOpts{
			Secret: encryption.Encrypt(key),
		},
	}
}