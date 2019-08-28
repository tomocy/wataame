package http

import "testing"

func TestAddress_Compensate(t *testing.T) {
	tests := map[string]struct {
		addr     Address
		expected string
	}{
		"with port": {
			"localhost:12345",
			"localhost:12345",
		},
		"without port": {
			"localhost",
			"localhost:80",
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
