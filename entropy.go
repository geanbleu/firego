package firego

import (
	"context"
	"net/http"
)

// PutEntropyDevice enables the virtio-rng entropy device
// (pre-boot, PUT /entropy).
// The guest kernel reads randomness from this device via /dev/hwrng.
func (c *Client) PutEntropyDevice(ctx context.Context, dev *EntropyDevice) error {
	return c.doJSON(ctx, http.MethodPut, "/entropy", dev, nil)
}
