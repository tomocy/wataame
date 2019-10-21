package client

import (
	"context"
	"fmt"
	"net"

	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/ip"
	"github.com/tomocy/wataame/tcp"
)

type Client struct {
	Dialer tcp.Dialer
}

func (c *Client) Do(ctx context.Context, r *http1_0.FullRequest) (*http1_0.FullResponse, error) {
	assureFullRequest(r)
	conn, err := c.dialForFullRequest(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to do full request: %s", err)
	}

	if _, err := r.WriteTo(conn); err != nil {
		return nil, fmt.Errorf("failed to do full request: %s", err)
	}

	respCh, errCh := c.receive(conn)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to do full request: %s", ctx.Err())
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, fmt.Errorf("failed to do full request: %s", err)
	}
}

func assureFullRequest(r *http1_0.FullRequest) {
	r.RequestLine.Version = &http1_0.Version{
		Major: 1, Minor: 0,
	}
}

func (c *Client) dialForFullRequest(ctx context.Context, r *http1_0.FullRequest) (net.Conn, error) {
	addr, err := ip.Addr(r.RequestLine.URI.Host).Compensate()
	if err != nil {
		return nil, fmt.Errorf("failed to dial for full request: %s", err)
	}

	conn, err := c.dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial for full request: %s", err)
	}

	return conn, nil
}

func (c *Client) dial(ctx context.Context, addr string) (net.Conn, error) {
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

func (c *Client) receive(conn net.Conn) (<-chan *http1_0.FullResponse, <-chan error) {
	respCh, errCh := make(chan *http1_0.FullResponse), make(chan error)
	go func() {
		defer func() {
			close(respCh)
			close(errCh)
		}()

		resp := new(http1_0.FullResponse)
		if _, err := resp.ReadFrom(conn); err != nil {
			errCh <- fmt.Errorf("failed to recieve: %s", err)
			return
		}

		respCh <- resp
	}()

	return respCh, errCh
}
