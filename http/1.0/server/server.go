package server

import (
	"fmt"
	"net"

	"github.com/tomocy/wataame/http"
	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/ip"
)

type Server struct {
	Addr     string
	Listener net.Listener
	Handler  Handler
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

	peekable := &http.PeekableConn{
		Conn: conn,
	}
	v, err := http.DetectVersion(peekable)

	if err != nil {
		fmt.Fprintf(peekable, "failed to handle: %s", err)
		return
	}

	switch v {
	case "0.9":
		s.handleSimpleRequest(peekable)
	case "1.0":
		s.handleFullRequest(peekable)
	default:
		panic(fmt.Sprintf("failed to handle: unsupported HTTP version: %s", v))
	}
}

func (s *Server) handleSimpleRequest(conn net.Conn) {
	defer conn.Close()

	if s.Handler == nil {
		fmt.Fprintln(conn, "failed to handle simple request: handler is not set")
		return
	}

	req := new(http1_0.SimpleRequest)
	if _, err := req.ReadFrom(conn); err != nil {
		fmt.Fprintf(conn, "failed to handle simple request: %s", err)
		return
	}

	resp := new(http1_0.SimpleResponse)
	s.Handler.HandleSimpleRequest(resp, req)

	resp.WriteTo(conn)
}

func (s *Server) handleFullRequest(conn net.Conn) {
	defer conn.Close()

	if s.Handler == nil {
		fmt.Fprintln(conn, "failed to handle full request: handler is not set")
		return
	}

	req := new(http1_0.FullRequest)
	if _, err := req.ReadFrom(conn); err != nil {
		fmt.Fprintf(conn, "failed to handle full request: %s", err)
		return
	}

	resp := new(http1_0.FullResponse)
	s.Handler.HandleFullRequest(resp, req)

	assureFullResponse(resp)
	resp.WriteTo(conn)
}

func assureFullResponse(r *http1_0.FullResponse) {
	r.StatusLine.Version = &http1_0.Version{
		Major: 1, Minor: 0,
	}
}

type Handler interface {
	FullHandler
	SimpleHandler
}

type HandlerFunc struct {
	FullHandlerFunc
	SimpleHandlerFunc
}

type FullHandler interface {
	HandleFullRequest(*http1_0.FullResponse, *http1_0.FullRequest)
}

type FullHandlerFunc func(*http1_0.FullResponse, *http1_0.FullRequest)

func (f FullHandlerFunc) HandleFullRequest(resp *http1_0.FullResponse, req *http1_0.FullRequest) {
	f(resp, req)
}

type SimpleHandler interface {
	HandleSimpleRequest(*http1_0.SimpleResponse, *http1_0.SimpleRequest)
}

type SimpleHandlerFunc func(*http1_0.SimpleResponse, *http1_0.SimpleRequest)

func (f SimpleHandlerFunc) HandleSimpleRequest(resp *http1_0.SimpleResponse, req *http1_0.SimpleRequest) {
	f(resp, req)
}
