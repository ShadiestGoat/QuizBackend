package parser

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/shadiestgoat/log"
	"whotfislucy.com/encryption"
)

const DEFAULT_FINALE_NAME = "$default"
var regMultiline = regexp.MustCompile(`\n{2,}`)

var finaleCache = map[string]*FinaleCache{}
// A cache of already parsed aliases
var aliasCache = map[string]bool{}

func readFinaleFile(f string) string {
	d, err := os.ReadFile(f)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return ""
	}

	log.FatalIfErr(err, "reading file '%v'", f)

	return strings.TrimSpace(string(regMultiline.ReplaceAll(d, []byte{'\n'})))
}

// figure out the finale name & key for it.
func parseFinaleName(name string) (realName, key string) {
	if name == "" {
		// sanity check
		return "", ""
	}

	if name[0] != '$' {
		return "", name
	}

	if name == DEFAULT_FINALE_NAME {
		return name, name
	}

	spl := strings.SplitN(name, "-", 2)
	
	realName = strings.TrimSpace(strings.ToLower(spl[0][1:]))

	if len(spl) == 1 {
		key = os.Getenv("KEYS_" + realName)
	} else {
		key = spl[1]
	}

	return
}

// Prase a finale given the name of the folder & the full path including the folder
func ParseFinale(name string, path string) *FinaleCache {
	name, key := parseFinaleName(name)

	if name == "" {
		return nil
	}
	if aliasCache[name] {
		log.Fatal("Finale name '%v' is already taken!", name)
	}
	if finaleCache[key] != nil {
		log.Fatal("Finale key '%v' is duplicated!", key)
	}
	if name != DEFAULT_FINALE_NAME && key == DEFAULT_FINALE_NAME {
		log.Fatal("Key '%v' is reserved!", key)
	}

	cache := &FinaleCache{}
	
	aliasCache[name] = true
	finaleCache[key] = cache

	cache.Essay = readFinaleFile(filepath.Join(path, "essay.md"))

	rawFAQ := readFinaleFile(filepath.Join(path, "faq.md"))
	cache.FAQ = parseFAQ(rawFAQ)

	return cache
}

func parseFAQ(rawFAQ string) [][2]string {
	oFAQ := [][2]string{}
	curFAQ := [2]string{}

	addFAQ := func() {
		curFAQ[1] = strings.TrimSpace(curFAQ[1])
		oFAQ = append(oFAQ, curFAQ)
		
		curFAQ = [2]string{}
	}

	for _, l := range strings.Split(rawFAQ, "\n") {
		l = strings.TrimRightFunc(l, unicode.IsSpace)
		if l == "" {
			continue
		}

		if l := strings.TrimLeftFunc(l, unicode.IsSpace); strings.HasPrefix(l, "# ") {
			addFAQ()
			curFAQ[0] = l[2:]
			continue
		}

		// At this point, its def not a h1. But if it *is* a heading, lets upgrade it by 1
		if strings.HasPrefix(l, "#") {
			l = l[1:]
		}

		curFAQ[0] += l + "\n"
	}

	addFAQ()

	return oFAQ
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