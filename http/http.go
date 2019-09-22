package http

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Addr string

func (a Addr) Compensate() (string, error) {
	if strings.HasPrefix(string(a), "[") {
		return a.compensateWith(ipv6Addr(a))
	}

	return a.compensateWith(ipv4Addr(a))
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
	host, port, err := parseHostPort(string(a))
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("[%s]:%s", host, port), nil
}

type ipv4Addr string

func (a ipv4Addr) proto() string {
	return "IPv4"
}

func (a ipv4Addr) compensate() (string, error) {
	host, port, err := parseHostPort(string(a))
	if err != nil {
		return "", nil
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func parseHostPort(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, "http://") {
		raw = "http://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	return parsed.Hostname(), parsed.Port(), nil
}

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	fname := filepath.Join(string(d), name)
	return os.Open(fname)
}

type FileSystem interface {
	Open(string) (*os.File, error)
}
