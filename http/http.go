package http

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

type Address string

func (a Address) Compensate() (string, error) {
	host, port, err := net.SplitHostPort(string(a))
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}
