package http

import (
	"net/url"
	"testing"
)

func TestRequest_String(t *testing.T) {
	uri, _ := url.Parse("http://localhost:1234/index.html")
	subject := &Request{
		Method: MethodGet, URI: uri,
	}
	expected := "GET /index.html\n"
	actual := subject.String()

	if actual != expected {
		t.Errorf("unexpected result of (*Request).String: got %s, expect %s\n", actual, expected)
	}
}
