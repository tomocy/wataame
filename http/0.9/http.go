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

func (r *Request) WriteTo(dst io.Writer) (int64, error) {
	n, err := fmt.Fprintln(dst, r)
	return int64(n), err
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.URI.Path)
}

func (r *Request) Scan(state fmt.ScanState, _ rune) error {
	var uri http.ScannableURL
	if _, err := fmt.Fscanf(state, "%s %v", &r.Method, &uri); err != nil {
		return fmt.Errorf("failed to scan request: %s", err)
	}
	casted := url.URL(uri)
	r.URI = &casted

	return nil
}

type Response []byte
