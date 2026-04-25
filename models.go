package firego

// ─── Rate limiting ────────────────────────────────────────────────────────────

// TokenBucket defines a token-bucket rate limiter used to cap bandwidth or IOPS.
// Size and RefillTime must both be ≥ 0.
type TokenBucket struct {
	// Size is the maximum number of tokens (bytes or ops) the bucket can hold.
	Size int64 `json:"size"`
	// RefillTime is the interval in milliseconds over which Size tokens are replenished.
	RefillTime int64 `json:"refill_time"`
	// OneTimeBurst is an optional initial burst allowance added on top of Size.
	OneTimeBurst *int64 `json:"one_time_burst,omitempty"`
}

// RateLimiter pairs independent token buckets for bandwidth (bytes/s) and ops/s limiting.
// Either or both fields may be set; omitting a field disables that dimension of limiting.
type RateLimiter struct {
	Bandwidth *TokenBucket `json:"bandwidth,omitempty"`
	Ops       *TokenBucket `json:"ops,omitempty"`
}

// ─── Instance ────────────────────────────────────────────────────────────────

// InstanceState enumerates the possible running states of a MicroVM.
type InstanceState string

const (
	InstanceStateNotStarted InstanceState = "Not started"
	InstanceStateRunning    InstanceState = "Running"
	InstanceStatePaused     InstanceState = "Paused"
)

// InstanceInfo describes a running MicroVM instance (GET /).
type InstanceInfo struct {
	// AppName is the name of the Firecracker application binary.
	AppName string `json:"app_name"`
	// ID is the unique identifier of this MicroVM instance.
	ID string `json:"id"`
	// State is the current running state of the instance.
	State InstanceState `json:"state"`
	// VmmVersion is the Firecracker build version string.
	VmmVersion string `json:"vmm_version"`
}

// FirecrackerVersion holds the build version (GET /version).
type FirecrackerVersion struct {
	FirecrackerVersion string `json:"firecracker_version"`
}

// ─── Actions ─────────────────────────────────────────────────────────────────

// ActionType enumerates the synchronous actions available via PUT /actions.
type ActionType string

const (
	// ActionFlushMetrics forces an immediate flush of in-memory metrics to disk.
	ActionFlushMetrics ActionType = "FlushMetrics"
	// ActionInstanceStart boots the MicroVM. All pre-boot configuration must be
	// complete before calling this action.
	ActionInstanceStart ActionType = "InstanceStart"
	// ActionSendCtrlAltDel sends the Ctrl+Alt+Del key sequence to the guest,
	// typically triggering a graceful reboot or shutdown.
	ActionSendCtrlAltDel ActionType = "SendCtrlAltDel"
)

// InstanceActionInfo wraps an action type for PUT /actions.
type InstanceActionInfo struct {
	ActionType ActionType `json:"action_type"`
}

// ─── Machine configuration ───────────────────────────────────────────────────

// CpuTemplate is a deprecated pre-defined CPU template. Prefer [CpuConfig] instead.
type CpuTemplate string

const (
	CpuTemplateC3   CpuTemplate = "C3"
	CpuTemplateT2   CpuTemplate = "T2"
	CpuTemplateT2S  CpuTemplate = "T2S"
	CpuTemplateT2CL CpuTemplate = "T2CL"
	CpuTemplateT2A  CpuTemplate = "T2A"
	CpuTemplateV1N1 CpuTemplate = "V1N1"
	CpuTemplateNone CpuTemplate = "None"
)

// HugePages enumerates the huge page backing options for VM memory.
type HugePages string

const (
	HugePagesNone HugePages = "None"
	HugePages2M   HugePages = "2M"
)

