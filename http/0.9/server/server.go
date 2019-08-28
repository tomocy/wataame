package server

import (
	"io"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Server struct {
	Addr    http.Address
	Handler Handler
}

type FileServer struct {
	Root http.FileSystem
}

func (s *FileServer) Handle(w io.Writer, r *http0_9.Request) {
	f, err := s.Root.Open(string(r.Addr))
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

	io.CopyN(w, f, stat.Size())
}

type Handler interface {
	Handle(io.Writer, *http0_9.Request)
}
