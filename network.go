package firego

import (
	"context"
	"net/http"
)

// PutNetworkInterface creates or replaces a network interface identified by ifaceID
// (pre-boot, PUT /network-interfaces/{iface_id}).
// The ifaceID in the path must match NetworkInterface.IfaceID.
func (c *Client) PutNetworkInterface(ctx context.Context, ifaceID string, iface *NetworkInterface) error {
	return c.doJSON(ctx, http.MethodPut, "/network-interfaces/"+ifaceID, iface, nil)
}

// PatchNetworkInterface updates the rate limiters of a network interface post-boot
// (PATCH /network-interfaces/{iface_id}). Only the rate limiter fields may be changed
// after the VM has started; structural changes require a new VM.
func (c *Client) PatchNetworkInterface(ctx context.Context, ifaceID string, update *PartialNetworkInterface) error {
	return c.doJSON(ctx, http.MethodPatch, "/network-interfaces/"+ifaceID, update, nil)
}
