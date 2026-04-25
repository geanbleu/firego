package firego

import (
	"context"
	"net/http"
)

// PatchVM requests a state transition for the running MicroVM (PATCH /vm).
// Use [PauseVM] and [ResumeVM] for the common pause/resume flow.
func (c *Client) PatchVM(ctx context.Context, state *Vm) error {
	return c.doJSON(ctx, http.MethodPatch, "/vm", state, nil)
}

// PauseVM suspends execution of all vCPUs. The VM remains in memory and all
// devices are still accessible from the host. Required before [CreateSnapshot].
func (c *Client) PauseVM(ctx context.Context) error {
	return c.PatchVM(ctx, &Vm{State: VmStatePaused})
}

// ResumeVM resumes a VM that was paused with [PauseVM] or loaded from a
// snapshot with ResumeVM set to false.
func (c *Client) ResumeVM(ctx context.Context) error {
	return c.PatchVM(ctx, &Vm{State: VmStateResumed})
}
