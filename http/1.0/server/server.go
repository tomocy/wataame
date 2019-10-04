package server

import (
	"io"

	http1_0 "github.com/tomocy/wataame/http/1.0"
)

type Handler interface {
	Handle(io.Writer, http1_0.Request)
}

type HandlerFunc func(io.Writer, http1_0.Request)
