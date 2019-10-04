package client

import (
	"context"
	"fmt"

	"github.com/tomocy/wataame/http/0.9/client"
	http1_0 "github.com/tomocy/wataame/http/1.0"
	"github.com/tomocy/wataame/tcp"
)

type Client struct {
	Dialer tcp.Dialer
}

func (c *Client) Do(ctx context.Context, r http1_0.Request) (http1_0.Response, error) {
	switch r := r.(type) {
	case *http1_0.SimpleRequest:
		resp, err := c.doSimpleRequest(ctx, r)
		if err != nil {
			return nil, fmt.Errorf("failed to do: %s", err)
		}
		return resp, nil
	default:
		panic(fmt.Sprintf("failed to do: unsupported request type: %T", r))
	}
}

func (c *Client) doSimpleRequest(ctx context.Context, r *http1_0.SimpleRequest) (*http1_0.SimpleResponse, error) {
	delegated := &client.Client{
		Dialer: c.Dialer,
	}
	resp, err := delegated.Do(ctx, &r.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to do simple request: %s", err)
	}

	return &http1_0.SimpleResponse{
		Response: resp,
	}, nil
}
