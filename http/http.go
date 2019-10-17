package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
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

func DetectVersion(conn *PeekableConn) (string, error) {
	line, err := conn.PeekLine()
	if err != nil {
		return "", fmt.Errorf("failed to detect version: %s", err)
	}

	idx := strings.LastIndex(line, "HTTP/")
	if idx < 0 {
		return "0.9", nil
	}
	if len(line) <= idx+5 {
		return "", fmt.Errorf("failed to detect version: invalid format of HTTP version: got %s, expect HTTP/majour.minor", line[idx:])
	}

	return strings.TrimRight(line[idx+5:], "\n"), nil
}

type PeekableConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *PeekableConn) PeekLine() (string, error) {
	c.r = bufio.NewReader(c.Conn)
	var (
		peeked []byte
		err    error
	)
	for i := 1; i < c.r.Size(); i++ {
		peeked, err = c.r.Peek(i)
		if err != nil || bytes.HasSuffix(peeked, []byte{'\n'}) {
			break
		}
	}

	return string(peeked), err
}

func (c *PeekableConn) Read(dst []byte) (int, error) {
	var r io.Reader = c.Conn
	if c.r != nil {
		r = c.r
	}

	return r.Read(dst)
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

func (u *ScannableURL) URL() *url.URL {
	casted := url.URL(*u)
	return &casted
}
