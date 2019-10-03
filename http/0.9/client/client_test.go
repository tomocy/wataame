package client

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/http/0.9/server"
)

func TestClient_Do(t *testing.T) {
	addr := ":1234"
	serv := &server.Server{
		Addr: addr, Handler: server.HandlerFunc(func(w io.Writer, r *http0_9.Request) {
			r.WriteTo(w)
		}),
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil {
			t.Fatalf("unexpected error from (*Server).ListenAndServe: got %s, expect nil\n", err)
		}
	}()

	path := "/index.html"
	expected := fmt.Sprintf("%s %s\n", http0_9.MethodGet, path)

	var client Client
	uri, _ := url.Parse("http://localhost" + addr + path)
	resp, err := client.Do(context.Background(), &http0_9.Request{
		Method: http0_9.MethodGet, URI: uri,
	})
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: got %s, expect nil\n", err)
	}
	actual := string(resp)

	if actual != expected {
		t.Errorf("unexpected Response from (*Client).Do: got %s, expect %s\n", actual, expected)
	}
}
