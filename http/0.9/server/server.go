package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/ip"
)

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) ListenAndServe() error {
	l, err := s.listen()
	if err != nil {
		return fmt.Errorf("failed to listen and serve: %s", err)
	}
	defer l.Close()

	if err := s.Serve(l); err != nil {
		return fmt.Errorf("failed to listen and serve: %s", err)
	}

	return nil
}

func (s *Server) listen() (net.Listener, error) {
	compensated, err := ip.Addr(s.Addr).Compensate()
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %s", err)
	}

	l, err := net.Listen("tcp", compensated)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %s", err)
	}

	return l, nil
}

func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to serve: %s", err)
		}
		go func() {
			defer conn.Close()
			r, err := readRequest(conn)
			if err != nil {
				fmt.Fprintf(conn, "failed to serve: %s\n", err)
				return
			}

			if r.Method != http0_9.MethodGet {
				fmt.Fprintln(conn, "failed to serve: method not allowed")
				return
			}

			s.Handler.Handle(conn, r)
			fmt.Fprintln(conn)
		}()
	}
}

func readRequest(conn net.Conn) (*http0_9.Request, error) {
	var b bytes.Buffer
	r := bufio.NewReader(conn)
	for line, isPrefix, err := r.ReadLine(); ; {
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %s", err)
		}
		b.Write(line)
		if !isPrefix {
			break
		}
	}

	req, err := parseRequest(b.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %s", err)
	}

	return req, nil
}

func parseRequest(bs []byte) (*http0_9.Request, error) {
	splited := strings.Split(string(bs), " ")
	if len(splited) < 2 {
		return nil, fmt.Errorf("failed to parse request: invalid format of request: got %s, expect method uri", string(bs))
	}
	uri, err := url.Parse("http://" + splited[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %s", err)
	}

	return &http0_9.Request{
		Method: splited[0],
		URI:    uri,
	}, nil
}

type FileServer struct {
	Root http.FileSystem
}

func (s *FileServer) Handle(w io.Writer, r *http0_9.Request) {
	f, err := s.Root.Open(r.URI.Path)
	if err != nil {
		fmt.Fprint(w, "not found")
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Fprint(w, "internal server error")
		return
	}
	if stat.IsDir() {
		fmt.Fprint(w, "not found")
		return
	}

	io.CopyN(w, f, stat.Size())
}

type Handler interface {
	Handle(io.Writer, *http0_9.Request)
}
