package firego

import (
	"context"
	"net/http"
)

// GetMachineConfig returns the current machine configuration (GET /machine-config).
func (c *Client) GetMachineConfig(ctx context.Context) (*MachineConfiguration, error) {
	var out MachineConfiguration
	return &out, c.doJSON(ctx, http.MethodGet, "/machine-config", nil, &out)
}

// PutMachineConfig sets the machine configuration, replacing any existing values
// (pre-boot, PUT /machine-config).
func (c *Client) PutMachineConfig(ctx context.Context, cfg *MachineConfiguration) error {
	return c.doJSON(ctx, http.MethodPut, "/machine-config", cfg, nil)
}

// PatchMachineConfig partially updates the machine configuration (pre-boot,
// PATCH /machine-config). Fields left nil in cfg retain their current values.
func (c *Client) PatchMachineConfig(ctx context.Context, cfg *MachineConfiguration) error {
	return c.doJSON(ctx, http.MethodPatch, "/machine-config", cfg, nil)
}
