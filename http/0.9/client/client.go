package client

import (
	"context"
	"net"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Client struct {
	Dialer Dialer
}

func (c *Client) Do(r *http0_9.Request) (http0_9.Response, error) {
	conn, err := c.dialForRequest(context.TODO(), r)
	if err != nil {
		return nil, err
	}

	if err := r.Write(conn); err != nil {
		return nil, err
	}
	resp := make(http0_9.Response, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}

	return resp[:n], nil
}

func (c *Client) dialForRequest(ctx context.Context, r *http0_9.Request) (net.Conn, error) {
	addr, err := http.Address(r.URI.Host).Compensate()
	if err != nil {
		return nil, err
	}

	return c.dial(ctx, "tcp", addr)
}

func (c *Client) dial(ctx context.Context, network, addr string) (net.Conn, error) {
	d := c.Dialer
	if d == nil {
		d = new(DefaultDialer)
	}

	return d.Dial(ctx, network, addr)
}

type Dialer interface {
	Dial(ctx context.Context, network, addr string) (net.Conn, error)
}

type DefaultDialer struct {
	net.Dialer
}

func (d *DefaultDialer) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, network, addr)
}
