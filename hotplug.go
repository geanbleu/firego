package firego

import (
	"context"
	"net/http"
)

// PutHotplugMemory configures the virtio-mem hotpluggable memory device
// (pre-boot, PUT /hotplug/memory).
// TotalSizeMib defines the upper bound of memory that can later be plugged in
// without restarting the VM.
func (c *Client) PutHotplugMemory(ctx context.Context, cfg *MemoryHotplugConfig) error {
	return c.doJSON(ctx, http.MethodPut, "/hotplug/memory", cfg, nil)
}

// PatchHotplugMemory requests a change to the amount of memory currently plugged
// into the guest (post-boot, PATCH /hotplug/memory).
// RequestedSizeMib must be a multiple of SlotSizeMib and cannot exceed TotalSizeMib.
func (c *Client) PatchHotplugMemory(ctx context.Context, update *MemoryHotplugSizeUpdate) error {
	return c.doJSON(ctx, http.MethodPatch, "/hotplug/memory", update, nil)
}

// GetHotplugMemory returns the current state of the virtio-mem device,
// including the plugged and requested sizes (GET /hotplug/memory).
func (c *Client) GetHotplugMemory(ctx context.Context) (*MemoryHotplugStatus, error) {
	var out MemoryHotplugStatus
	return &out, c.doJSON(ctx, http.MethodGet, "/hotplug/memory", nil, &out)
}
