package client

import (
	"fmt"
	"io"
	"net/url"
	"testing"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"golang.org/x/net/nettest"
)

func TestClient_Do(t *testing.T) {
	l, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatalf("unexpected error from nettest.NewLocalListener: got %s, expect nil\n", err)
	}
	defer l.Close()
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("unexpected error from (*Listener).Accept: %s\n", err)
		}
		defer conn.Close()

		io.Copy(conn, conn)
	}()

	addr := l.Addr().String()
	expected := fmt.Sprintf("GET %s\n", addr)

	var client Client
	uri, _ := url.Parse("http://" + addr)
	resp, err := client.Do(&http0_9.Request{
		URI: uri,
	})
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: %s\n", err)
	}
	actual := string(resp)

	if actual != expected {
		t.Errorf("unexpected response from (*Client).Do: got %s, expect %s\n", actual, expected)
	}
}
