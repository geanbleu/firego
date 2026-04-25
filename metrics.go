package firego

import (
	"context"
	"net/http"
)

// PutMetrics initializes the metrics subsystem (PUT /metrics).
// MetricsPath must be an existing named pipe or file. Metrics are written as
// newline-delimited JSON objects. Call [Client.FlushMetrics] to force an
// immediate flush outside of the normal periodic schedule.
func (c *Client) PutMetrics(ctx context.Context, cfg *Metrics) error {
	return c.doJSON(ctx, http.MethodPut, "/metrics", cfg, nil)
}
