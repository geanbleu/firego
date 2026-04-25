package firego

import (
	"context"
	"net/http"
)

// PutSerial configures the emulated serial console output (PUT /serial).
// SerialOutPath must be an existing named pipe or file opened for writing.
func (c *Client) PutSerial(ctx context.Context, dev *SerialDevice) error {
	return c.doJSON(ctx, http.MethodPut, "/serial", dev, nil)
}
