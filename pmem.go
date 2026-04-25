package firego

import (
	"context"
	"net/http"
)

// PutPmem creates or replaces a persistent memory device identified by id
// (pre-boot, PUT /pmem/{id}).
// The id in the path must match Pmem.ID.
func (c *Client) PutPmem(ctx context.Context, id string, pmem *Pmem) error {
	return c.doJSON(ctx, http.MethodPut, "/pmem/"+id, pmem, nil)
}

// PatchPmem updates the rate limiter of a pmem device post-boot
// (PATCH /pmem/{id}). Only the rate limiter can be changed after boot.
func (c *Client) PatchPmem(ctx context.Context, id string, update *PartialPmem) error {
	return c.doJSON(ctx, http.MethodPatch, "/pmem/"+id, update, nil)
}
