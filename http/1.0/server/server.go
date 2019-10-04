package server

import (
	"io"

	http1_0 "github.com/tomocy/wataame/http/1.0"
)

type Server struct {
	Addr    string
	Handler Handler
}

type Handler interface {
	Handle(io.Writer, http1_0.Request)
}

type HandlerFunc func(io.Writer, http1_0.Request)

func (f HandlerFunc) Handle(w io.Writer, r http1_0.Request) {
	f(w, r)
}
