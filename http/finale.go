package http

import "whotfislucy.com/parser"

type ReqFinale struct {
	Secret string `json:"finaleSecret"`
}

type RespFinale struct {
	parser.FinaleCache
}

func postFinale(b *ReqFinale) (*RespFinale, error) {
	resp := parser.GetFinale(b.Secret)

	if resp == nil {
		return nil, ErrBadSecret
	}

	return &RespFinale{*resp}, nil
}
