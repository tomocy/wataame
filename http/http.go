package http

import (
	"os"
	"path/filepath"
)

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}
