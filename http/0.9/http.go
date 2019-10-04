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

func (r *Request) ReadFrom(src io.Reader) (int64, error) {
	n, err := fmt.Fscan(src, r)
	return int64(n), err
}

func (r *Request) Scan(state fmt.ScanState, _ rune) error {
	var uri http.ScannableURL
	if _, err := fmt.Fscanf(state, "%s %v", &r.Method, &uri); err != nil {
		return fmt.Errorf("failed to scan request: %s", err)
	}
	r.URI = uri.URL()

	return nil
}

type Response []byte

func (r *Response) Write(src []byte) (int, error) {
	*r = append(*r, src...)
	return len(src), nil
}

func (r *Response) ReadFrom(src io.Reader) (int64, error) {
	n, err := fmt.Fscan(src, r)
	return int64(n), err
}

func (r *Response) Scan(state fmt.ScanState, _ rune) error {
	var reads []rune
	for {
		read, _, err := state.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to scan response: %s", err)
		}

		reads = append(reads, read)
	}

	*r = Response(string(reads))

	return nil
}
