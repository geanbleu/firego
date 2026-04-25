package firego

import (
	"context"
	"net/http"
)

// PutLogger initializes the logging subsystem (PUT /logger).
// This can be called before or after boot. The LogPath must be an existing
// named pipe or file; Firecracker does not create it automatically.
func (c *Client) PutLogger(ctx context.Context, cfg *Logger) error {
	return c.doJSON(ctx, http.MethodPut, "/logger", cfg, nil)
}
