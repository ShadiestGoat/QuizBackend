package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"whotfislucy.com/parser"
)

func Router(info *parser.SectionState) chi.Router {
	r := chi.NewRouter()

	r.Post(`/`, func(w http.ResponseWriter, r *http.Request) {
		body := &Req{}
		
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, 400, "Failed to parse body: " + err.Error())
			return
		}

		if body.CurrentSec == "" {
			writeJson(w, 200, info.Sections[0])
			return
		}
		if body.CurrentSec == parser.FINAL_SECTION_NAME {
			writeErr(w, 400, "Finale has no next section")
			return
		}

		s, ok := info.SectionID[body.CurrentSec]
		if !ok {
			writeErr(w, 404, "Section unknown")
			return
		}

		nextPage(w, s, body.Answer, info, r.URL.Query().Get("key"))
	})

	return r
}
