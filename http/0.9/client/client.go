package client

import (
	"net"

	http "github.com/tomocy/wataame/http/0.9"
)

type Client struct {
	Dialer Dialer
}

func (c *Client) Do(r *http.Request) (http.Response, error) {
	conn, err := c.dialForRequest(r)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := r.Write(conn); err != nil {
		return nil, err
	}
	resp := make(http.Response, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}

	return resp[:n], nil
}

func (c *Client) dialForRequest(r *http.Request) (net.Conn, error) {
	addr, err := r.Addr.Compensate()
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
