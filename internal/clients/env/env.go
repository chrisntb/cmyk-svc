package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

// String returns a string representation of the client
func (c Client) String() string {
	return fmt.Sprintf("MOCK_MODE=%s, IsMockMode=%t", c.MockModeEnv(), c.IsMockMode())
}

// MockMode returns the MOCK_MODE env value
func (c Client) MockModeEnv() string {
	v := strings.TrimSpace(os.Getenv("MOCK_MODE"))
	if len(v) > 0 {
		return v
	}
	return "NOT_SET"
}

// IsMockMode returns true if MOCK_MODE env var is set to a truthy value
// -> Truthy: 1, t, T, TRUE, true, True
// -> NOT Truthy: 0, f, F, FALSE, false, False
func (c Client) IsMockMode() bool {
	val, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("MOCK_MODE")))
	return val
}
