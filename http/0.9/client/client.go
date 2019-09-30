package client

import (
	"context"
	"io/ioutil"
	"net"

	http "github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

type Client struct {
	Dialer Dialer
}

func (c *Client) Do(ctx context.Context, r *http0_9.Request) (http0_9.Response, error) {
	conn, err := c.dialForRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	if _, err := r.WriteTo(conn); err != nil {
		return nil, err
	}
	respCh, errCh := c.receive(conn)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (c *Client) dialForRequest(ctx context.Context, r *http0_9.Request) (net.Conn, error) {
	addr, err := http.Addr(r.URI.Host).Compensate()
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

func (c *Client) receive(conn net.Conn) (<-chan http0_9.Response, <-chan error) {
	respCh, errCh := make(chan http0_9.Response), make(chan error)
	go func() {
		defer func() {
			close(respCh)
			close(errCh)
		}()

		resp, err := ioutil.ReadAll(conn)
		if err != nil {
			errCh <- err
			return
		}

		respCh <- resp
	}()

	return respCh, errCh
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
