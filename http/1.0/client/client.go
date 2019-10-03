package client

import (
	"context"
	"fmt"

	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/tcp"
)

type Client struct {
	Dialer tcp.Dialer
}

func (c *Client) Do(ctx context.Context, r http1_0.Request) (http1_0.Response, error) {
	switch r.(type) {
	default:
		panic(fmt.Sprintf("failed to do: unsupported request type: %T", r))
	}
}
