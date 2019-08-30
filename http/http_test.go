package http

import "testing"

func TestAddress_Compensate(t *testing.T) {
	tests := map[string]struct {
		addr     Addr
		expected string
	}{
		"ipv4 with port": {
			"127.0.0.1:12345",
			"127.0.0.1:12345",
		},
		"ipv4 without port": {
			"127.0.0.1",
			"127.0.0.1:80",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := test.addr.Compensate()
			if err != nil {
				t.Fatalf("unexpected error from (Address).Compensate: got %s, expect nil\n", err)
			}
			if actual != test.expected {
				t.Errorf("unexpected address from (Address).Compensate: got %s, expect %s\n", actual, test.expected)
			}
		})
	}
}
