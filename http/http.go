package http

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}

type ScannableURL url.URL

func (u *ScannableURL) Scan(state fmt.ScanState, _ rune) error {
	var raw string
	if _, err := fmt.Fscan(state, &raw); err != nil {
		return fmt.Errorf("failed to scan uri: %s", err)
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("failed to scan uri: %s", err)
	}

	*u = ScannableURL(*parsed)

	return nil
}
