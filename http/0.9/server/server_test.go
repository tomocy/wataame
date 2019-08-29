package server

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	http "github.com/tomocy/wataame/http"
)

func TestServer_ListenAndServe(t *testing.T) {
	addr := "localhost:12345"
	s := &Server{
		Addr: addr, Handler: &FileServer{
			Root: http.Dir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/tomocy/wataame/testdata")),
		},
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
			"GET",
			"http://" + filepath.Join(addr, "/index.html"),
			"<h1>Hello world</h1>",
		},
		"not found": {
			"GET",
			"http://" + filepath.Join(addr, "/"),
			"not found",
		},
		"method not allowed": {
			"HEAD",
			"http://" + filepath.Join(addr, "/index.html"),
			"method not allowed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := receiveTestResponse("tcp", test.method, test.uri)
			if err != nil {
				t.Fatalf("unexpected error from receiveTestResponse: got %s, expect nil\n", err)
			}
			if actual != test.expected {
				t.Errorf("unexpected response from receiveTestResponse: got %s, expect %s\n", actual, test.expected)
			}
		})
	}
}

func receiveTestResponse(network, method, uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	addr, err := http.Address(parsed.Host).Compensate()
	if err != nil {
		return "", err
	}
	conn, err := net.Dial(network, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s %s\n", method, filepath.Join(addr, parsed.Path))
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return "", err
	}

	return string(resp[:n]), nil
}
