package http

import (
	"fmt"
	"io"

	"github.com/tomocy/wataame/http"
)

type Request struct {
	Addr http.Address
}

func (r *Request) Write(dst io.Writer) error {
	addr, err := r.Addr.Compensate()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(dst, "GET %s\n", addr)
	return err
}

type Response []byte
