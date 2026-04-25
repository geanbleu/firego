// Package firego provides a Go client for the Firecracker microVM HTTP API.
//
// Firecracker exposes its management API over a Unix domain socket using HTTP/JSON.
// This package wraps that API with idiomatic Go types and methods, covering all
// endpoints of the Firecracker REST API.
//
// # Creating a client
//
//	client := firego.New("/run/firecracker.sock")
//
// # VM lifecycle
//
// Operations fall into two categories:
//   - Pre-boot: configuration applied before [Client.StartInstance] is called.
//   - Post-boot: adjustments allowed while the VM is running (rate limiters, balloon, snapshots).
//
// A typical boot sequence:
//
//	ctx := context.Background()
//
//	client.PutBootSource(ctx, &firego.BootSource{
//	    KernelImagePath: "/path/to/vmlinux",
//	    BootArgs:        firego.Ptr("console=ttyS0 reboot=k panic=1 pci=off"),
//	})
//
//	client.PutMachineConfig(ctx, &firego.MachineConfiguration{
//	    VcpuCount:  2,
//	    MemSizeMib: 512,
//	})
//
//	client.PutDrive(ctx, "rootfs", &firego.Drive{
//	    DriveID:      "rootfs",
//	    IsRootDevice: true,
//	    PathOnHost:   firego.Ptr("/path/to/rootfs.ext4"),
//	})
//
//	client.PutNetworkInterface(ctx, "eth0", &firego.NetworkInterface{
//	    IfaceID:     "eth0",
//	    HostDevName: "tap0",
//	})
//
//	client.StartInstance(ctx)
//
// # Error handling
//
// All methods return an error. Firecracker API errors are surfaced as [*APIError],
// which carries the HTTP status code and the fault_message from the response body.
//
//	if err := client.StartInstance(ctx); err != nil {
//	    var apiErr *firego.APIError
//	    if errors.As(err, &apiErr) {
//	        fmt.Printf("firecracker %d: %s\n", apiErr.StatusCode, apiErr.FaultMessage)
//	    }
//	}
//
// # Optional fields
//
// Many struct fields are optional and represented as pointers. Use the [Ptr] helper
// to set them inline without declaring intermediate variables:
//
//	src := &firego.BootSource{
//	    KernelImagePath: "/vmlinux",
//	    BootArgs:        firego.Ptr("console=ttyS0"),
//	}
package firego
