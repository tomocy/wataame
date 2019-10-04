package client

import (
	"context"
	"io"
	"net/url"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/http/0.9/server"
	http1_0 "github.com/tomocy/wataame/http/1.0"
)

func TestClient_Do(t *testing.T) {
	testers := map[string]func(t *testing.T){
		"simple request": testClient_DoSimpleRequest,
	}

	for name, tester := range testers {
		t.Run(name, tester)
	}
}

func testClient_DoSimpleRequest(t *testing.T) {
	addr := ":1234"
	uri, _ := url.Parse("http://localhost" + addr + "/index.html")
	var client Client
	serv := &server.Server{
		Addr: addr,
		Handler: server.HandlerFunc(func(w io.Writer, r *http0_9.Request) {
			r.WriteTo(w)
		}),
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
