package qdrant

import (
	"fmt"
	"sync/atomic"
)

// Client is a high-level client for Qdrant.
// It can manage a single connection or a pool of connections, chosen by setting
// PoolSize in the Config.
type Client struct {
	clients []*ZapClient
	next    uint32
}

// NewClient creates a new Qdrant client.
// It checks Config.PoolSize to determine whether to create a single client
// or a pool of clients. If PoolSize > 1, requests are distributed across
// the connections in a round-robin fashion.
func NewClient(config *Config) (*Client, error) {
	// Ensure config is not modified for the caller by cloning.
	cfgCopy := *config
	if cfgCopy.PoolSize == 0 {
		cfgCopy.PoolSize = 3
	}
	// Create the client, with an inner connection pool of ZAP clients
	client := &Client{
		clients: make([]*ZapClient, 0, cfgCopy.PoolSize),
	}
	// Iterate over the pool size to create the individual client.
	for i := range cfgCopy.PoolSize {
		if i > 0 {
			// In case of a pool, we only want to check compatibility once.
			cfgCopy.SkipCompatibilityCheck = true
		}
		zapClient, err := NewZapClient(&cfgCopy)
		if err != nil {
			// Close already opened clients before returning an error
			client.Close()
			return nil, fmt.Errorf("failed to create client %d in pool: %w", i, err)
		}
		client.clients = append(client.clients, zapClient)
	}
	// Return the client
	return client, nil
}

// DefaultClient creates a client with the default configuration.
// Connects to localhost:6335 (ZAP port).
func DefaultClient() (*Client, error) {
	return NewClient(&Config{})
}

// get returns the next ZapClient from the pool in a round-robin fashion.
func (c *Client) get() *ZapClient {
	if len(c.clients) == 1 {
		return c.clients[0]
	}
	// Atomically increment and wrap around the counter
	idx := atomic.AddUint32(&c.next, 1) - 1
	return c.clients[idx%uint32(len(c.clients))]
}

// Close tears down all underlying connections.
func (c *Client) Close() error {
	var lastErr error
	for _, client := range c.clients {
		if err := client.Close(); err != nil {
			lastErr = err
		}
	}
	c.clients = nil // Clear the slice
	return lastErr
}

// Creates a pointer to a value of any type.
func PtrOf[T any](t T) *T {
	return &t
}
