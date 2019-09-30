package http

import (
	"net/url"
	"strings"
	"testing"
)

func TestRequest_WriteTo(t *testing.T) {
	uri, _ := url.Parse("http://localhost:1234/index.html")
	subject := &Request{
		Method: MethodGet, URI: uri,
	}
	expected := "GET /index.html\n"
	var b strings.Builder
	subject.WriteTo(&b)
	actual := b.String()

	if actual != expected {
		t.Errorf("unexpected (*Request).WriteTo: got %s, expect %s\n", actual, expected)
	}
}
