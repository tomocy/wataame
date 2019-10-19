package client

import (
	"context"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/http/1.0/server"
)

func TestClient_Do(t *testing.T) {
	testers := map[string]func(t *testing.T){
		// "simple request": testClientDoSimpleRequest,
		"full request": testClientDoFullRequest,
	}

	for name, tester := range testers {
		t.Run(name, tester)
	}
}

func testClientDoSimpleRequest(t *testing.T) {
	addr := ":1234"
	uri, _ := url.Parse("http://localhost" + addr + "/index.html")
	var client Client
	serv := &server.Server{
		Addr: addr,
		Handler: server.HandlerFunc{
			SimpleHandlerFunc: func(res *http1_0.SimpleResponse, req *http1_0.SimpleRequest) {
				req.WriteTo(res)
			},
		},
	}
	go serv.ListenAndServe()

	input := &http1_0.SimpleRequest{
		Request: http0_9.Request{
			Method: http1_0.MethodGet, URI: uri,
		},
	}
	expected := "GET /index.html\n"

	resp, err := client.Do(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	actual, ok := resp.(*http1_0.SimpleResponse)
	if !ok {
		t.Fatalf("unexpected type of Response from (*Clinet): got %T, expect *SimpleResponse\n", resp)
	}
	if string(actual.Response) != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %s, expect %s\n", string(actual.Response), expected)
	}
}

func testClientDoFullRequest(t *testing.T) {
	addr, path := ":1234", "/index.html"
	uri, _ := url.Parse("http://localhost" + addr + path)
	var client Client
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
	go serv.ListenAndServe()

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

	resp, err := client.Do(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	actual, ok := resp.(*http1_0.FullResponse)
	if !ok {
		t.Fatalf("unexpected type of Response from (*Clinet): got %T, expect *SimpleResponse\n", resp)
	}
	if actual.String() != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %q, expect %q\n", actual, expected)
	}
}
