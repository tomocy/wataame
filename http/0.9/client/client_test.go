package client

import (
	"context"
	"io"
	"net/url"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/http/0.9/server"
)

func TestClient_Do(t *testing.T) {
	addr := ":1234"
	uri, _ := url.Parse("http://localhost" + addr + "/index.html")
	var client Client
	serv := &server.Server{
		Addr: addr, Handler: server.HandlerFunc(func(w io.Writer, r *http0_9.Request) {
			r.WriteTo(w)
		}),
	}
	go serv.ListenAndServe()

	input := &http0_9.Request{
		Method: http0_9.MethodGet, URI: uri,
	}
	expected := "GET /index.html\n"

	actual, err := client.Do(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	if string(actual) != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %s, expect %s\n", string(actual), expected)
	}
}
