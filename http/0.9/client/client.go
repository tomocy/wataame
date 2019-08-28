package client

import "net"

type Client struct {
	Dialer Dialer
}

type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
}
