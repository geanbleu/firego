package firego

import (
	"context"
	"net/http"
)

// GetInstanceInfo returns the current state and identity of the MicroVM instance
// (GET /).
func (c *Client) GetInstanceInfo(ctx context.Context) (*InstanceInfo, error) {
	var out InstanceInfo
	return &out, c.doJSON(ctx, http.MethodGet, "/", nil, &out)
}

// GetVersion returns the Firecracker build version string (GET /version).
func (c *Client) GetVersion(ctx context.Context) (*FirecrackerVersion, error) {
	var out FirecrackerVersion
	return &out, c.doJSON(ctx, http.MethodGet, "/version", nil, &out)
}

// GetVMConfig returns a snapshot of the complete VM configuration as currently
// applied (GET /vm/config).
func (c *Client) GetVMConfig(ctx context.Context) (*FullVmConfiguration, error) {
	var out FullVmConfiguration
	return &out, c.doJSON(ctx, http.MethodGet, "/vm/config", nil, &out)
}