// MachineConfiguration defines vCPU count, memory size, and related hardware options.
// All fields must be set before calling [Client.StartInstance] (pre-boot).
type MachineConfiguration struct {
	// MemSizeMib is the amount of memory available to the guest, in mebibytes.
	MemSizeMib int `json:"mem_size_mib"`
	// VcpuCount is the number of virtual CPUs (must be 1 or an even number, 1–32).
	VcpuCount int `json:"vcpu_count"`
	// CpuTemplate selects a deprecated built-in CPU template. Prefer CpuConfig.
	CpuTemplate *CpuTemplate `json:"cpu_template,omitempty"`
	// Smt enables Simultaneous Multi-Threading (x86_64 only, default false).
	Smt *bool `json:"smt,omitempty"`
	// TrackDirtyPages enables dirty-page tracking, required for diff snapshots.
	TrackDirtyPages *bool `json:"track_dirty_pages,omitempty"`
	// HugePages selects the huge-page size used to back guest memory.
	HugePages *HugePages `json:"huge_pages,omitempty"`
}

// ─── CPU configuration ───────────────────────────────────────────────────────

// CpuidRegister enumerates the CPUID registers that can be modified (x86_64).
type CpuidRegister string

const (
	CpuidRegisterEax CpuidRegister = "eax"
	CpuidRegisterEbx CpuidRegister = "ebx"
	CpuidRegisterEcx CpuidRegister = "ecx"
	CpuidRegisterEdx CpuidRegister = "edx"
)

// CpuidRegisterModifier modifies a single CPUID register within a leaf (x86_64).
type CpuidRegisterModifier struct {
	// Register identifies which CPUID output register to modify.
	Register CpuidRegister `json:"register"`
	// Bitmap is a 32-character binary string prefixed with "0b" describing the
	// bits to set, clear, or leave unchanged.
	Bitmap string `json:"bitmap"`
}

// CpuidLeafModifier modifies a CPUID leaf/subleaf combination (x86_64).
type CpuidLeafModifier struct {
	// Leaf is the CPUID leaf index (hexadecimal, binary, or decimal string).
	Leaf string `json:"leaf"`
	// Subleaf is the CPUID subleaf index (same format as Leaf).
	Subleaf string `json:"subleaf"`
	// Flags holds the KVM feature flags for this leaf.
	Flags int32 `json:"flags"`
	// Modifiers lists per-register modifications to apply.
	Modifiers []CpuidRegisterModifier `json:"modifiers"`
}

// MsrModifier modifies a Model-Specific Register (x86_64).
type MsrModifier struct {
	// Addr is the MSR address (hexadecimal, binary, or decimal string).
	Addr string `json:"addr"`
	// Bitmap is a 64-character binary string describing the bits to modify.
	Bitmap string `json:"bitmap"`
}

// ArmRegisterModifier modifies an ARM system register (aarch64).
type ArmRegisterModifier struct {
	// Addr is the 64-bit register address.
	Addr string `json:"addr"`
	// Bitmap is a binary string of up to 128 bits describing the bits to modify.
	Bitmap string `json:"bitmap"`
}

// VcpuFeatures modifies a vCPU feature flag (aarch64).
type VcpuFeatures struct {
	// Index is the position in the kvm_vcpu_init.features array.
	Index int32 `json:"index"`
	// Bitmap is a 32-character binary string describing the bits to modify.
	Bitmap string `json:"bitmap"`
}

// CpuConfig defines per-vCPU feature flag overrides (pre-boot, PUT /cpu-config).
// Fields are architecture-specific: use CpuidModifiers and MsrModifiers on x86_64,
// RegModifiers and VcpuFeatures on aarch64.
type CpuConfig struct {
	KvmCapabilities []string              `json:"kvm_capabilities,omitempty"`
	CpuidModifiers  []CpuidLeafModifier   `json:"cpuid_modifiers,omitempty"`
	MsrModifiers    []MsrModifier         `json:"msr_modifiers,omitempty"`
	RegModifiers    []ArmRegisterModifier  `json:"reg_modifiers,omitempty"`
	VcpuFeatures    []VcpuFeatures        `json:"vcpu_features,omitempty"`
}

// ─── Boot source ─────────────────────────────────────────────────────────────

