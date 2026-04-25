package firego

import (
	"context"
	"net/http"
)

// PutDrive creates or replaces a block device identified by driveID
// (pre-boot, PUT /drives/{drive_id}).
// The driveID in the path must match Drive.DriveID.
func (c *Client) PutDrive(ctx context.Context, driveID string, drive *Drive) error {
	return c.doJSON(ctx, http.MethodPut, "/drives/"+driveID, drive, nil)
}

// PatchDrive updates a drive's host path or rate limiter while the VM is running
// (post-boot, PATCH /drives/{drive_id}). Only the fields set in update are changed.
func (c *Client) PatchDrive(ctx context.Context, driveID string, update *PartialDrive) error {
	return c.doJSON(ctx, http.MethodPatch, "/drives/"+driveID, update, nil)
}
