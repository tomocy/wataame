package client

import (
	"fmt"
	"io"
	"testing"

	http "github.com/tomocy/wataame/http/0.9"
	"golang.org/x/net/nettest"
)

func TestClient_Do(t *testing.T) {
	l, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatalf("failed to generate new local listener: %s\n", err)
	}
	defer l.Close()
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("failed for listener to accept: %s\n", err)
		}
		defer conn.Close()

		io.Copy(conn, conn)
	}()

	addr := l.Addr().String()
	expected := fmt.Sprintf("GET %s\n", addr)

	var client Client
	resp, err := client.Do(&http.Request{
		URI: addr,
	})
	if err != nil {
		t.Fatalf("failed for client to do: %s\n", err)
	}
	actual := string(resp)

	if actual != expected {
		t.Errorf("unexpected response: got %s, expect %s\n", actual, expected)
	}
}
