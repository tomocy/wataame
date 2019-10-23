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

func (l *GoListener) Accept(ctx context.Context) (net.Conn, error) {
	connCh, errCh := l.acceptAsnycly()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case conn := <-connCh:
		return conn, nil
	case err := <-errCh:
		return nil, err
	}
}

func (l *GoListener) acceptAsnycly() (<-chan net.Conn, <-chan error) {
	connCh, errCh := make(chan net.Conn), make(chan error)
	go func() {
		defer func() {
			close(connCh)
			close(errCh)
		}()

		conn, err := l.TCPListener.Accept()
		if err != nil {
			errCh <- err
			return
		}

		connCh <- conn
	}()

	return connCh, errCh
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
