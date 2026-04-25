package firego

import (
	"context"
	"net/http"
)

// CreateSnapshot saves the current VM state and memory to disk
// (post-boot, PUT /snapshot/create).
// The VM should be paused with [Client.PauseVM] before calling this to ensure
// a consistent snapshot. Use [Client.ResumeVM] afterwards to continue execution.
func (c *Client) CreateSnapshot(ctx context.Context, params *SnapshotCreateParams) error {
	return c.doJSON(ctx, http.MethodPut, "/snapshot/create", params, nil)
}

// LoadSnapshot restores a VM from a previously created snapshot
// (pre-boot, PUT /snapshot/load).
// Exactly one of SnapshotLoadParams.MemFilePath or SnapshotLoadParams.MemBackend
// must be set. Set ResumeVM to true to start the VM immediately after loading.
func (c *Client) LoadSnapshot(ctx context.Context, params *SnapshotLoadParams) error {
	return c.doJSON(ctx, http.MethodPut, "/snapshot/load", params, nil)
}
