package firego

import (
	"context"
	"net/http"
)

// GetMMDS retrieves the current contents of the MMDS data store (GET /mmds).
func (c *Client) GetMMDS(ctx context.Context) (MmdsContentsObject, error) {
	var out MmdsContentsObject
	return out, c.doJSON(ctx, http.MethodGet, "/mmds", nil, &out)
}

// PutMMDS creates or fully replaces the MMDS data store (PUT /mmds).
// The guest can query this data via the MMDS endpoint (default 169.254.169.254).
func (c *Client) PutMMDS(ctx context.Context, contents MmdsContentsObject) error {
	return c.doJSON(ctx, http.MethodPut, "/mmds", contents, nil)
}

// PatchMMDS performs a recursive merge update of the MMDS data store (PATCH /mmds).
// Existing keys not present in contents are preserved.
func (c *Client) PatchMMDS(ctx context.Context, contents MmdsContentsObject) error {
	return c.doJSON(ctx, http.MethodPatch, "/mmds", contents, nil)
}

// PutMMDSConfig configures the MMDS service itself (pre-boot, PUT /mmds/config).
// This must be called before [Client.StartInstance].
func (c *Client) PutMMDSConfig(ctx context.Context, cfg *MmdsConfig) error {
	return c.doJSON(ctx, http.MethodPut, "/mmds/config", cfg, nil)
}
