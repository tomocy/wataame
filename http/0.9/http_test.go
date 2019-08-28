package http

import "testing"

func TestRequest_Address(t *testing.T) {
	tests := map[string]struct {
		r        *Request
		expected string
	}{
		"with port": {
			&Request{
				URI: "localhost:12345",
			},
			"localhost:12345",
		},
		"without port": {
			&Request{
				URI: "localhost",
			},
			"localhost:80",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := test.r.Address()
			if err != nil {
				t.Fatalf("unexpected error from (*Request).Address: got %s, expect nil\n", err)
			}
			if actual != test.expected {
				t.Errorf("unexpected address from (*Request).Address: got %s, expect %s\n", actual, test.expected)
			}
		})
	}
}
