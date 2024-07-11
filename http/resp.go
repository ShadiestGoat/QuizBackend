package http

import (
	"encoding/json"
	"net/http"

	"github.com/shadiestgoat/log"
)

type RespErr struct {
	Error string `json:"error"`
}

func writeJson(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)
	log.ErrorIfErr(json.NewEncoder(w).Encode(body), "writing json body")
}

type HttpFunc[Resp any] func () (*Resp, error)
type HttpBottomFunc func (w http.ResponseWriter, r *http.Request) (any, int)

func httpBotWrap(h HttpBottomFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, status := h(w, r)
		writeJson(w, status, resp)
	}
}

func httpBotErrWrap[Resp any](h func (w http.ResponseWriter, r *http.Request) (*Resp, error)) http.HandlerFunc {
	return httpBotWrap(func(w http.ResponseWriter, r *http.Request) (any, int) {
		resp, err := h(w, r)

		if err == nil {
			return resp, 200
		}

		var httpErr HttpErr
		if err, ok := err.(HttpErr); ok {
			httpErr = err
		} else {
			httpErr = ErrUnknown
		}

		return &RespErr{httpErr.Error()}, httpErr.Status()
	})
}

func httpWrap[Resp any](h HttpFunc[Resp]) http.HandlerFunc {
	return httpBotErrWrap(func(w http.ResponseWriter, r *http.Request) (*Resp, error) {
		return h()
	})
}

func httpBotWrapWithBody[Req any, Resp any](h func (w http.ResponseWriter, r *http.Request, b *Req) (*Resp, error)) http.HandlerFunc {
	return httpBotErrWrap(func(w http.ResponseWriter, r *http.Request) (*Resp, error) {
		var body Req
		
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, ErrBadSyntax
		}

		return h(w, r, &body)
	})
}

func httpWrapWithBody[Req any, Resp any](h func (b *Req) (*Resp, error)) http.HandlerFunc {
	return httpBotWrapWithBody(func(w http.ResponseWriter, r *http.Request, b *Req) (*Resp, error) {
		return h(b)
	})
}
