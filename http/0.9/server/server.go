package server

import (
	"fmt"
	"io"
	"net"

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
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	r := new(http0_9.Request)
	if _, err := r.ReadFrom(conn); err != nil {
		fmt.Fprintf(conn, "failed to serve: %s\n", err)
		return
	}

	if r.Method != http0_9.MethodGet {
		fmt.Fprint(conn, "method not allowed")
		return
	}

	s.Handler.Handle(conn, r)
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

type HandlerFunc func(io.Writer, *http0_9.Request)

func (f HandlerFunc) Handle(w io.Writer, r *http0_9.Request) {
	f(w, r)
}
