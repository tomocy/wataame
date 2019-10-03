package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"

	http0_9 "github.com/tomocy/wataame/http/0.9"
	"github.com/tomocy/wataame/ip"
	"github.com/tomocy/wataame/tcp"
)

type Client struct {
	Dialer tcp.Dialer
}

func (c *Client) Do(ctx context.Context, r *http0_9.Request) (http0_9.Response, error) {
	conn, err := c.dialForRequest(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to do: %s", err)
	}

	if _, err := r.WriteTo(conn); err != nil {
		return nil, fmt.Errorf("failed to do: %s", err)
	}

	respCh, errCh := c.receive(conn)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to do: %s", ctx.Err())
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, fmt.Errorf("failed to do: %s", err)
	}
}

func (c *Client) dialForRequest(ctx context.Context, r *http0_9.Request) (tcp.Conn, error) {
	addr, err := ip.Addr(r.URI.Host).Compensate()
	if err != nil {
		return nil, fmt.Errorf("failed to dial for request: %s", err)
	}

	conn, err := c.dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial for request: %s", err)
	}

	return conn, nil
}

func (c *Client) dial(ctx context.Context, addr string) (tcp.Conn, error) {
	d := c.Dialer
	if d == nil {
		d = new(tcp.GoDialer)
	}

	conn, err := d.Dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	return conn, nil
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
			errCh <- fmt.Errorf("failed to recieve: %s", err)
			return
		}

		respCh <- resp
	}()

	return respCh, errCh
}