// BootSource configures the kernel and optional initrd used to boot the guest
// (pre-boot, PUT /boot-source).
type BootSource struct {
	// KernelImagePath is the host filesystem path to the uncompressed kernel binary.
	KernelImagePath string `json:"kernel_image_path"`
	// BootArgs are the kernel command-line parameters passed at boot.
	BootArgs *string `json:"boot_args,omitempty"`
	// InitrdPath is the optional host path to an initrd/initramfs image.
	InitrdPath *string `json:"initrd_path,omitempty"`
}

// ─── Drives ──────────────────────────────────────────────────────────────────

// DriveCacheType controls the virtio-blk cache mode.
type DriveCacheType string

const (
	DriveCacheUnsafe    DriveCacheType = "Unsafe"
	DriveCacheWriteback DriveCacheType = "Writeback"
)

// DriveIoEngine selects the I/O engine used by the virtio-blk device.
type DriveIoEngine string

const (
	DriveIoEngineSync  DriveIoEngine = "Sync"
	DriveIoEngineAsync DriveIoEngine = "Async"
)

// Drive describes a guest block device (pre-boot, PUT /drives/{drive_id}).
// Either PathOnHost (virtio-block) or Socket (vhost-user-block) must be set.
type Drive struct {
	// DriveID uniquely identifies the drive. It is also used as the URL path segment.
	DriveID string `json:"drive_id"`
	// IsRootDevice marks this drive as the root device (passed as root= to the kernel).
	IsRootDevice bool `json:"is_root_device"`
	// PartUUID is the optional unique partition UUID used for root device identification.
	PartUUID *string `json:"partuuid,omitempty"`
	// CacheType selects the virtio-blk cache mode (default Unsafe).
	CacheType *DriveCacheType `json:"cache_type,omitempty"`
	// IsReadOnly mounts the drive in read-only mode (virtio-block only).
	IsReadOnly *bool `json:"is_read_only,omitempty"`
	// PathOnHost is the host filesystem path to the backing file or block device (virtio-block).
	PathOnHost *string `json:"path_on_host,omitempty"`
	// RateLimiter limits the I/O bandwidth and operations per second for this drive.
	RateLimiter *RateLimiter `json:"rate_limiter,omitempty"`
	// IoEngine selects between synchronous and asynchronous I/O (virtio-block, default Sync).
	IoEngine *DriveIoEngine `json:"io_engine,omitempty"`
	// Socket is the host path to the vhost-user-block Unix domain socket.
	Socket *string `json:"socket,omitempty"`
}

// PartialDrive updates a drive's path or rate limiter post-boot
// (PATCH /drives/{drive_id}).
type PartialDrive struct {
	DriveID     string       `json:"drive_id"`
	PathOnHost  *string      `json:"path_on_host,omitempty"`
	RateLimiter *RateLimiter `json:"rate_limiter,omitempty"`
}

// ─── Persistent memory (pmem) ────────────────────────────────────────────────

// Pmem describes a virtio-pmem persistent memory device
// (pre-boot, PUT /pmem/{id}).
type Pmem struct {
	// ID uniquely identifies the pmem device.
	ID string `json:"id"`
	// PathOnHost is the host path to the backing file for the persistent memory region.
	PathOnHost string `json:"path_on_host"`
	// RootDevice marks this pmem device as the root device.
	RootDevice *bool `json:"root_device,omitempty"`
	// ReadOnly maps the memory region as read-only inside the guest.
	ReadOnly *bool `json:"read_only,omitempty"`
	// RateLimiter limits I/O on this pmem device.
	RateLimiter *RateLimiter `json:"rate_limiter,omitempty"`
}

// PartialPmem updates the rate limiter of a pmem device post-boot
// (PATCH /pmem/{id}).
type PartialPmem struct {
	ID          string       `json:"id"`
	RateLimiter *RateLimiter `json:"rate_limiter,omitempty"`
}

// ─── Network interfaces ──────────────────────────────────────────────────────

