package http

import (
	"fmt"
	"io"
	"net/url"
)

const MethodGet = "GET"

type Request struct {
	Method string
	URI    *url.URL
}

func (r *Request) WriteTo(dst io.Writer) (int64, error) {
	n, err := fmt.Fprintf(dst, "%s %s\n", r.Method, r.URI.EscapedPath())
	return int64(n), err
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.URI.Path)
}

type Response []byte
