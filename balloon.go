package firego

import (
	"context"
	"net/http"
)

// GetBalloon returns the current balloon device configuration (GET /balloon).
func (c *Client) GetBalloon(ctx context.Context) (*Balloon, error) {
	var out Balloon
	return &out, c.doJSON(ctx, http.MethodGet, "/balloon", nil, &out)
}

// PutBalloon creates or replaces the balloon device (pre-boot, PUT /balloon).
// Statistics collection is disabled when StatsPollingIntervalS is omitted or zero.
func (c *Client) PutBalloon(ctx context.Context, cfg *Balloon) error {
	return c.doJSON(ctx, http.MethodPut, "/balloon", cfg, nil)
}

// PatchBalloon updates the target size of the balloon device post-boot
// (PATCH /balloon). The guest kernel will inflate or deflate the balloon
// to approach the new target.
func (c *Client) PatchBalloon(ctx context.Context, update *BalloonUpdate) error {
	return c.doJSON(ctx, http.MethodPatch, "/balloon", update, nil)
}

// GetBalloonStats returns the latest statistics reported by the balloon device
// (GET /balloon/statistics). Statistics must be enabled via [PutBalloon] or
// [PatchBalloonStats] before this endpoint is available.
func (c *Client) GetBalloonStats(ctx context.Context) (*BalloonStats, error) {
	var out BalloonStats
	return &out, c.doJSON(ctx, http.MethodGet, "/balloon/statistics", nil, &out)
}

// PatchBalloonStats changes the balloon statistics polling interval post-boot
// (PATCH /balloon/statistics). Set to 0 to disable statistics collection.
func (c *Client) PatchBalloonStats(ctx context.Context, update *BalloonStatsUpdate) error {
	return c.doJSON(ctx, http.MethodPatch, "/balloon/statistics", update, nil)
}

// StartBalloonHinting initiates a free page hinting run
// (PATCH /balloon/hinting/start). The guest kernel must support the feature.
func (c *Client) StartBalloonHinting(ctx context.Context, cmd *BalloonStartCmd) error {
	return c.doJSON(ctx, http.MethodPatch, "/balloon/hinting/start", cmd, nil)
}

// GetBalloonHintingStatus returns the current state of the free page hinting run
// (GET /balloon/hinting/status).
func (c *Client) GetBalloonHintingStatus(ctx context.Context) (*BalloonHintingStatus, error) {
	var out BalloonHintingStatus
	return &out, c.doJSON(ctx, http.MethodGet, "/balloon/hinting/status", nil, &out)
}

// StopBalloonHinting halts an in-progress free page hinting run
// (PATCH /balloon/hinting/stop).
func (c *Client) StopBalloonHinting(ctx context.Context) error {
	return c.doJSON(ctx, http.MethodPatch, "/balloon/hinting/stop", struct{}{}, nil)
}