// NetworkInterface defines a guest virtio-net device backed by a host TAP interface
// (pre-boot, PUT /network-interfaces/{iface_id}).
type NetworkInterface struct {
	// IfaceID uniquely identifies the interface. It is also the URL path segment.
	IfaceID string `json:"iface_id"`
	// HostDevName is the name of the host TAP device to attach to this interface.
	HostDevName string `json:"host_dev_name"`
	// GuestMAC sets the MAC address visible inside the guest.
	// If omitted, Firecracker generates one automatically.
	GuestMAC *string `json:"guest_mac,omitempty"`
	// RxRateLimiter limits inbound (guest-receive) traffic.
	RxRateLimiter *RateLimiter `json:"rx_rate_limiter,omitempty"`
	// TxRateLimiter limits outbound (guest-transmit) traffic.
	TxRateLimiter *RateLimiter `json:"tx_rate_limiter,omitempty"`
}

// PartialNetworkInterface updates rate limiters on a network interface post-boot
// (PATCH /network-interfaces/{iface_id}).
type PartialNetworkInterface struct {
	IfaceID       string       `json:"iface_id"`
	RxRateLimiter *RateLimiter `json:"rx_rate_limiter,omitempty"`
	TxRateLimiter *RateLimiter `json:"tx_rate_limiter,omitempty"`
}

// ─── Balloon device ──────────────────────────────────────────────────────────

// Balloon describes a virtio-balloon device (pre-boot, PUT /balloon).
type Balloon struct {
	// AmountMib is the target size of the balloon, in mebibytes.
	AmountMib int `json:"amount_mib"`
	// DeflateOnOom automatically deflates the balloon when the guest is under
	// memory pressure (OOM condition).
	DeflateOnOom bool `json:"deflate_on_oom"`
	// StatsPollingIntervalS sets how often the balloon device reports statistics,
	// in seconds. 0 disables statistics collection.
	StatsPollingIntervalS *int `json:"stats_polling_interval_s,omitempty"`
	// FreePageHinting enables the free page hinting feature.
	FreePageHinting *bool `json:"free_page_hinting,omitempty"`
	// FreePageReporting enables the free page reporting feature.
	FreePageReporting *bool `json:"free_page_reporting,omitempty"`
}

// BalloonUpdate changes the target size of the balloon device post-boot
// (PATCH /balloon).
type BalloonUpdate struct {
	// AmountMib is the new target balloon size, in mebibytes.
	AmountMib int `json:"amount_mib"`
}

// BalloonStats holds the latest statistics reported by the balloon device
// (GET /balloon/statistics).
type BalloonStats struct {
	TargetPages     int    `json:"target_pages"`
	ActualPages     int    `json:"actual_pages"`
	TargetMib       int    `json:"target_mib"`
	ActualMib       int    `json:"actual_mib"`
	SwapIn          *int64 `json:"swap_in,omitempty"`
	SwapOut         *int64 `json:"swap_out,omitempty"`
	MajorFaults     *int64 `json:"major_faults,omitempty"`
	MinorFaults     *int64 `json:"minor_faults,omitempty"`
	FreeMemory      *int64 `json:"free_memory,omitempty"`
	TotalMemory     *int64 `json:"total_memory,omitempty"`
	AvailableMemory *int64 `json:"available_memory,omitempty"`
	DiskCaches      *int64 `json:"disk_caches,omitempty"`
	HugetlbAllocs   *int64 `json:"hugetlb_allocations,omitempty"`
	HugetlbFailures *int64 `json:"hugetlb_failures,omitempty"`
	OomKill         *int64 `json:"oom_kill,omitempty"`
	AllocStall      *int64 `json:"alloc_stall,omitempty"`
	AsyncScan       *int64 `json:"async_scan,omitempty"`
	DirectScan      *int64 `json:"direct_scan,omitempty"`
	AsyncReclaim    *int64 `json:"async_reclaim,omitempty"`
	DirectReclaim   *int64 `json:"direct_reclaim,omitempty"`
}

// BalloonStatsUpdate changes the statistics polling interval post-boot
// (PATCH /balloon/statistics).
type BalloonStatsUpdate struct {
	// StatsPollingIntervalS is the new polling interval in seconds.
	StatsPollingIntervalS int `json:"stats_polling_interval_s"`
}

