package server

import http1_0 "github.com/tomocy/wataame/http/1.0"

type Server struct {
	Addr          string
	SimpleHandler SimpleHandler
}

type SimpleHandler interface {
	Handle(*http1_0.SimpleResponse, *http1_0.SimpleRequest)
}

type SimpleHandlerFunc func(*http1_0.SimpleResponse, *http1_0.SimpleRequest)
