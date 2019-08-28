package server

import (
	"io"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Server struct {
	Addr    string
	Handler Handler
}

type FileServer struct {
	Root http.FileSystem
}

type Handler interface {
	Handle(io.Writer, *http0_9.Request)
}
