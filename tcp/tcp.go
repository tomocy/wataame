package tcp

import (
	"context"
	"net"
)

type Listener interface {
	Addr() net.Addr
	Accept(context.Context) (net.Conn, error)
	Close() error
}

type GoListener struct {
	net.TCPListener
}

type Dialer interface {
	Dial(ctx context.Context, addr string) (net.Conn, error)
}

type DialerFunc func(context.Context, string) (net.Conn, error)

func (f DialerFunc) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return f(ctx, addr)
}

type GoDialer struct {
	net.Dialer
}

func (d *GoDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, "tcp", addr)
}
