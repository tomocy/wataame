package http_test

import (
	"context"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	client0_9 "github.com/tomocy/wataame/http/0.9/client"
	http1_0 "github.com/tomocy/wataame/http/1.0"
	client1_0 "github.com/tomocy/wataame/http/1.0/client"
	"github.com/tomocy/wataame/http/1.0/server"
)

func TestServerClient_SimpleRequest(t *testing.T) {
	addr := ":1234"
	uri, _ := url.Parse("http://localhost" + addr + "/index.html")
	var client client0_9.Client
	serv := &server.Server{
		Addr: addr,
		Handler: server.HandlerFunc{
			SimpleHandlerFunc: func(res *http1_0.SimpleResponse, req *http1_0.SimpleRequest) {
				req.WriteTo(res)
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go serv.ListenAndServe(ctx)

	input := &http0_9.Request{
		Method: http1_0.MethodGet, URI: uri,
	}
	expected := "GET /index.html\n"

	actual, err := client.Do(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	if actual.String() != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %s, expect %s\n", actual, expected)
	}
}

func TestServerClient_FullRequest(t *testing.T) {
	addr, path := ":1234", "/index.html"
	uri, _ := url.Parse("http://localhost" + addr + path)
	var client client1_0.Client
	serv := &server.Server{
		Addr: addr,
		Handler: server.HandlerFunc{
			FullHandlerFunc: func(res *http1_0.FullResponse, req *http1_0.FullRequest) {
				if req.RequestLine.Method != http1_0.MethodPost || req.RequestLine.URI.Path != path {
					return
				}
				uas, ok := req.Header["User-Agent"]
				if !ok {
					return
				}
				if uas[0] != "wataame" {
					return
				}
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					return
				}
				if string(body) != "foo=bar" {
					return
				}

				if res.StatusLine == nil {
					res.StatusLine = new(http1_0.StatusLine)
				}
				res.StatusLine.Status = &http1_0.Status{
					Code: 201, Phrase: "Created",
				}
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go serv.ListenAndServe(ctx)

	input := &http1_0.FullRequest{
		RequestLine: &http1_0.RequestLine{
			Method: http1_0.MethodPost, URI: uri,
		},
		Header: http1_0.Header{
			"User-Agent": []string{"wataame"},
		},
		Body: ioutil.NopCloser(strings.NewReader("foo=bar")),
	}
	expected := "HTTP/1.0 201 Created\n\n"

	actual, err := client.Do(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	if actual.String() != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %q, expect %q\n", actual, expected)
	}
}
