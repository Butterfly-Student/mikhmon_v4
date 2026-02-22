package routeros

import (
	"fmt"

	"github.com/go-routeros/routeros/v3"
)

// Connect opens a RouterOS API connection to the given host (host:port).
func Connect(host, user, password string) (*routeros.Client, error) {
	c, err := routeros.Dial(host, user, password)
	if err != nil {
		return nil, fmt.Errorf("routeros connect %s: %w", host, err)
	}
	return c, nil
}

// RunArgs executes a RouterOS command given as a flat []string:
//
//	["/ip/hotspot/user/print", "?name=foo", "=comment=bar", ...]
//
// The go-routeros/v3 RunArgs method expects []string (not variadic).
// This wrapper is variadic for convenience at call sites.
func RunArgs(c *routeros.Client, args ...string) ([]map[string]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("routeros: no args provided")
	}
	reply, err := c.RunArgs(args) // pass as []string (not variadic)
	if err != nil {
		return nil, err
	}
	var result []map[string]string
	for _, re := range reply.Re {
		m := make(map[string]string)
		for _, p := range re.List {
			m[p.Key] = p.Value
		}
		result = append(result, m)
	}
	return result, nil
}
