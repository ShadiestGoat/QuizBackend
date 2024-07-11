package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"whotfislucy.com/parser"
)

func Router(info *parser.SectionState) chi.Router {
	r := chi.NewRouter()

	r.Use(
		cors.AllowAll().Handler,
		middleware.CleanPath,
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json"),
	)

	r.Post(`/`, httpWrapWithBody(func(b *ReqNextSec) (*RespNextSec, error) {
		return postNextSec(b, info)
	}))

	r.Post(`/finale`, httpWrapWithBody(postFinale))

	return r
}
