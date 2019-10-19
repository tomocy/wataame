package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestFullRequest_WriteTo(t *testing.T) {
	uri, _ := url.Parse("http://localhost:1234/index.html")
	tests := map[string]struct {
		subject  *FullRequest
		expected string
	}{
		"GET method": {
			subject: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodGet, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"Date": []string{"Tue, 15 Nov 1994 08:12:31 GMT", "Wed, 16 Nov 1994 08:12:31 GMT"},
				},
			},
			expected: `GET /index.html HTTP/1.0
Date: Tue, 15 Nov 1994 08:12:31 GMT
Date: Wed, 16 Nov 1994 08:12:31 GMT
Content-Length: 0

`,
		},
		"HEAD method": {
			subject: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodHead, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"User-Agent": []string{"CERN-LineMode/2.15", "libwww/2.17b3"},
				},
			},
			expected: `HEAD /index.html HTTP/1.0
User-Agent: CERN-LineMode/2.15
User-Agent: libwww/2.17b3
Content-Length: 0

`,
		},
		"POST method": {
			subject: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodPost, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"Content-Type": []string{"application/x-www-form-urlencoded"},
				},
				Body: ioutil.NopCloser(strings.NewReader(url.Values{
					"name": []string{"foo"}, "password": []string{"bar"},
				}.Encode())),
			},
			expected: `POST /index.html HTTP/1.0
Content-Length: 21
Content-Type: application/x-www-form-urlencoded

name=foo&password=bar`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var b strings.Builder
			if _, err := test.subject.WriteTo(&b); err != nil {
				t.Fatalf("unexpected error from (*FullRequest).WriteTo: got %s, expect nil\n", err)
			}
			actual := b.String()
			if actual != test.expected {
				t.Errorf("unexpected FullRequest by (*FullRequest).WriteTo: got %s, expect %s\n", actual, test.expected)
			}
		})
	}
}

func TestFullRequest_ReadFrom(t *testing.T) {
	uri, _ := url.Parse("http:///index.html")
	tests := map[string]struct {
		input    string
		expected *FullRequest
	}{
		"GET method": {
			input: `GET /index.html HTTP/1.0
Date: Tue, 15 Nov 1994 08:12:31 GMT
Date: Wed, 16 Nov 1994 08:12:31 GMT

`,
			expected: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodGet, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"Date": []string{"Tue, 15 Nov 1994 08:12:31 GMT", "Wed, 16 Nov 1994 08:12:31 GMT"},
				},
			},
		},
		"HEAD method": {
			input: `HEAD /index.html HTTP/1.0
User-Agent: CERN-LineMode/2.15
User-Agent: libwww/2.17b3

`,
			expected: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodHead, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"User-Agent": []string{"CERN-LineMode/2.15", "libwww/2.17b3"},
				},
			},
		},
		"POST method": {
			input: `POST /index.html HTTP/1.0
Content-Length: 21
Content-Type: application/x-www-form-urlencoded

name=foo&password=bar`,
			expected: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodPost, URI: uri, Version: &Version{Major: 1, Minor: 0},
				},
				Header: Header{
					"Content-Length": []string{"21"},
					"Content-Type":   []string{"application/x-www-form-urlencoded"},
				},
				Body: ioutil.NopCloser(strings.NewReader(url.Values{
					"name": []string{"foo"}, "password": []string{"bar"},
				}.Encode())),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := new(FullRequest)
			if _, err := actual.ReadFrom(strings.NewReader(test.input)); err != nil {
				t.Fatalf("unexpected error from (*FullRequest).ReadFroom: got %s, expect nil\n", err)
			}
			if err := assertFullRequest(actual, test.expected); err != nil {
				t.Errorf("unexpected FullRequest by (*FullRequest).ReadFrom: %s\n", err)
			}
		})
	}
}

func assertFullRequest(actual, expected *FullRequest) error {
	if err := assertRequestLine(actual.RequestLine, expected.RequestLine); err != nil {
		return fmt.Errorf("unexpected request line of full request: %s", err)
	}
	if !reflect.DeepEqual(actual.Header, expected.Header) {
		return fmt.Errorf("unexpected header of full request: got %v, expect %v", actual.Header, expected.Header)
	}
	if err := assertBody(actual.Body, expected.Body); err != nil {
		return fmt.Errorf("unexpected body of full request: %s", err)
	}

	return nil
}

func assertRequestLine(actual, expected *RequestLine) error {
	if actual.Method != expected.Method {
		return fmt.Errorf("unexpected method of request line: got %s, expect %s", actual.Method, expected.Method)
	}
	if actual.URI.String() != expected.URI.String() {
		return fmt.Errorf("unexpected uri of request line: got %s, expect %s", actual.URI, expected.URI)
	}
	if err := assertVersion(actual.Version, expected.Version); err != nil {
		return fmt.Errorf("unexpected version of request line: %s", err)
	}

	return nil
}

func assertVersion(actual, expected *Version) error {
	if actual.Major != expected.Major {
		return fmt.Errorf("unexpected major of version: got %d, expect %d", actual.Major, expected.Major)
	}
	if actual.Minor != expected.Minor {
		return fmt.Errorf("unexpected minor of version: got %d, expect %d", actual.Minor, expected.Minor)
	}

	return nil
}

func assertBody(actual, expected io.ReadCloser) error {
	if actual == nil && expected == nil {
		return nil
	}

	defer func() {
		actual.Close()
		expected.Close()
	}()
	actualBody, _ := ioutil.ReadAll(actual)
	expectedBody, _ := ioutil.ReadAll(expected)
	if string(actualBody) != string(expectedBody) {
		return fmt.Errorf("unexpected body: got %s, expect %s", string(actualBody), string(expectedBody))
	}

	return nil
}
