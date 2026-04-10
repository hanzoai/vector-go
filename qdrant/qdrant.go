package qdrant

import (
	"context"
)

// Check liveliness of the service.
func (c *Client) HealthCheck(ctx context.Context) (*HealthCheckReply, error) {
	resp := &HealthCheckReply{}
	err := c.get().call(ctx, opHealthCheck, &HealthCheckRequest{}, resp)
	if err != nil {
		return nil, newQdrantErr(err, "HealthCheck")
	}
	return resp, nil
}
