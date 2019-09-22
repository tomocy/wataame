package http

import (
	"io"
	"net/url"
)

type Request struct {
	Method string
	URI    *url.URL
	Body   io.ReadCloser
}
