package firego

import (
	"context"
	"net/http"
)

// PutCPUConfig sets fine-grained CPU feature flag overrides for all vCPUs
// (pre-boot, PUT /cpu-config). This is the preferred alternative to the
// deprecated [CpuTemplate] field in [MachineConfiguration].
func (c *Client) PutCPUConfig(ctx context.Context, cfg *CpuConfig) error {
	return c.doJSON(ctx, http.MethodPut, "/cpu-config", cfg, nil)
}
