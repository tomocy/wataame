package server

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) ListenAndServe() error {
	l, err := s.listen()
	if err != nil {
		return err
	}
	defer l.Close()

	return s.Serve(l)
}

func (s *Server) listen() (net.Listener, error) {
	compensated, err := http.Addr(s.Addr).Compensate()
	if err != nil {
		return nil, err
	}

	return net.Listen("tcp", compensated)
}

func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer conn.Close()
			r, err := readRequest(conn)
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}

			if r.Method != http.MethodGet {
				fmt.Fprintln(conn, "method not allowed")
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
			return nil, err
		}
		b.Write(line)
		if !isPrefix {
			break
		}
	}

	return parseRequest(b.Bytes())
}

func parseRequest(bs []byte) (*http0_9.Request, error) {
	splited := strings.Split(string(bs), " ")
	if len(splited) < 2 {
		return nil, errors.New("invalid format of request: missing space between method and uri")
	}
	uri, err := url.Parse("http://" + splited[1])
	if err != nil {
		return nil, err
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
