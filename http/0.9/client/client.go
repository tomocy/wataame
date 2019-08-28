package client

import "net"

type Client struct{}

type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
}
