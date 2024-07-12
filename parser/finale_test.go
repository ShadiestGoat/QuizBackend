package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
	"whotfislucy.com/parser"
)

type FinaleTest struct {
	Essay string            `yaml:"essay"`
	FAQ   []*FinaleQuestion `yaml:"faq"`
}

type FinaleQuestion struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
}

const FINALE_TEST_DIR = "finale_tests"

func TestParseFinale(t *testing.T) {
	dir, err := os.ReadDir(FINALE_TEST_DIR)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]int{}

	for _, f := range dir {
		n := f.Name()

		if strings.HasSuffix(n, ".yaml") {
			tests[n[:len(n)-5]]++
		} else if strings.HasSuffix(n, ".md") {
			tests[n[:len(n)-3]]++
		}
	}

	for n, d := range tests {
		if d != 2 {
			continue
		}

		t.Run(n, func(t *testing.T) {
			raw_finale, err := os.ReadFile(filepath.Join(FINALE_TEST_DIR, n+".md"))
			if err != nil {
				t.Fatalf("Failed to read raw file: %v", err)
			}

			raw_answer, err := os.ReadFile(filepath.Join(FINALE_TEST_DIR, n+".yaml"))
			if err != nil {
				t.Fatalf("Failed to read answer file: %v", err)
			}

			expected := &FinaleTest{}
			if err := yaml.Unmarshal(raw_answer, expected); err != nil {
				t.Fatalf("Failed to parse answer: %v", err)
			}

			real := parser.ParseFinale(n, string(raw_finale))
			if real == nil {
				t.Fatal("Finale failed to parse")
			}

			t.Logf("----- Expected -----\n%v\n--- End Expected ---", expected.Essay)
			t.Logf("------- Real -------\n%v\n----- End Real -----", real.Essay)

			if expected.Essay != real.Essay {
				t.Errorf("Unequal essay")
			}

			if len(expected.FAQ) != len(real.FAQ) {
				t.Errorf("Bad FAQ Length: exp: %d; got: %d", len(expected.FAQ), len(real.FAQ))
			} else {
				for i := 0; i < len(expected.FAQ); i++ {
					e := expected.FAQ[i]
					r := real.FAQ[i]

					t.Logf("----- Expected -----\nQuestion: %v\nAnswer: %v\n--- End Expected ---", e.Question, e.Answer)
					t.Logf("------- Real -------\nQuestion: %v\nAnswer: %v\n----- End Real -----", r[0], r[1])

					if e.Question != r[0] {
						t.Error("Question mismatch")
					}
					if e.Answer != r[1] {
						t.Error("Answer mismatch")
					}
				}
			}
		})
	}
}
