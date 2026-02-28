package socks5

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"golang.org/x/net/proxy"
)

type Client struct {
	dialer proxy.Dialer
}

func New(proxyURL string) (*Client, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed parsing SOCKS5 proxy URL: %w", err)
	}

	dialer, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed creating SOCKS5 dialer: %w", err)
	}

	return &Client{dialer: dialer}, nil
}

func (c *Client) Dial(_ context.Context, network, addr string) (net.Conn, error) {
	return c.dialer.Dial(network, addr)
}
