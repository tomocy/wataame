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

func (r *Request) Write(dst io.Writer) error {
	var err error
	_, err = fmt.Fprintf(dst, "%s %s\n", r.Method, r.URI.EscapedPath())

	return err
}

type Response []byte
