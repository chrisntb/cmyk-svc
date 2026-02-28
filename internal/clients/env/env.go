package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Client struct{}

var (
	notSet = "NOT_SET"
)

func New() *Client {
	return &Client{}
}

// String returns a string representation of the client
func (c Client) String() string {
	return fmt.Sprintf("MOCK_MODE=%s, IsMockMode=%t, HTTP_PROXY=%s, HTTPS_PROXY=%s, SOCKS5_PROXY=%s", c.MockModeEnv(), c.IsMockMode(), c.HttpProxyEnv(), c.HttpsProxyEnv(), c.Socks5ProxyEnv())
}

// MockMode returns the MOCK_MODE env value
func (c Client) MockModeEnv() string {
	v := strings.TrimSpace(os.Getenv("MOCK_MODE"))
	if len(v) > 0 {
		return v
	}
	return notSet
}

// IsMockMode returns true if MOCK_MODE env var is set to a truthy value
// -> Truthy: 1, t, T, TRUE, true, True
// -> NOT Truthy: 0, f, F, FALSE, false, False
func (c Client) IsMockMode() bool {
	val, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("MOCK_MODE")))
	return val
}

// HttpProxyEnv returns the HTTP_PROXY env value
func (c Client) HttpProxyEnv() string {
	v := strings.TrimSpace(os.Getenv("HTTP_PROXY"))
	if len(v) > 0 {
		return v
	}
	return notSet
}

// HttpProxyMode returns true if HTTP_PROXY env var is set else false
func (c Client) HttpProxyMode() bool {
	return strings.TrimSpace(os.Getenv("HTTP_PROXY")) != ""
}

// HttpsProxyEnv returns the HTTPS_PROXY env value
func (c Client) HttpsProxyEnv() string {
	v := strings.TrimSpace(os.Getenv("HTTPS_PROXY"))
	if len(v) > 0 {
		return v
	}
	return notSet
}

// HttpsProxyMode returns true if HTTPS_PROXY env var is set else false
func (c Client) HttpsProxyMode() bool {
	return strings.TrimSpace(os.Getenv("HTTPS_PROXY")) != ""
}

// Socks5ProxyEnv returns the SOCKS5_PROXY env value
func (c Client) Socks5ProxyEnv() string {
	v := strings.TrimSpace(os.Getenv("SOCKS5_PROXY"))
	if len(v) > 0 {
		return v
	}
	return notSet
}

// Socks5ProxyMode returns true if SOCKS5_PROXY env var is set else false
func (c Client) Socks5ProxyMode() bool {
	return strings.TrimSpace(os.Getenv("SOCKS5_PROXY")) != ""
}
