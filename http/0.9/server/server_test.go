package server

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/http/0.9/client"
)

func TestServer_ListenAndServe(t *testing.T) {
	addr := "localhost:1234"
	var client client.Client
	s := &Server{
		Addr: addr, Handler: HandlerFunc(func(w io.Writer, r *http0_9.Request) {
			if r.URI.Path != "/index.html" {
				fmt.Fprint(w, "not found")
				return
			}

			fmt.Fprint(w, "<h1>Hello world</h1>")
		}),
	}
	go s.ListenAndServe()

	tests := map[string]struct {
		input    *http0_9.Request
		expected string
	}{
		"ok": {
			&http0_9.Request{
				Method: http0_9.MethodGet, URI: func() *url.URL {
					parsed, _ := url.Parse("http://" + filepath.Join(addr, "index.html"))
					return parsed
				}(),
			},
			"<h1>Hello world</h1>",
		},
		"not found": {
			&http0_9.Request{
				Method: http0_9.MethodGet, URI: func() *url.URL {
					parsed, _ := url.Parse("http://" + addr + "/")
					return parsed
				}(),
			},
			"not found",
		},
		"method not allowed": {
			&http0_9.Request{
				Method: "HEAD", URI: func() *url.URL {
					parsed, _ := url.Parse("http://" + filepath.Join(addr, "index.html"))
					return parsed
				}(),
			},
			"method not allowed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := client.Do(context.Background(), test.input)
			if err != nil {
				t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
			}
			if string(actual) != test.expected {
				t.Errorf("unexpected Response from (*Client).Do: got %q, expect %q\n", string(actual), test.expected)
			}
		})
	}
}
