package http

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MethodGet = "GET"

type Addr string

func (a Addr) Compensate() (string, error) {
	if strings.HasPrefix(string(a), "[") {
		return ipv6Addr(a).compensate()
	}

	return ipv4Addr(a).compensate()
}

func (a Addr) compensateWith(c interface {
	proto() string
	compensate() (string, error)
}) (string, error) {
	addr, err := c.compensate()
	if err != nil {
		return "", fmt.Errorf("failed to compensate address for %s: %s", c.proto(), err)
	}

	return addr, nil
}

type ipv6Addr string

func (a ipv6Addr) proto() string {
	return "IPv6"
}

func (a ipv6Addr) compensate() (string, error) {
	host, port, err := a.parse()
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "80"
	}

	return strings.Join([]string{host, port}, ":"), nil
}

func (a ipv6Addr) parse() (string, string, error) {
	if len(a) <= 0 {
		return "", "", errors.New("address is empty")
	}
	host, err := a.parseHost()
	if err != nil {
		return "", "", err
	}
	port, err := a.parsePort(host)
	if err != nil {
		return "", "", err
	}

	return host, port, nil
}

func (a ipv6Addr) parseHost() (string, error) {
	if a[0] != '[' {
		return "", errors.New("[ is missing")
	}
	end := strings.IndexByte(string(a)[1:], ']') + 1
	if end < 0 {
		return "", errors.New("] is missing")
	}

	return string(a)[0 : end+1], nil
}

func (a ipv6Addr) parsePort(host string) (string, error) {
	splited := strings.Split(string(a), host)
	if len(splited) < 2 {
		return "", nil
	}
	port := splited[1]
	if 1 <= len(port) && port[0] == ':' {
		port = port[1:]
	}

	return port, nil
}

type ipv4Addr string

func (a ipv4Addr) proto() string {
	return "IPv4"
}

func (a ipv4Addr) compensate() (string, error) {
	splited := strings.Split(string(a), ":")
	if len(splited) < 2 {
		return splited[0] + ":80", nil
	}

	return strings.Join(splited[:2], ":"), nil
}

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}
