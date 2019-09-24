package http

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	http0_9 "github.com/tomocy/wataame/http/0.9"
)

const (
	MethodGet  = http0_9.MethodGet
	MethodHead = "HEAD"
	MethodPost = "POST"
)

type SimpleRequest struct {
	http0_9.Request
}

type FullRequest struct {
	RequestLine *RequestLine
	Header      Header
	Body        io.ReadCloser
}

type RequestLine struct {
	Method  string
	URI     *url.URL
	Version *Version
}

func (l RequestLine) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s %s", l.Method, l.URI.Path, l.Version)

	return b.String()
}

type SimpleResponse http0_9.Response

type FullResponse struct {
	StatusLine *StatusLine
	Header     Header
	Body       io.ReadCloser
}

type StatusLine struct {
	Version *Version
	Status  *Status
}

type Version struct {
	Major, Minor int
}

func (v Version) String() string {
	return fmt.Sprintf("HTTP/%d.%d", v.Major, v.Minor)
}

type Status struct {
	Code   int
	Phrase string
}

type Header map[string][]string
