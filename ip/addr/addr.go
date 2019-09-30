package addr

import (
	"net/url"
	"strings"
)

func ParseHostPort(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, "http://") {
		raw = "http://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	return parsed.Hostname(), parsed.Port(), nil
}
