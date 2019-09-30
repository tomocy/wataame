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

type Conn interface {
	net.Conn
}