// BalloonStartCmd initiates a free page hinting run (PATCH /balloon/hinting/start).
type BalloonStartCmd struct {
	// AcknowledgeOnStop automatically acknowledges the stop command sent by the guest.
	AcknowledgeOnStop *bool `json:"acknowledge_on_stop,omitempty"`
}

// BalloonHintingStatus describes the current state of a free page hinting run
// (GET /balloon/hinting/status).
type BalloonHintingStatus struct {
	// HostCmd is the last command issued by the host.
	HostCmd int `json:"host_cmd"`
	// GuestCmd is the last command received from the guest, if any.
	GuestCmd *int `json:"guest_cmd,omitempty"`
}

// ─── Snapshots ───────────────────────────────────────────────────────────────

// SnapshotType selects between a full and an incremental (diff) snapshot.
type SnapshotType string

const (
	SnapshotTypeFull SnapshotType = "Full"
	SnapshotTypeDiff SnapshotType = "Diff"
)

// MemoryBackendType selects how guest memory is loaded from a snapshot.
type MemoryBackendType string

const (
	// MemoryBackendFile loads memory directly from a regular file.
	MemoryBackendFile MemoryBackendType = "File"
	// MemoryBackendUffd loads memory on demand via a userfaultfd handler.
	MemoryBackendUffd MemoryBackendType = "Uffd"
)

// MemoryBackend configures how guest memory is loaded when restoring a snapshot.
type MemoryBackend struct {
	// BackendType selects between file-based and userfaultfd-based loading.
	BackendType MemoryBackendType `json:"backend_type"`
	// BackendPath is either the path to the memory file (File) or the path to
	// the userfaultfd Unix domain socket (Uffd).
	BackendPath string `json:"backend_path"`
}

// NetworkOverride replaces the backing TAP device of a network interface when
// restoring a snapshot. Useful when the host TAP device name has changed.
type NetworkOverride struct {
	// IfaceID identifies the network interface to override.
	IfaceID string `json:"iface_id"`
	// HostDevName is the new host TAP device name to attach to this interface.
	HostDevName string `json:"host_dev_name"`
}

// VsockOverride replaces the backing Unix domain socket of the vsock device
// when restoring a snapshot.
type VsockOverride struct {
	// UDSPath is the new host path for the vsock proxy socket.
	UDSPath string `json:"uds_path"`
}

// SnapshotCreateParams configures snapshot creation (post-boot, PUT /snapshot/create).
// The VM should be paused with [Client.PauseVM] before creating a snapshot.
type SnapshotCreateParams struct {
	// MemFilePath is the host path where guest memory will be saved.
	MemFilePath string `json:"mem_file_path"`
	// SnapshotPath is the host path where the MicroVM state file will be written.
	SnapshotPath string `json:"snapshot_path"`
	// SnapshotType selects a full or incremental diff snapshot (default Full).
	// Diff snapshots require TrackDirtyPages to be enabled in [MachineConfiguration].
	SnapshotType *SnapshotType `json:"snapshot_type,omitempty"`
}

// SnapshotLoadParams configures snapshot restoration (pre-boot, PUT /snapshot/load).
// Exactly one of MemFilePath or MemBackend must be set.
type SnapshotLoadParams struct {
	// SnapshotPath is the host path to the MicroVM state file.
	SnapshotPath string `json:"snapshot_path"`
	// EnableDiffSnapshots is deprecated; use TrackDirtyPages instead.
	EnableDiffSnapshots *bool `json:"enable_diff_snapshots,omitempty"`
	// TrackDirtyPages enables dirty-page tracking after restore (required for diff snapshots).
	TrackDirtyPages *bool `json:"track_dirty_pages,omitempty"`
	// MemFilePath is the host path to the guest memory file (deprecated; use MemBackend).
	MemFilePath *string `json:"mem_file_path,omitempty"`
	// MemBackend configures the memory loading backend (preferred over MemFilePath).
	MemBackend *MemoryBackend `json:"mem_backend,omitempty"`
	// ResumeVM automatically resumes the VM after a successful snapshot load.
	ResumeVM *bool `json:"resume_vm,omitempty"`
	// NetworkOverrides replace TAP device names for any network interfaces.
	NetworkOverrides []NetworkOverride `json:"network_overrides,omitempty"`
	// VsockOverride replaces the vsock backing socket path.
	VsockOverride *VsockOverride `json:"vsock_override,omitempty"`
	// ClockRealtime advances the guest kvmclock by the elapsed wall-clock time
	// since the snapshot was taken (x86_64 only).
	ClockRealtime *bool `json:"clock_realtime,omitempty"`
}

