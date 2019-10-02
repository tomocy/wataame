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
Content-Type: application/x-www-form-urlencoded

name=foo&password=bar`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var b strings.Builder
			test.subject.WriteTo(&b)
			actual := b.String()
			if actual != test.expected {
				t.Errorf("unexpected (*FullRequest).WriteTo: got %s, expect %s\n", actual, test.expected)
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
