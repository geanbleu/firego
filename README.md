# firego

Go client for the [Firecracker](https://github.com/firecracker-microvm/firecracker) microVM HTTP API.

Firecracker exposes its management API over a Unix domain socket. This library wraps every endpoint with idiomatic Go types and methods, covering the full API surface.

## Installation

```bash
go get github.com/geanbleu/firego
```

Requires Go 1.22+.

## Quick start

```go
package main

import (
    "context"
    "log"

    "github.com/geanbleu/firego"
)

func main() {
    ctx := context.Background()
    c := firego.New("/run/firecracker.sock")

    c.PutBootSource(ctx, &firego.BootSource{
        KernelImagePath: "/var/lib/firecracker/vmlinux",
        BootArgs:        firego.Ptr("console=ttyS0 reboot=k panic=1 pci=off"),
    })

    c.PutMachineConfig(ctx, &firego.MachineConfiguration{
        VcpuCount:  2,
        MemSizeMib: 512,
    })

    c.PutDrive(ctx, "rootfs", &firego.Drive{
        DriveID:      "rootfs",
        IsRootDevice: true,
        PathOnHost:   firego.Ptr("/var/lib/firecracker/rootfs.ext4"),
    })

    c.PutNetworkInterface(ctx, "eth0", &firego.NetworkInterface{
        IfaceID:     "eth0",
        HostDevName: "tap0",
    })

    if err := c.StartInstance(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## API coverage

All Firecracker REST endpoints are supported. Methods are grouped by feature:

| Feature | Methods |
|---|---|
| **Instance** | `GetInstanceInfo`, `GetVersion`, `GetVMConfig` |
| **Actions** | `StartInstance`, `FlushMetrics`, `SendCtrlAltDel` |
| **Machine config** | `GetMachineConfig`, `PutMachineConfig`, `PatchMachineConfig` |
| **Boot source** | `PutBootSource` |
| **CPU config** | `PutCPUConfig` |
| **Drives** | `PutDrive`, `PatchDrive` |
| **Network interfaces** | `PutNetworkInterface`, `PatchNetworkInterface` |
| **Balloon** | `GetBalloon`, `PutBalloon`, `PatchBalloon`, `GetBalloonStats`, `PatchBalloonStats`, `StartBalloonHinting`, `GetBalloonHintingStatus`, `StopBalloonHinting` |
| **Snapshots** | `CreateSnapshot`, `LoadSnapshot` |
| **VM state** | `PatchVM`, `PauseVM`, `ResumeVM` |
| **MMDS** | `GetMMDS`, `PutMMDS`, `PatchMMDS`, `PutMMDSConfig` |
| **Vsock** | `PutVsock` |
| **Persistent memory** | `PutPmem`, `PatchPmem` |
| **Memory hotplug** | `PutHotplugMemory`, `PatchHotplugMemory`, `GetHotplugMemory` |
| **Logger** | `PutLogger` |
| **Metrics** | `PutMetrics` |
| **Entropy device** | `PutEntropyDevice` |
| **Serial console** | `PutSerial` |

### Pre-boot vs post-boot

Most configuration must be applied before `StartInstance` is called (pre-boot). The following operations are available while the VM is running (post-boot):

- Rate limiter updates — `PatchDrive`, `PatchNetworkInterface`, `PatchPmem`
- Balloon resizing — `PatchBalloon`, `PatchBalloonStats`
- Memory hotplug — `PatchHotplugMemory`
- Snapshot creation — `CreateSnapshot` (pause the VM first with `PauseVM`)
- VM state transitions — `PauseVM`, `ResumeVM`
- Metrics — `FlushMetrics`
- Guest signals — `SendCtrlAltDel`

## Error handling

API errors are returned as `*firego.APIError`, which carries the HTTP status code and the `fault_message` from Firecracker:

```go
if err := c.StartInstance(ctx); err != nil {
    var apiErr *firego.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("firecracker %d: %s\n", apiErr.StatusCode, apiErr.FaultMessage)
    }
}
```

## Optional fields

Many request structs have optional pointer fields. Use the `Ptr` helper to set them inline:

```go
&firego.BootSource{
    KernelImagePath: "/vmlinux",
    BootArgs:        firego.Ptr("console=ttyS0"),
}
```

## Snapshots

```go
// Pause → snapshot → resume
c.PauseVM(ctx)
c.CreateSnapshot(ctx, &firego.SnapshotCreateParams{
    MemFilePath:  "/snapshots/vm.mem",
    SnapshotPath: "/snapshots/vm.state",
})
c.ResumeVM(ctx)

// Restore on another instance
c.LoadSnapshot(ctx, &firego.SnapshotLoadParams{
    SnapshotPath: "/snapshots/vm.state",
    MemBackend: &firego.MemoryBackend{
        BackendType: firego.MemoryBackendFile,
        BackendPath: "/snapshots/vm.mem",
    },
    ResumeVM: firego.Ptr(true),
})
```

## Examples

Runnable example programs are in [`examples/`](./examples):

| Example | Description |
|---|---|
| [`examples/boot`](./examples/boot) | Boot a VM from kernel, rootfs, and TAP device |
| [`examples/snapshot`](./examples/snapshot) | Create and restore snapshots |
| [`examples/balloon`](./examples/balloon) | Configure and resize the balloon device |

```bash
go run ./examples/boot -kernel /vmlinux -rootfs /rootfs.ext4
go run ./examples/snapshot -action create -mem /tmp/vm.mem -state /tmp/vm.state
go run ./examples/balloon -target 256 -stats
```

## License

MIT
