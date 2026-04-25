// snapshot demonstrates pausing a running VM, creating a snapshot, and
// restoring it on a second Firecracker instance.
//
// Usage — create:
//
//	go run ./examples/snapshot -socket /run/fc.sock \
//	  -action create -mem /tmp/vm.mem -state /tmp/vm.state
//
// Usage — load:
//
//	go run ./examples/snapshot -socket /run/fc2.sock \
//	  -action load -mem /tmp/vm.mem -state /tmp/vm.state -resume
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/geanbleu/firego"
)

func main() {
	socket := flag.String("socket", "/run/firecracker.sock", "Firecracker Unix socket path")
	action := flag.String("action", "create", "Action to perform: create or load")
	memPath := flag.String("mem", "/tmp/vm.mem", "Path to the guest memory file")
	statePath := flag.String("state", "/tmp/vm.state", "Path to the MicroVM state file")
	diff := flag.Bool("diff", false, "Create a diff (incremental) snapshot instead of full")
	resume := flag.Bool("resume", false, "Automatically resume the VM after loading (load only)")
	flag.Parse()

	ctx := context.Background()
	c := firego.New(*socket)

	switch *action {
	case "create":
		if err := create(ctx, c, *memPath, *statePath, *diff); err != nil {
			log.Fatal(err)
		}
	case "load":
		if err := load(ctx, c, *memPath, *statePath, *resume); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown action %q: must be create or load\n", *action)
		os.Exit(1)
	}
}

func create(ctx context.Context, c *firego.Client, memPath, statePath string, diff bool) error {
	log.Println("pausing VM...")
	if err := c.PauseVM(ctx); err != nil {
		return fmt.Errorf("pause: %w", err)
	}

	snapType := firego.SnapshotTypeFull
	if diff {
		snapType = firego.SnapshotTypeDiff
	}

	log.Printf("creating %s snapshot: mem=%s state=%s", snapType, memPath, statePath)
	if err := c.CreateSnapshot(ctx, &firego.SnapshotCreateParams{
		MemFilePath:  memPath,
		SnapshotPath: statePath,
		SnapshotType: &snapType,
	}); err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	log.Println("resuming VM...")
	if err := c.ResumeVM(ctx); err != nil {
		return fmt.Errorf("resume: %w", err)
	}

	log.Println("snapshot created successfully")
	return nil
}

func load(ctx context.Context, c *firego.Client, memPath, statePath string, resume bool) error {
	log.Printf("loading snapshot: mem=%s state=%s resume=%v", memPath, statePath, resume)
	if err := c.LoadSnapshot(ctx, &firego.SnapshotLoadParams{
		SnapshotPath: statePath,
		MemBackend: &firego.MemoryBackend{
			BackendType: firego.MemoryBackendFile,
			BackendPath: memPath,
		},
		ResumeVM: &resume,
	}); err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	log.Println("snapshot loaded successfully")
	return nil
}
