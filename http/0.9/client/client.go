package client

import (
	"net"
)

type Client struct {
	Dialer Dialer
}

func (c *Client) dial(network, addr string) (net.Conn, error) {
	d := c.Dialer
	if d == nil {
		d = new(net.Dialer)
	}

	return d.Dial(network, addr)
}

type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
}
