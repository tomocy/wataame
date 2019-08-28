package http

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Address string

func (a Address) Compensate() (string, error) {
	host, port, err := a.split()
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func (a Address) split() (string, string, error) {
	splited := strings.Split(string(a), ":")
	if len(splited) > 2 {
		return "", "", errors.New("invalid format of uri: the format should be host[:port]")
	}
	var host, port string
	host = splited[0]
	if len(splited) == 2 {
		port = splited[1]
	}

	return host, port, nil
}

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}
