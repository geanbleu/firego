// balloon demonstrates dynamic memory management via the virtio-balloon device:
// configuring the balloon pre-boot, querying statistics, and resizing it at runtime.
//
// Usage:
//
//	go run ./examples/balloon -socket /run/firecracker.sock -target 256
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/geanbleu/firego"
)

func main() {
	socket := flag.String("socket", "/run/firecracker.sock", "Firecracker Unix socket path")
	target := flag.Int("target", 0, "Target balloon size in MiB (0 = configure only, do not resize)")
	stats := flag.Bool("stats", false, "Print current balloon statistics")
	configure := flag.Bool("configure", false, "Configure balloon device (pre-boot)")
	flag.Parse()

	ctx := context.Background()
	c := firego.New(*socket)

	if *configure {
		if err := configureBalloon(ctx, c); err != nil {
			log.Fatal(err)
		}
	}

	if *target > 0 {
		if err := resize(ctx, c, *target); err != nil {
			log.Fatal(err)
		}
	}

	if *stats {
		if err := printStats(ctx, c); err != nil {
			log.Fatal(err)
		}
	}
}

func configureBalloon(ctx context.Context, c *firego.Client) error {
	log.Println("configuring balloon device (pre-boot)...")
	return c.PutBalloon(ctx, &firego.Balloon{
		AmountMib:             0,
		DeflateOnOom:          true,
		StatsPollingIntervalS: firego.Ptr(1),
	})
}

func resize(ctx context.Context, c *firego.Client, targetMib int) error {
	balloon, err := c.GetBalloon(ctx)
	if err != nil {
		return fmt.Errorf("get balloon: %w", err)
	}
	log.Printf("current balloon size: %d MiB → target: %d MiB", balloon.AmountMib, targetMib)

	if err := c.PatchBalloon(ctx, &firego.BalloonUpdate{AmountMib: targetMib}); err != nil {
		return fmt.Errorf("resize balloon: %w", err)
	}
	log.Printf("balloon resized to %d MiB", targetMib)
	return nil
}

func printStats(ctx context.Context, c *firego.Client) error {
	s, err := c.GetBalloonStats(ctx)
	if err != nil {
		return fmt.Errorf("get balloon stats: %w", err)
	}
	fmt.Printf("balloon statistics:\n")
	fmt.Printf("  target: %d MiB (%d pages)\n", s.TargetMib, s.TargetPages)
	fmt.Printf("  actual: %d MiB (%d pages)\n", s.ActualMib, s.ActualPages)
	if s.FreeMemory != nil {
		fmt.Printf("  free memory:      %d bytes\n", *s.FreeMemory)
	}
	if s.AvailableMemory != nil {
		fmt.Printf("  available memory: %d bytes\n", *s.AvailableMemory)
	}
	if s.TotalMemory != nil {
		fmt.Printf("  total memory:     %d bytes\n", *s.TotalMemory)
	}
	return nil
}
