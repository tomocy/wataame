package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

func (r *FullRequest) WriteTo(dst io.Writer) (int64, error) {
	n, err := fmt.Fprint(dst, r)
	return int64(n), err
}

func (r *FullRequest) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, r.RequestLine)
	fmt.Fprintln(&b, r.Header)
	fmt.Fprintln(&b)
	if r.Body != nil {
		teed := r.teeBody()
		io.Copy(&b, teed)
	}

	return b.String()
}

func (r *FullRequest) teeBody() io.Reader {
	var b bytes.Buffer
	teed := io.TeeReader(r.Body, &b)
	r.Body = ioutil.NopCloser(&b)

	return teed
}

type RequestLine struct {
	Method  string
	URI     *url.URL
	Version *Version
}

func (l RequestLine) String() string {
	return fmt.Sprintf("%s %s %s", l.Method, l.URI.Path, l.Version)
}

type SimpleResponse struct {
	http0_9.Response
}

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

func (h Header) String() string {
	var b strings.Builder
	for k, vs := range h {
		for _, v := range vs {
			fmt.Fprintf(&b, "%s: %s\n", k, v)
		}
	}

	return strings.TrimSuffix(b.String(), "\n")
}

type headerField struct {
	key  string
	vals []string
}
