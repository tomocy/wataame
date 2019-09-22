package http

import (
	"fmt"
	"io"
	"net/url"

	"github.com/tomocy/wataame/http"
)

const MethodGet = "GET"

type Request struct {
	Method string
	URI    *url.URL
}

func (r *Request) Write(dst io.Writer) error {
	addr, err := http.Addr(r.URI.Host).Compensate()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(dst, "%s %s\n", MethodGet, addr)
	return err
}

type Response []byte
