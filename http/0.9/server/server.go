package server

import (
	"errors"
	"io"
	"net"
	"net/url"
	"strings"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Server struct {
	Addr    http.Address
	Handler Handler
}

func (s *Server) listen() (net.Listener, error) {
	compensated, err := s.Addr.Compensate()
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
				conn.Write([]byte(err.Error()))
				return
			}

			s.Handler.Handle(conn, r)
		}()
	}
}

func readRequest(conn net.Conn) (*http0_9.Request, error) {
	r := make([]byte, 1024)
	n, err := conn.Read(r)
	if err != nil {
		return nil, err
	}

	return parseRequest(r[:n])
}

func parseRequest(bs []byte) (*http0_9.Request, error) {
	splited := strings.Split(string(bs), "\n")
	if len(splited) < 2 {
		return nil, errors.New("invalid format of request: finishing without a new line")
	}
	splited = strings.Split(splited[0], " ")
	if len(splited) < 2 {
		return nil, errors.New("invalid format of request: missing space between method and uri")
	}
	uri, err := url.Parse("http://" + splited[1])
	if err != nil {
		return nil, err
	}

	return &http0_9.Request{
		URI: uri,
	}, nil
}

type FileServer struct {
	Root http.FileSystem
}

func (s *FileServer) Handle(w io.Writer, r *http0_9.Request) {
	f, err := s.Root.Open(r.URI.Path)
	if err != nil {
		w.Write([]byte("not found"))
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		w.Write([]byte("interanl server error"))
		return
	}
	if stat.IsDir() {
		w.Write([]byte("not found"))
		return
	}

	io.CopyN(w, f, stat.Size())
}

type Handler interface {
	Handle(io.Writer, *http0_9.Request)
}
