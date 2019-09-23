package client

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"path/filepath"
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

		r := bufio.NewReader(conn)
		read, _, err := r.ReadLine()
		if err != nil {
			t.Fatalf("unexpected error from (*Reader).Read: %s\n", err)
		}

		fmt.Fprintln(conn, string(read))
	}()

	path := "/index.html"
	expected := fmt.Sprintf("%s %s\n", http0_9.MethodGet, path)

	var client Client
	uri, _ := url.Parse("http://" + filepath.Join(l.Addr().String(), path))
	resp, err := client.Do(context.Background(), &http0_9.Request{
		Method: http0_9.MethodGet, URI: uri,
	})
	if err != nil {
		t.Fatalf("unexpected error from (*Client).Do: %s\n", err)
	}
	actual := string(resp)

	if actual != expected {
		t.Errorf("unexpected response from (*Client).Do: got %s, expect %s\n", actual, expected)
	}
}