// ─── VM state ────────────────────────────────────────────────────────────────

// VmState enumerates the states that can be requested via PATCH /vm.
type VmState string

const (
	VmStatePaused  VmState = "Paused"
	VmStateResumed VmState = "Resumed"
)

// Vm requests a state transition for the running MicroVM (PATCH /vm).
type Vm struct {
	// State is the target state: Paused or Resumed.
	State VmState `json:"state"`
}

// ─── Logging and metrics ─────────────────────────────────────────────────────

// LogLevel enumerates the available log verbosity levels.
type LogLevel string

const (
	LogLevelError   LogLevel = "Error"
	LogLevelWarning LogLevel = "Warning"
	LogLevelInfo    LogLevel = "Info"
	LogLevelDebug   LogLevel = "Debug"
	LogLevelTrace   LogLevel = "Trace"
	LogLevelOff     LogLevel = "Off"
)

// Logger configures where and how Firecracker writes log messages
// (PUT /logger).
type Logger struct {
	// Level controls the minimum log severity that is emitted (default Info).
	Level *LogLevel `json:"level,omitempty"`
	// LogPath is the host path to the named pipe or file that receives log output.
	LogPath *string `json:"log_path,omitempty"`
	// ShowLevel includes the log level in each emitted log line.
	ShowLevel *bool `json:"show_level,omitempty"`
	// ShowLogOrigin includes the source file path and line number in each log line.
	ShowLogOrigin *bool `json:"show_log_origin,omitempty"`
	// Module filters log output to messages originating from the given Rust module path.
	Module *string `json:"module,omitempty"`
}

// Metrics configures where Firecracker writes JSON-formatted metrics
// (PUT /metrics).
type Metrics struct {
	// MetricsPath is the host path to the named pipe or file that receives metrics output.
	MetricsPath string `json:"metrics_path"`
}

// ─── MMDS ────────────────────────────────────────────────────────────────────

// MmdsVersion selects the MMDS protocol version exposed to the guest.
type MmdsVersion string

const (
	MmdsVersionV1 MmdsVersion = "V1"
	MmdsVersionV2 MmdsVersion = "V2"
)

// MmdsConfig configures the Microvm Metadata Service (pre-boot, PUT /mmds/config).
type MmdsConfig struct {
	// NetworkInterfaces lists the IDs of network interfaces that the guest can
	// use to reach the MMDS endpoint.
	NetworkInterfaces []string `json:"network_interfaces"`
	// Version selects the MMDS API version (default V1).
	// V2 adds session-oriented token-based requests and is more secure.
	Version *MmdsVersion `json:"version,omitempty"`
	// IPv4Address is the link-local IPv4 address of the MMDS endpoint inside the guest
	// (default 169.254.169.254).
	IPv4Address *string `json:"ipv4_address,omitempty"`
	// ImdsCompat enables EC2 IMDS compatibility mode.
	ImdsCompat *bool `json:"imds_compat,omitempty"`
}

// MmdsContentsObject is the free-form JSON document stored in and served by the MMDS.
type MmdsContentsObject map[string]interface{}

// ─── Vsock ───────────────────────────────────────────────────────────────────

