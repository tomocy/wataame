package ip

import (
	"fmt"
	"strings"

	ip4 "github.com/tomocy/wataame/ip/4"
	ip6 "github.com/tomocy/wataame/ip/6"
)

type Addr string

func (a Addr) Compensate() (string, error) {
	if strings.HasPrefix(string(a), "[") {
		return a.compensateWith(ip6.Addr(a))
	}

	return a.compensateWith(ip4.Addr(a))
}

func (a Addr) compensateWith(c interface {
	Proto() string
	Compensate() (string, error)
}) (string, error) {
	addr, err := c.Compensate()
	if err != nil {
		return "", fmt.Errorf("failed to compensate address for %s: %s", c.Proto(), err)
	}

	return addr, nil
}
