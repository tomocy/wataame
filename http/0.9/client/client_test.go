package client

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"testing"

	http "github.com/tomocy/wataame/http"
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

		r := bufio.NewReader(conn)
		read, _, err := r.ReadLine()
		if err != nil {
			t.Fatalf("unexpected error from (*Reader).Read: %s\n", err)
		}

		fmt.Fprintln(conn, string(read))
	}()

	addr := l.Addr().String()
	expected := fmt.Sprintf("%s %s\n", http.MethodGet, addr)

	var client Client
	uri, _ := url.Parse("http://" + addr)
	resp, err := client.Do(context.Background(), &http0_9.Request{
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
