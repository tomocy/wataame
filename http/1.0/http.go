package http

import (
	"io"
	"net/url"

	http0_9 "github.com/tomocy/wataame/http/0.9"
)

const (
	MethodGet  = http0_9.MethodGet
	MethodHead = "HEAD"
	MethodPost = "POST"
)

type Request struct {
	Method string
	URI    *url.URL
	Body   io.ReadCloser
}

type Header map[string][]string
