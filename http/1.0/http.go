package http

import (
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

const (
	MethodGet  = http0_9.MethodGet
	MethodHead = "HEAD"
	MethodPost = "POST"
)

type SimpleRequest http0_9.Request

type Response struct {
	StatusLine *StatusLine
	Header     Header
	Body       []byte
}

type StatusLine struct {
	Version string
	Status  *Status
}

type Status struct {
	Code   int
	Phrase string
}

type Header map[string][]string
