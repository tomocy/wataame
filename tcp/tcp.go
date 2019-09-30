package tcp

import (
	"net"
)

type Conn interface {
	net.Conn
}