// Vsock defines a virtio-vsock device backed by a Unix Domain Socket
// (pre-boot, PUT /vsock).
type Vsock struct {
	// GuestCID is the Context Identifier assigned to the guest (minimum 3;
	// 0–2 are reserved).
	GuestCID int `json:"guest_cid"`
	// UDSPath is the host path to the Unix domain socket used to proxy vsock connections.
	UDSPath string `json:"uds_path"`
	// VsockID is a deprecated identifier kept for backwards compatibility.
	VsockID *string `json:"vsock_id,omitempty"`
}

// ─── Entropy device ──────────────────────────────────────────────────────────

// EntropyDevice enables the virtio-rng entropy device
// (pre-boot, PUT /entropy).
type EntropyDevice struct {
	// RateLimiter limits the rate at which the device serves entropy to the guest.
	RateLimiter *RateLimiter `json:"rate_limiter,omitempty"`
}

// ─── Serial console ──────────────────────────────────────────────────────────

// SerialDevice configures the emulated serial console (PUT /serial).
type SerialDevice struct {
	// SerialOutPath is the host path to the named pipe or file that receives
	// serial output from the guest.
	SerialOutPath *string `json:"serial_out_path,omitempty"`
	// RateLimiter limits the bandwidth of serial output written to SerialOutPath.
	RateLimiter *TokenBucket `json:"rate_limiter,omitempty"`
}

// ─── Memory hotplug ──────────────────────────────────────────────────────────

// MemoryHotplugConfig configures the virtio-mem hotpluggable memory device
// (pre-boot, PUT /hotplug/memory).
type MemoryHotplugConfig struct {
	// TotalSizeMib is the total amount of memory that can be hotplugged, in MiB.
	TotalSizeMib *int `json:"total_size_mib,omitempty"`
	// SlotSizeMib is the granularity at which memory is hotplugged (default 128 MiB, minimum 128 MiB).
	SlotSizeMib *int `json:"slot_size_mib,omitempty"`
	// BlockSizeMib is the logical block size used by the virtio-mem device (default 2 MiB, minimum 2 MiB).
	BlockSizeMib *int `json:"block_size_mib,omitempty"`
}

// MemoryHotplugSizeUpdate requests a change to the plugged memory region size
// (post-boot, PATCH /hotplug/memory).
type MemoryHotplugSizeUpdate struct {
	// RequestedSizeMib is the new target plugged size in MiB. Must be a multiple of SlotSizeMib.
	RequestedSizeMib *int `json:"requested_size_mib,omitempty"`
}

// MemoryHotplugStatus describes the current state of the virtio-mem device
// (GET /hotplug/memory).
type MemoryHotplugStatus struct {
	TotalSizeMib     *int `json:"total_size_mib,omitempty"`
	SlotSizeMib      *int `json:"slot_size_mib,omitempty"`
	BlockSizeMib     *int `json:"block_size_mib,omitempty"`
	PluggedSizeMib   *int `json:"plugged_size_mib,omitempty"`
	RequestedSizeMib *int `json:"requested_size_mib,omitempty"`
}

// ─── Full configuration ──────────────────────────────────────────────────────

// FullVmConfiguration is the complete snapshot of a VM's configuration
// (GET /vm/config).
type FullVmConfiguration struct {
	Balloon           *Balloon              `json:"balloon,omitempty"`
	Drives            []Drive               `json:"drives,omitempty"`
	BootSource        *BootSource           `json:"boot-source,omitempty"`
	CpuConfig         *CpuConfig            `json:"cpu-config,omitempty"`
	Logger            *Logger               `json:"logger,omitempty"`
	MachineConfig     *MachineConfiguration `json:"machine-config,omitempty"`
	Metrics           *Metrics              `json:"metrics,omitempty"`
	MemoryHotplug     *MemoryHotplugConfig  `json:"memory-hotplug,omitempty"`
	MmdsConfig        *MmdsConfig           `json:"mmds-config,omitempty"`
	NetworkInterfaces []NetworkInterface    `json:"network-interfaces,omitempty"`
	Pmem              []Pmem                `json:"pmem,omitempty"`
	Vsock             *Vsock                `json:"vsock,omitempty"`
	Entropy           *EntropyDevice        `json:"entropy,omitempty"`
}
