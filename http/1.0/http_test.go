package http

import (
	"io/ioutil"
	"net/url"
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
					"Date": []string{"Tue, 15 Nov 1994 08:12:31 GMT"},
				},
			},
			expected: `GET /index.html HTTP/1.0
Date: Tue, 15 Nov 1994 08:12:31 GMT

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
User-Agent: CERN-LineMode/2.15 libwww/2.17b3

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
