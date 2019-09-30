package tcp

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, addr string) (Conn, error)
}

type GoDialer struct {
	net.Dialer
}

func (d *GoDialer) Dial(ctx context.Context, addr string) (Conn, error) {
	return d.Dialer.DialContext(ctx, "tcp", addr)
}

type Conn interface {
	net.Conn
}
