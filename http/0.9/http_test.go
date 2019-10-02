package http

import (
	"fmt"
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

func TestRequest_ReadFrom(t *testing.T) {
	uri, _ := url.Parse("http:///index.html")
	input := "GET /index.html"
	expected := &Request{
		Method: MethodGet, URI: uri,
	}

	actual := new(Request)
	if _, err := actual.ReadFrom(strings.NewReader(input)); err != nil {
		t.Errorf("unexpected error from (*Request).ReadFrom: got %s, expect nil\n", err)
		return
	}
	if err := assertRequest(actual, expected); err != nil {
		t.Errorf("unexpected (*Request).ReadFrom: %s\n", err)
	}
}

func assertRequest(actual, expected *Request) error {
	if actual.Method != expected.Method {
		return fmt.Errorf("unexpected method of request: got %s, expect %s", actual.Method, expected.Method)
	}
	if actual.URI.String() != expected.URI.String() {
		return fmt.Errorf("unexpected uri of request: got %s, expect %s", actual.URI, expected.URI)
	}

	return nil
}
