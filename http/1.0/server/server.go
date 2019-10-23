package server

import (
	"context"
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

func (s *Server) ListenAndServe(ctx context.Context) error {
	if s.Listener == nil {
		if err := s.listen(); err != nil {
			return fmt.Errorf("failed to listen and serve: %s", err)
		}
	}
	defer s.Listener.Close()

	if err := s.Serve(ctx); err != nil {
		return fmt.Errorf("failed to listen and serve: %s", err)
	}

	return nil
}

func (s *Server) listen() error {
	compensated, err := ip.Addr(s.Addr).Compensate()
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	l, err := net.Listen("tcp", compensated)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}
	s.Listener = l

	return nil
}

func (s *Server) Serve(ctx context.Context) error {
	connCh, errCh := s.accept(ctx)
	for {
		select {
		case conn := <-connCh:
			go s.handle(conn)
		case err := <-errCh:
			if err == context.Canceled {
				return nil
			}
			return fmt.Errorf("failed to serve: %s", err)
		}
	}
}

func (s *Server) accept(ctx context.Context) (<-chan net.Conn, <-chan error) {
	connCh, errCh := make(chan net.Conn), make(chan error)
	go func() {
		defer func() {
			close(connCh)
			close(errCh)
		}()
		defer s.Listener.Close()

		go func() {
			<-ctx.Done()
			errCh <- ctx.Err()
		}()

		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				errCh <- fmt.Errorf("failed to accept: %s", err)
				continue
			}

			connCh <- conn
		}
	}()

	return connCh, errCh
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
