package firego

import (
	"context"
	"net/http"
)

// PerformAction executes a synchronous action on the MicroVM (PUT /actions).
// Use the typed helpers [StartInstance], [FlushMetrics], and [SendCtrlAltDel]
// for the most common cases.
func (c *Client) PerformAction(ctx context.Context, actionType ActionType) error {
	return c.doJSON(ctx, http.MethodPut, "/actions", &InstanceActionInfo{ActionType: actionType}, nil)
}

// StartInstance boots the MicroVM. All pre-boot configuration (boot source,
// machine config, drives, network interfaces) must be applied before this call.
func (c *Client) StartInstance(ctx context.Context) error {
	return c.PerformAction(ctx, ActionInstanceStart)
}

// FlushMetrics forces an immediate write of in-memory metrics to the configured
// metrics output. Can be called at any time while the instance is running.
func (c *Client) FlushMetrics(ctx context.Context) error {
	return c.PerformAction(ctx, ActionFlushMetrics)
}

// SendCtrlAltDel sends the Ctrl+Alt+Del key sequence to the guest OS.
// On most Linux guests this triggers a graceful reboot.
func (c *Client) SendCtrlAltDel(ctx context.Context) error {
	return c.PerformAction(ctx, ActionSendCtrlAltDel)
}
