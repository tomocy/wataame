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
	n, err := fmt.Fprint(dst, r)
	return int64(n), err
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s\n", r.Method, r.URI.Path)
}

type Response []byte
