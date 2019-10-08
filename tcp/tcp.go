package tcp

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, addr string) (net.Conn, error)
}

type GoDialer struct {
	net.Dialer
}

func (d *GoDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, "tcp", addr)
}
