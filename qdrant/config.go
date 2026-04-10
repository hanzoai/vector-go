package qdrant

import (
	"fmt"
)

const (
	apiKeyHeader    = "api-key"
	defaultHost     = "localhost"
	defaultZapPort_ = 6335
)

// Configuration options for the client.
type Config struct {
	// Hostname of the Qdrant server. Defaults to "localhost".
	Host string
	// ZAP port of the Qdrant server. Defaults to 6335.
	Port int
	// API key to use for authentication. Defaults to "".
	APIKey string
	// Whether to check compatibility between server's version and client's. Defaults to false.
	SkipCompatibilityCheck bool
	// PoolSize specifies the number of connections to create.
	// If 0, the default of 3 will be used.
	// If 1 a single connection is used (aka no pool).
	// If greater than 1, a pool of connections is created and requests are distributed in a round-robin fashion.
	PoolSize uint
}

// getZapAddr returns the ZAP address string.
func (c *Config) getZapAddr() string {
	host := c.Host
	if host == "" {
		host = defaultHost
	}
	port := c.Port
	if port == 0 {
		port = defaultZapPort_
	}
	return fmt.Sprintf("%s:%d", host, port)
}
