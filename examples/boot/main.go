// boot demonstrates the minimal sequence required to start a Firecracker MicroVM.
//
// Usage:
//
//	go run ./examples/boot \
//	  -socket /run/firecracker.sock \
//	  -kernel /var/lib/firecracker/vmlinux \
//	  -rootfs  /var/lib/firecracker/rootfs.ext4 \
//	  -tap     tap0
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
	kernel := flag.String("kernel", "", "Path to the uncompressed kernel image (required)")
	rootfs := flag.String("rootfs", "", "Path to the root filesystem image (required)")
	tap := flag.String("tap", "tap0", "Host TAP device name for the network interface")
	vcpus := flag.Int("vcpus", 2, "Number of vCPUs")
	memMib := flag.Int("mem", 512, "Guest memory in MiB")
	flag.Parse()

	if *kernel == "" || *rootfs == "" {
		fmt.Fprintln(os.Stderr, "error: -kernel and -rootfs are required")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	c := firego.New(*socket)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"boot source", func() error {
			return c.PutBootSource(ctx, &firego.BootSource{
				KernelImagePath: *kernel,
				BootArgs:        firego.Ptr("console=ttyS0 reboot=k panic=1 pci=off"),
			})
		}},
		{"machine config", func() error {
			return c.PutMachineConfig(ctx, &firego.MachineConfiguration{
				VcpuCount:  *vcpus,
				MemSizeMib: *memMib,
			})
		}},
		{"root drive", func() error {
			return c.PutDrive(ctx, "rootfs", &firego.Drive{
				DriveID:      "rootfs",
				IsRootDevice: true,
				IsReadOnly:   firego.Ptr(false),
				PathOnHost:   rootfs,
			})
		}},
		{"network interface", func() error {
			return c.PutNetworkInterface(ctx, "eth0", &firego.NetworkInterface{
				IfaceID:     "eth0",
				HostDevName: *tap,
			})
		}},
		{"start", func() error {
			return c.StartInstance(ctx)
		}},
	}

	for _, s := range steps {
		log.Printf("configuring %s...", s.name)
		if err := s.fn(); err != nil {
			log.Fatalf("%s: %v", s.name, err)
		}
	}

	info, err := c.GetInstanceInfo(ctx)
	if err != nil {
		log.Fatalf("get instance info: %v", err)
	}
	log.Printf("VM started: id=%s state=%s vmm=%s", info.ID, info.State, info.VmmVersion)
}
