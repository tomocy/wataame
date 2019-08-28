package server

import (
	"io"

	http "github.com/tomocy/wataame/http/0.9"
)

type Server struct {
	Addr string
}

type Handler interface {
	Handle(io.Writer, *http.Request)
}
