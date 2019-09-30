package ip

import (
	"fmt"

	"github.com/tomocy/wataame/ip/addr"
)

type Addr string

func (a Addr) Proto() string {
	return "IPv4"
}

func (a Addr) Compensate() (string, error) {
	host, port, err := addr.ParseHostPort(string(a))
	if err != nil {
		return "", nil
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
