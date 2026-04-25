# firego

Go client for the [Firecracker](https://github.com/firecracker-microvm/firecracker) microVM HTTP API.

Firecracker exposes its management API over a Unix domain socket. This library wraps every endpoint with idiomatic Go types and methods, covering the full API surface.

## Installation

```bash
go get github.com/geanbleu/firego
```

Requires Go 1.22+.

## Table of contents

- [Quick start](#quick-start)
- [Client](#client)
- [Error handling](#error-handling)
- [Machine configuration](#machine-configuration)
- [Boot source](#boot-source)
- [CPU configuration](#cpu-configuration)
- [Drives](#drives)
- [Network interfaces](#network-interfaces)
- [VM state](#vm-state)
- [Snapshots](#snapshots)
- [Balloon device](#balloon-device)
- [MMDS](#mmds)
- [Vsock](#vsock)
- [Persistent memory (pmem)](#persistent-memory-pmem)
- [Memory hotplug](#memory-hotplug)
- [Logger](#logger)
- [Metrics](#metrics)
- [Entropy device](#entropy-device)
- [Serial console](#serial-console)
- [Instance info](#instance-info)
- [Examples](#examples)

---

## Quick start

```go
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
```

---

## Client

Create a client by passing the path to the Firecracker Unix socket:

```go
c := firego.New("/run/firecracker.sock")
```

The socket does not need to exist at construction time; the connection is established on each API call. All methods accept a `context.Context` for cancellation and deadline control.

### Optional fields

Many request structs have optional pointer fields. Use the `Ptr` helper to set them inline without intermediate variables:

```go
&firego.BootSource{
    KernelImagePath: "/vmlinux",
    BootArgs:        firego.Ptr("console=ttyS0"),
}
```

---

## Error handling

All methods return an `error`. Firecracker API errors are returned as `*APIError`, which carries the HTTP status code and the `fault_message` from the response body.

```go
if err := c.StartInstance(ctx); err != nil {
    var apiErr *firego.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.FaultMessage)
    }
    return err
}
```

---

## Machine configuration

Defines the vCPU count, memory size, and hardware options. Must be applied **pre-boot**.

```go
c.PutMachineConfig(ctx, &firego.MachineConfiguration{
    VcpuCount:       2,       // 1 or an even number, up to 32
    MemSizeMib:      1024,
    Smt:             firego.Ptr(false),         // simultaneous multi-threading (x86_64 only)
    TrackDirtyPages: firego.Ptr(true),          // required for diff snapshots
    HugePages:       firego.Ptr(firego.HugePages2M), // None (default) or 2M
})
```

Partially update individual fields without replacing the whole config:

```go
c.PatchMachineConfig(ctx, &firego.MachineConfiguration{
    TrackDirtyPages: firego.Ptr(true),
})
```

Read the current configuration:

```go
cfg, err := c.GetMachineConfig(ctx)
```

---

## Boot source

Configures the guest kernel and optional initrd. Must be applied **pre-boot**.

```go
c.PutBootSource(ctx, &firego.BootSource{
    KernelImagePath: "/var/lib/firecracker/vmlinux", // uncompressed kernel binary
    BootArgs:        firego.Ptr("console=ttyS0 reboot=k panic=1 pci=off"),
    InitrdPath:      firego.Ptr("/var/lib/firecracker/initrd.img"), // optional
})
```

---

## CPU configuration

Sets fine-grained CPU feature flags per vCPU. This is the preferred alternative to the deprecated `CpuTemplate` field. Must be applied **pre-boot**.

**x86_64 — CPUID modifiers:**

```go
c.PutCPUConfig(ctx, &firego.CpuConfig{
    CpuidModifiers: []firego.CpuidLeafModifier{
        {
            Leaf:    "0x1",
            Subleaf: "0x0",
            Flags:   0,
            Modifiers: []firego.CpuidRegisterModifier{
                {
                    Register: firego.CpuidRegisterEcx,
                    // Each character is a bit: '0' clears, '1' sets, 'x' leaves unchanged.
                    Bitmap: "0bxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
                },
            },
        },
    },
})
```

**x86_64 — MSR modifiers:**

```go
c.PutCPUConfig(ctx, &firego.CpuConfig{
    MsrModifiers: []firego.MsrModifier{
        {Addr: "0x10a", Bitmap: "0b" + strings.Repeat("x", 64)},
    },
})
```

**aarch64 — register and feature modifiers:**

```go
c.PutCPUConfig(ctx, &firego.CpuConfig{
    RegModifiers: []firego.ArmRegisterModifier{
        {Addr: "0x603000000013c064", Bitmap: "0b" + strings.Repeat("x", 64)},
    },
    VcpuFeatures: []firego.VcpuFeatures{
        {Index: 0, Bitmap: "0b" + strings.Repeat("x", 32)},
    },
})
```

---

## Drives

Attaches block devices to the guest. At least one drive with `IsRootDevice: true` is required to boot.

### virtio-block (file or block device)

```go
// Pre-boot: attach root drive
c.PutDrive(ctx, "rootfs", &firego.Drive{
    DriveID:      "rootfs",
    IsRootDevice: true,
    PathOnHost:   firego.Ptr("/var/lib/firecracker/rootfs.ext4"),
    IsReadOnly:   firego.Ptr(false),
    CacheType:    firego.Ptr(firego.DriveCacheUnsafe), // Unsafe (default) or Writeback
    IoEngine:     firego.Ptr(firego.DriveIoEngineSync), // Sync (default) or Async
})

// Pre-boot: attach a secondary read-only data drive
c.PutDrive(ctx, "data", &firego.Drive{
    DriveID:      "data",
    IsRootDevice: false,
    PathOnHost:   firego.Ptr("/var/lib/firecracker/data.ext4"),
    IsReadOnly:   firego.Ptr(true),
})
```

### vhost-user-block

```go
c.PutDrive(ctx, "vhost-data", &firego.Drive{
    DriveID:      "vhost-data",
    IsRootDevice: false,
    Socket:       firego.Ptr("/run/vhost-user-blk.sock"),
})
```

### Rate limiting (post-boot)

```go
c.PatchDrive(ctx, "data", &firego.PartialDrive{
    DriveID: "data",
    RateLimiter: &firego.RateLimiter{
        Bandwidth: &firego.TokenBucket{
            Size:         104857600, // 100 MiB bucket
            RefillTime:   1000,      // refilled every 1 s
            OneTimeBurst: firego.Ptr(int64(10485760)), // 10 MiB initial burst
        },
        Ops: &firego.TokenBucket{
            Size:       1000,
            RefillTime: 1000,
        },
    },
})
```

---

## Network interfaces

Attaches virtio-net devices backed by host TAP interfaces.

### Attach an interface (pre-boot)

```go
c.PutNetworkInterface(ctx, "eth0", &firego.NetworkInterface{
    IfaceID:     "eth0",
    HostDevName: "tap0",
    GuestMAC:    firego.Ptr("AA:BB:CC:DD:EE:FF"), // optional; auto-generated if omitted
})
```

### Rate limiting (post-boot)

```go
c.PatchNetworkInterface(ctx, "eth0", &firego.PartialNetworkInterface{
    IfaceID: "eth0",
    TxRateLimiter: &firego.RateLimiter{
        Bandwidth: &firego.TokenBucket{
            Size:       125000000, // 125 MB/s (~1 Gbit/s)
            RefillTime: 1000,
        },
    },
})
```

---

## VM state

Pause and resume all vCPUs while the VM remains in memory.

```go
c.PauseVM(ctx)  // suspend execution
c.ResumeVM(ctx) // resume execution
```

These are convenience wrappers for `PatchVM`, which accepts any `VmState` value:

```go
c.PatchVM(ctx, &firego.Vm{State: firego.VmStatePaused})
```

---

## Snapshots

Snapshots capture the full VM state (CPU registers, memory, device state) to disk.
`TrackDirtyPages` must be enabled in `MachineConfiguration` to use diff snapshots.

### Create a snapshot (post-boot)

```go
// Always pause the VM first for a consistent snapshot.
c.PauseVM(ctx)

c.CreateSnapshot(ctx, &firego.SnapshotCreateParams{
    MemFilePath:  "/snapshots/vm.mem",
    SnapshotPath: "/snapshots/vm.state",
    SnapshotType: firego.Ptr(firego.SnapshotTypeFull), // Full (default) or Diff
})

c.ResumeVM(ctx)
```

### Load a snapshot (pre-boot)

```go
c.LoadSnapshot(ctx, &firego.SnapshotLoadParams{
    SnapshotPath: "/snapshots/vm.state",

    // Preferred: configure the memory loading backend explicitly.
    MemBackend: &firego.MemoryBackend{
        BackendType: firego.MemoryBackendFile, // File or Uffd (userfaultfd)
        BackendPath: "/snapshots/vm.mem",
    },

    ResumeVM:        firego.Ptr(true),  // start immediately after load
    TrackDirtyPages: firego.Ptr(true),  // enable for subsequent diff snapshots

    // Override TAP and vsock paths when restoring on a different host.
    NetworkOverrides: []firego.NetworkOverride{
        {IfaceID: "eth0", HostDevName: "tap1"},
    },
    VsockOverride: &firego.VsockOverride{UDSPath: "/run/vsock-new.sock"},
})
```

---

## Balloon device

The balloon device dynamically adjusts the amount of memory available to the guest by inflating (reclaiming from guest) or deflating (returning to guest).

### Configure (pre-boot)

```go
c.PutBalloon(ctx, &firego.Balloon{
    AmountMib:             0,    // initial size; 0 means fully deflated
    DeflateOnOom:          true, // automatically deflate under memory pressure
    StatsPollingIntervalS: firego.Ptr(1), // collect stats every second (0 = disabled)
})
```

### Resize at runtime (post-boot)

```go
// Inflate: reclaim 256 MiB from the guest.
c.PatchBalloon(ctx, &firego.BalloonUpdate{AmountMib: 256})

// Deflate: return all memory to the guest.
c.PatchBalloon(ctx, &firego.BalloonUpdate{AmountMib: 0})
```

### Statistics

```go
stats, err := c.GetBalloonStats(ctx)
fmt.Printf("guest free memory: %d bytes\n", *stats.FreeMemory)
fmt.Printf("balloon actual:    %d MiB\n", stats.ActualMib)

// Adjust polling interval at runtime.
c.PatchBalloonStats(ctx, &firego.BalloonStatsUpdate{StatsPollingIntervalS: 5})
```

### Free page hinting

Free page hinting allows the host to reclaim guest memory pages that the guest OS has marked as free, without inflating the balloon.

```go
c.StartBalloonHinting(ctx, &firego.BalloonStartCmd{
    AcknowledgeOnStop: firego.Ptr(true),
})

status, _ := c.GetBalloonHintingStatus(ctx)
fmt.Printf("host cmd: %d, guest cmd: %v\n", status.HostCmd, status.GuestCmd)

c.StopBalloonHinting(ctx)
```

---

## MMDS

The Microvm Metadata Service (MMDS) is a key-value store served to the guest at a link-local IP address (default `169.254.169.254`), similar to the EC2 Instance Metadata Service.

### Configure (pre-boot)

```go
c.PutMMDSConfig(ctx, &firego.MmdsConfig{
    // Interfaces through which the guest can reach the MMDS endpoint.
    NetworkInterfaces: []string{"eth0"},
    Version:     firego.Ptr(firego.MmdsVersionV2), // V1 (default) or V2
    IPv4Address: firego.Ptr("169.254.169.254"),     // default
    ImdsCompat:  firego.Ptr(false),                 // EC2 IMDS compatibility
})
```

### Populate the data store

```go
// Create or fully replace the data store.
c.PutMMDS(ctx, firego.MmdsContentsObject{
    "latest": map[string]interface{}{
        "meta-data": map[string]interface{}{
            "instance-id": "i-1234567890",
            "hostname":    "my-microvm",
        },
    },
})

// Merge-update: existing keys not present in the patch are preserved.
c.PatchMMDS(ctx, firego.MmdsContentsObject{
    "latest": map[string]interface{}{
        "meta-data": map[string]interface{}{
            "hostname": "updated-hostname",
        },
    },
})

// Read back the full store.
data, err := c.GetMMDS(ctx)
```

---

## Vsock

The vsock device provides a communication channel between the host and the guest using the VM Sockets (AF_VSOCK) protocol, proxied through a Unix domain socket on the host.

```go
// Pre-boot: attach vsock device.
c.PutVsock(ctx, &firego.Vsock{
    GuestCID: 3,                        // minimum 3; 0–2 are reserved
    UDSPath:  "/run/firecracker-vsock.sock",
})
```

The host-side socket at `UDSPath` must be created by the caller before `StartInstance` is called. Firecracker connects to it at boot time.

---

## Persistent memory (pmem)

The virtio-pmem device exposes a host file as persistent memory (NVDIMM) to the guest.

```go
// Pre-boot: attach a pmem device.
c.PutPmem(ctx, "pmem0", &firego.Pmem{
    ID:         "pmem0",
    PathOnHost: "/var/lib/firecracker/pmem.img",
    RootDevice: firego.Ptr(false),
    ReadOnly:   firego.Ptr(false),
})

// Post-boot: update the rate limiter.
c.PatchPmem(ctx, "pmem0", &firego.PartialPmem{
    ID: "pmem0",
    RateLimiter: &firego.RateLimiter{
        Bandwidth: &firego.TokenBucket{Size: 52428800, RefillTime: 1000},
    },
})
```

---

## Memory hotplug

The virtio-mem device allows memory to be added to a running guest without restarting it.
`TotalSizeMib` defines the upper bound; actual plugging happens post-boot via `PatchHotplugMemory`.

### Configure (pre-boot)

```go
c.PutHotplugMemory(ctx, &firego.MemoryHotplugConfig{
    TotalSizeMib: firego.Ptr(4096), // maximum hotpluggable memory
    SlotSizeMib:  firego.Ptr(128),  // granularity (default and minimum: 128 MiB)
    BlockSizeMib: firego.Ptr(2),    // logical block size (default and minimum: 2 MiB)
})
```

### Plug and unplug memory (post-boot)

```go
// Plug 512 MiB into the guest.
c.PatchHotplugMemory(ctx, &firego.MemoryHotplugSizeUpdate{
    RequestedSizeMib: firego.Ptr(512),
})

// Unplug all hotplugged memory.
c.PatchHotplugMemory(ctx, &firego.MemoryHotplugSizeUpdate{
    RequestedSizeMib: firego.Ptr(0),
})

// Check current status.
status, err := c.GetHotplugMemory(ctx)
fmt.Printf("plugged: %d MiB / total: %d MiB\n", *status.PluggedSizeMib, *status.TotalSizeMib)
```

---

## Logger

Initializes the Firecracker logging subsystem. The target file or named pipe must exist before this call is made.

```go
c.PutLogger(ctx, &firego.Logger{
    Level:         firego.Ptr(firego.LogLevelDebug), // Error|Warning|Info|Debug|Trace|Off
    LogPath:       firego.Ptr("/run/firecracker.log"),
    ShowLevel:     firego.Ptr(true),  // prefix each line with the log level
    ShowLogOrigin: firego.Ptr(false), // include source file and line number
    Module:        firego.Ptr(""),    // filter to a specific Rust module path
})
```

---

## Metrics

Initializes the Firecracker metrics subsystem. Metrics are written as newline-delimited JSON objects to the target file or named pipe.

```go
c.PutMetrics(ctx, &firego.Metrics{
    MetricsPath: "/run/firecracker-metrics.log",
})

// Force an immediate flush outside of the normal periodic schedule.
c.FlushMetrics(ctx)
```

---

## Entropy device

The virtio-rng device provides hardware-backed entropy to the guest, readable via `/dev/hwrng`.

```go
c.PutEntropyDevice(ctx, &firego.EntropyDevice{
    RateLimiter: &firego.RateLimiter{
        Bandwidth: &firego.TokenBucket{
            Size:       1024,
            RefillTime: 1000,
        },
    },
})
```

---

## Serial console

Redirects guest serial output to a host file or named pipe.

```go
c.PutSerial(ctx, &firego.SerialDevice{
    SerialOutPath: firego.Ptr("/run/firecracker-serial.log"),
    RateLimiter: &firego.TokenBucket{
        Size:       1048576, // 1 MiB bucket
        RefillTime: 1000,
    },
})
```

---

## Instance info

```go
// Retrieve instance identity and current state.
info, err := c.GetInstanceInfo(ctx)
fmt.Printf("id=%s state=%s vmm=%s\n", info.ID, info.State, info.VmmVersion)

// Retrieve the Firecracker build version.
v, err := c.GetVersion(ctx)
fmt.Println(v.FirecrackerVersion)

// Retrieve a snapshot of the complete VM configuration.
cfg, err := c.GetVMConfig(ctx)
```

---

## Examples

Runnable example programs are in [`examples/`](./examples):

| Example | Description |
|---|---|
| [`examples/boot`](./examples/boot) | Boot a VM from a kernel, rootfs image, and TAP device |
| [`examples/snapshot`](./examples/snapshot) | Create and restore snapshots |
| [`examples/balloon`](./examples/balloon) | Configure and resize the balloon device at runtime |

```bash
go run ./examples/boot     -kernel /vmlinux -rootfs /rootfs.ext4 -tap tap0
go run ./examples/snapshot -action create -mem /tmp/vm.mem -state /tmp/vm.state
go run ./examples/snapshot -action load   -mem /tmp/vm.mem -state /tmp/vm.state -resume
go run ./examples/balloon  -configure -target 256 -stats
```

---

## License

MIT
