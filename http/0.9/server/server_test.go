package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"path/filepath"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/ip"
)

func TestServer_ListenAndServe(t *testing.T) {
	addr := "localhost:12345"
	s := &Server{
		Addr: addr, Handler: HandlerFunc(func(w io.Writer, r *http0_9.Request) {
			if r.URI.Path != "/index.html" {
				fmt.Fprint(w, "not found")
				return
			}

			fmt.Fprint(w, "<h1>Hello world</h1>")
		}),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			t.Fatalf("unexpected error from (*Server).ListenAndServe: got %s, expect nil\n", err)
		}
	}()

	tests := map[string]struct {
		method, uri string
		expected    string
	}{
		"ok": {
			http0_9.MethodGet,
			"http://" + filepath.Join(addr, "/index.html"),
			"<h1>Hello world</h1>\n",
		},
		"not found": {
			http0_9.MethodGet,
			"http://" + filepath.Join(addr, "/"),
			"not found\n",
		},
		"method not allowed": {
			"HEAD",
			"http://" + filepath.Join(addr, "/index.html"),
			"method not allowed\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := receiveTestResponse("tcp", test.method, test.uri)
			if err != nil {
				t.Fatalf("unexpected error from (*Server).Serve: got %s, expect nil\n", err)
			}
			if actual != test.expected {
				t.Errorf("unexpected Response from (*Server).Serve: got %q, expect %q\n", actual, test.expected)
			}
		})
	}
}

func receiveTestResponse(network, method, uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	addr, err := ip.Addr(parsed.Host).Compensate()
	if err != nil {
		return "", err
	}
	conn, err := net.Dial(network, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s %s\n", method, filepath.Join(addr, parsed.Path))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
