package http

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
)

func TestFullRequest_Write(t *testing.T) {
	uri, _ := url.Parse("http://localhost:1234/index.html")
	tests := map[string]struct {
		req      *FullRequest
		expected string
	}{
		"GET method": {
			req: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodGet, URI: uri,
				},
				Header: Header{
					"Date": []string{"Tue, 15 Nov 1994 08:12:31 GMT"},
				},
			},
			expected: `GET /index.html HTTP/1.0
Date: Tue, 15 Nov 1994 08:12:31 GMT

`,
		},
		"HEAD method": {
			req: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodGet, URI: uri,
				},
				Header: Header{
					"User-Agent": []string{"CERN-LineMode/2.15", "libwww/2.17b3"},
				},
			},
			expected: `HEAD /index.html HTTP/1.0
User-Agent: CERN-LineMode/2.15 libwww/2.17b3

`,
		},
		"POST method": {
			req: &FullRequest{
				RequestLine: &RequestLine{
					Method: MethodGet, URI: uri,
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
			var actual bytes.Buffer
			if err := test.req.Write(&actual); err != nil {
				t.Errorf("unexpected error from (*Request).Write: got %s, expect nil\n", err)
				return
			}
			if actual.String() != test.expected {
				t.Errorf("unexpected request format from (*Request).Write: got %s, expect %s\n", &actual, test.expected)
			}
		})
	}
}
