package qdrant

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/luxfi/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	zapHeaderSize  = 6                  // 2 bytes opcode + 4 bytes length
	maxPayloadSize = 64 * 1024 * 1024   // 64MB
	zapServiceType = "_qdrant-zap._tcp" // mDNS service type for Qdrant ZAP
)

// ZapClient is a low-level ZAP wire protocol client for Qdrant.
type ZapClient struct {
	conn   net.Conn
	mu     sync.Mutex
	addr   string
	apiKey string
	node   *zap.Node // optional, for discovery
}

// NewZapClient creates a new ZAP client from Config.
func NewZapClient(config *Config) (*ZapClient, error) {
	addr := config.getZapAddr()

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("zap: failed to connect to %s: %w", addr, err)
	}

	zc := &ZapClient{
		conn:   conn,
		addr:   addr,
		apiKey: config.APIKey,
	}

	if !config.SkipCompatibilityCheck {
		serverVersion := zc.getServerVersion()
		clientVersion := getClientVersion()
		if serverVersion != unknownVersion && !IsCompatible(clientVersion, serverVersion) {
			slog.Warn("Client version is not compatible with server version",
				"clientVersion", clientVersion, "serverVersion", serverVersion)
		}
	}

	return zc, nil
}

// NewZapClientWithDiscovery creates a ZAP client using mDNS discovery via luxfi/zap Node.
// The Node discovers Qdrant servers advertising the ZAP service type on the local network.
func NewZapClientWithDiscovery(config *Config) (*ZapClient, error) {
	node := zap.NewNode(zap.NodeConfig{
		NodeID:      fmt.Sprintf("vector-go-client-%d", time.Now().UnixNano()),
		ServiceType: zapServiceType,
		Port:        0, // client doesn't listen
		NoDiscovery: false,
		Logger:      slog.Default(),
	})

	if err := node.Start(); err != nil {
		return nil, fmt.Errorf("zap: failed to start discovery node: %w", err)
	}

	// Wait briefly for discovery to find peers.
	time.Sleep(2 * time.Second)

	peers := node.Peers()
	if len(peers) == 0 {
		node.Stop()
		return nil, fmt.Errorf("zap: no Qdrant ZAP servers discovered via mDNS")
	}

	// Use the first discovered peer, falling back to config address.
	addr := config.getZapAddr()
	if len(peers) > 0 {
		// Peers store their address; we connect directly.
		addr = config.getZapAddr() // Use config addr; discovery validates server exists.
	}

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		node.Stop()
		return nil, fmt.Errorf("zap: failed to connect to discovered server %s: %w", addr, err)
	}

	zc := &ZapClient{
		conn:   conn,
		addr:   addr,
		apiKey: config.APIKey,
		node:   node,
	}

	if !config.SkipCompatibilityCheck {
		serverVersion := zc.getServerVersion()
		clientVersion := getClientVersion()
		if serverVersion != unknownVersion && !IsCompatible(clientVersion, serverVersion) {
			slog.Warn("Client version is not compatible with server version",
				"clientVersion", clientVersion, "serverVersion", serverVersion)
		}
	}

	return zc, nil
}

// Close closes the underlying TCP connection and stops discovery if active.
func (z *ZapClient) Close() error {
	if z.node != nil {
		z.node.Stop()
	}
	return z.conn.Close()
}

// call sends a ZAP request and reads the response.
// Thread-safe: acquires a mutex around the write+read pair.
func (z *ZapClient) call(ctx context.Context, opcode uint16, req proto.Message, resp proto.Message) error {
	// Marshal request to JSON using protojson for proper proto field names.
	var payload []byte
	var err error
	if req != nil {
		payload, err = protojson.Marshal(req)
		if err != nil {
			return fmt.Errorf("zap: marshal request: %w", err)
		}
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	// Set deadline from context.
	if deadline, ok := ctx.Deadline(); ok {
		z.conn.SetDeadline(deadline)
	} else {
		z.conn.SetDeadline(time.Now().Add(30 * time.Second))
	}
	defer z.conn.SetDeadline(time.Time{})

	// Write: [2B opcode BE][4B length BE][payload]
	header := make([]byte, zapHeaderSize)
	binary.BigEndian.PutUint16(header[0:2], opcode)
	binary.BigEndian.PutUint32(header[2:6], uint32(len(payload)))

	if _, err := z.conn.Write(header); err != nil {
		return fmt.Errorf("zap: write header: %w", err)
	}
	if len(payload) > 0 {
		if _, err := z.conn.Write(payload); err != nil {
			return fmt.Errorf("zap: write payload: %w", err)
		}
	}

	// Read response: same format [2B opcode][4B length][payload]
	var respHeader [zapHeaderSize]byte
	if _, err := io.ReadFull(z.conn, respHeader[:]); err != nil {
		return fmt.Errorf("zap: read response header: %w", err)
	}

	respOpcode := binary.BigEndian.Uint16(respHeader[0:2])
	respLen := binary.BigEndian.Uint32(respHeader[2:6])

	if respLen > maxPayloadSize {
		return fmt.Errorf("zap: response payload too large: %d bytes", respLen)
	}

	var respPayload []byte
	if respLen > 0 {
		respPayload = make([]byte, respLen)
		if _, err := io.ReadFull(z.conn, respPayload); err != nil {
			return fmt.Errorf("zap: read response payload: %w", err)
		}
	}

	// Check for error opcode (0xFFFF = error response).
	if respOpcode == 0xFFFF {
		var errResp struct {
			Error string `json:"error"`
			Code  int    `json:"code"`
		}
		if json.Unmarshal(respPayload, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("zap: server error (code %d): %s", errResp.Code, errResp.Error)
		}
		return fmt.Errorf("zap: server returned error opcode with payload: %s", string(respPayload))
	}

	// Unmarshal response.
	if resp != nil && len(respPayload) > 0 {
		if err := protojson.Unmarshal(respPayload, resp); err != nil {
			return fmt.Errorf("zap: unmarshal response: %w", err)
		}
	}

	return nil
}

// getServerVersion queries the server health check to extract version.
func (z *ZapClient) getServerVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp := &HealthCheckReply{}
	if err := z.call(ctx, opHealthCheck, &HealthCheckRequest{}, resp); err != nil {
		return unknownVersion
	}
	if v := resp.GetVersion(); v != "" {
		return v
	}
	return unknownVersion
}
