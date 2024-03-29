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
	if _, err := subject.WriteTo(&b); err != nil {
		t.Fatalf("unexpected error from (*Request).WriteTo: got %s, expect nil\n", err)
	}
	actual := b.String()
	if actual != expected {
		t.Errorf("unexpected Request by (*Request).WriteTo: got %s, expect %s\n", actual, expected)
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
		t.Fatalf("unexpected error from (*Request).ReadFrom: got %s, expect nil\n", err)
		return
	}
	if err := assertRequest(actual, expected); err != nil {
		t.Errorf("unexpected Request by (*Request).ReadFrom: %s\n", err)
	}
}

func assertRequest(actual, expected *Request) error {
	if actual.Method != expected.Method {
		return fmt.Errorf("unexpected Method of Request: got %s, expect %s", actual.Method, expected.Method)
	}
	if actual.URI.String() != expected.URI.String() {
		return fmt.Errorf("unexpected URI of Request: got %s, expect %s", actual.URI, expected.URI)
	}

	return nil
}

func TestResponse_ReadFrom(t *testing.T) {
	input := "hello world"
	expected := input

	var actual Response
	if _, err := actual.ReadFrom(strings.NewReader(input)); err != nil {
		t.Fatalf("unexpected error from (*Response).ReadFrom: got %s, expect nil\n", err)
		return
	}
	if string(actual) != expected {
		t.Errorf("unexpected Response by (*Response).ReadFrom: got %s, expect %s\n", string(actual), expected)
	}
}
