package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/http/1.0/client"
)

func TestServer_ListenAndServe(t *testing.T) {
	testers := map[string]func(t *testing.T){
		"simple request": testServerHandleSimpleRequest,
	}

	for name, tester := range testers {
		t.Run(name, tester)
	}
}

func testServerHandleSimpleRequest(t *testing.T) {
	addr := ":1234"
	var client client.Client
	serv := &Server{
		Addr: addr,
		SimpleHandler: SimpleHandlerFunc(func(resp *http1_0.SimpleResponse, req *http1_0.SimpleRequest) {
			if req.Method != http.MethodGet || req.URI.Path != "/index.html" {
				fmt.Fprint(resp, "not found")
				return
			}

			fmt.Fprint(resp, "hello world")
		}),
	}
	go serv.ListenAndServe()
	time.Sleep(1 * time.Second)

	tests := map[string]struct {
		input    http1_0.Request
		expected string
	}{
		"ok": {
			&http1_0.SimpleRequest{
				Request: http0_9.Request{
					Method: http0_9.MethodGet, URI: func() *url.URL {
						parsed, _ := url.Parse("http://" + filepath.Join(addr, "index.html"))
						return parsed
					}(),
				},
			},
			"hello world",
		},
		"not found": {
			&http1_0.SimpleRequest{
				Request: http0_9.Request{
					Method: http0_9.MethodGet, URI: func() *url.URL {
						parsed, _ := url.Parse("http://" + addr + "/")
						return parsed
					}(),
				},
			},
			"not found",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.Do(context.Background(), test.input)
			if err != nil {
				t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
			}
			actual, ok := resp.(*http1_0.SimpleResponse)
			if !ok {
				t.Fatalf("unexpected type of Response from (*Clinet): got %T, expect *SimpleResponse\n", resp)
			}
			if string(actual.Response) != test.expected {
				t.Errorf("unexpected Response from (*Client).Do: got %s, expect %s\n", string(actual.Response), test.expected)
			}
		})
	}
}
