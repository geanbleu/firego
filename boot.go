package firego

import (
	"context"
	"net/http"
)

// PutBootSource configures the guest kernel and optional initrd
// (pre-boot, PUT /boot-source).
// This must be called before [Client.StartInstance].
func (c *Client) PutBootSource(ctx context.Context, src *BootSource) error {
	return c.doJSON(ctx, http.MethodPut, "/boot-source", src, nil)
}
