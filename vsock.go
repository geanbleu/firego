package firego

import (
	"context"
	"net/http"
)

// PutVsock creates or replaces the virtio-vsock device
// (pre-boot, PUT /vsock).
// The host-side proxy socket at UDSPath must be created by the caller before
// the VM starts; Firecracker connects to it at boot time.
func (c *Client) PutVsock(ctx context.Context, vsock *Vsock) error {
	return c.doJSON(ctx, http.MethodPut, "/vsock", vsock, nil)
}
