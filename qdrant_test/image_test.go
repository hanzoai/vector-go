package qdrant_test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const TestImage string = "qdrant/qdrant:v1.16.1"

// skipIfShort opts the caller out of running under `go test -short`.
// Integration tests in this package spin up a Qdrant container via
// testcontainers-go, which panics on hosts without a reachable Docker
// daemon (e.g. most CI lint jobs and dev workstations). Skipping under
// -short lets `go test -short ./...` succeed without Docker while leaving
// the full `go test ./...` path unchanged for CI and local runs that do
// have Docker available.
func skipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping container-backed integration test in -short mode (requires Docker)")
	}
}

// We use an instance with distributed mode enabled
// to test methods like CreateShardKey(), DeleteShardKey().
func distributedQdrant(ctx context.Context, apiKey string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        TestImage,
		ExposedPorts: []string{"6334/tcp"},
		Env: map[string]string{
			"QDRANT__CLUSTER__ENABLED": "true",
			"QDRANT__SERVICE__API_KEY": apiKey,
		},
		Cmd: []string{"./qdrant", "--uri", "http://qdrant_node_1:6335"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6334/tcp").WithStartupTimeout(5 * time.Second),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	return container, err
}

func standaloneQdrant(ctx context.Context, apiKey string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        TestImage,
		ExposedPorts: []string{"6334/tcp"},
		Env: map[string]string{
			"QDRANT__SERVICE__API_KEY": apiKey,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6334/tcp").WithStartupTimeout(5 * time.Second),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	return container, err
}
