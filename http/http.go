package http

import (
	"os"
	"path/filepath"
	"strings"
)

const MethodGet = "GET"

type Addr string

func (a Addr) Compensate() (string, error) {
	splited := strings.Split(string(a), ":")
	if len(splited) < 2 {
		return splited[0] + ":80", nil
	}

	return strings.Join(splited[:2], ":"), nil
}

type ipv4Addr string

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}
