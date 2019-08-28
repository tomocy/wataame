package client

import (
	"net"

	http "github.com/tomocy/wataame/http/0.9"
)

type Client struct {
	Dialer Dialer
}

func (c *Client) dialForRequest(r *http.Request) (net.Conn, error) {
	addr, err := r.Address()
	if err != nil {
		return nil, err
	}

	return c.dial("tcp", addr)
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
