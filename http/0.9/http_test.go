package http

import (
	"net/url"
	"strings"
	"testing"
)

func TestRequest_WriteTo(t *testing.T) {
	uri, _ := url.Parse("http://golang.org/index.html")
	subject := &Request{
		Method: MethodGet, URI: uri,
	}
	expected := "GET /index.html\n"
	var actual strings.Builder
	subject.WriteTo(&actual)

	if actual.String() != expected {
		t.Errorf("unexpected result of (*Request).WriteTo: got %s, expect %s\n", &actual, expected)
	}
}
