package models

import (
	"time"
	"github.com/google/uuid"
)

// Hardware 硬件设备模型
type Hardware struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"` // server, gpu, bmc
	Status       string    `json:"status"` // discovered, provisioning, ready, error
	IPAddress    string    `json:"ip_address"`
	MACAddress   string    `json:"mac_address"`
	SerialNumber string    `json:"serial_number"`
	Model        string    `json:"model"`
	Specs        HardwareSpecs `json:"specs"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// HardwareSpecs 硬件规格
type HardwareSpecs struct {
	CPU        CPUSpec     `json:"cpu,omitempty"`
	Memory     MemorySpec  `json:"memory,omitempty"`
	Storage    StorageSpec `json:"storage,omitempty"`
	GPU        []GPUSpec   `json:"gpu,omitempty"`
	Network    NetworkSpec `json:"network,omitempty"`
}

// CPUSpec CPU规格
type CPUSpec struct {
	Model       string `json:"model"`
	Cores       int    `json:"cores"`
	Threads     int    `json:"threads"`
	Frequency   string `json:"frequency"`
}

// MemorySpec 内存规格
type MemorySpec struct {
	TotalGB    int `json:"total_gb"`
	Slots      int `json:"slots"`
	UsedSlots  int `json:"used_slots"`
}

// StorageSpec 存储规格
type StorageSpec struct {
	Disks []DiskSpec `json:"disks"`
}

// DiskSpec 磁盘规格
type DiskSpec struct {
	Device     string `json:"device"`
	SizeGB     int    `json:"size_gb"`
	Type       string `json:"type"` // ssd, hdd, nvme
	Interface  string `json:"interface"` // sata, sas, nvme
}

// GPUSpec GPU规格
type GPUSpec struct {
	ID            string `json:"id"`
	Model         string `json:"model"`
	MemoryGB      int    `json:"memory_gb"`
	DriverVersion string `json:"driver_version"`
	CUDAVersion   string `json:"cuda_version"`
	Status        string `json:"status"` // available, allocated, in_use, error
	SerialNumber  string `json:"serial_number"`
	Temperature   int    `json:"temperature"`
	PowerUsage    int    `json:"power_usage"`
	Utilization   int    `json:"utilization"`
}

// NetworkSpec 网络规格
type NetworkSpec struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

// NetworkInterface 网络接口
type NetworkInterface struct {
	Name       string `json:"name"`
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Speed      string `json:"speed"`
	Duplex     string `json:"duplex"`
}

// NewHardware 创建新的硬件设备
func NewHardware(name, hwType, model string) *Hardware {
	return &Hardware{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      hwType,
		Status:    "discovered",
		Model:     model,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewServer 创建新的服务器
func NewServer(name, model string) *Hardware {
	server := NewHardware(name, "server", model)
	server.Specs = HardwareSpecs{
		CPU: CPUSpec{
			Model:     "Intel Xeon E5-2680 v4",
			Cores:     14,
			Threads:   28,
			Frequency: "2.4GHz",
		},
		Memory: MemorySpec{
			TotalGB:   256,
			Slots:     16,
			UsedSlots: 8,
		},
		Storage: StorageSpec{
			Disks: []DiskSpec{
				{
					Device:    "/dev/sda",
					SizeGB:    1000,
					Type:      "ssd",
					Interface: "sata",
				},
			},
		},
		Network: NetworkSpec{
			Interfaces: []NetworkInterface{
				{
					Name:       "eth0",
					MACAddress: "00:15:5d:01:ca:05",
					IPAddress:  "192.168.1.100",
					Speed:      "1Gbps",
					Duplex:     "full",
				},
			},
		},
	}
	return server
}

// AddGPU 添加GPU到服务器
func (h *Hardware) AddGPU(gpuSpec GPUSpec) {
	gpuSpec.ID = uuid.New().String()
	gpuSpec.Status = "available"
	h.Specs.GPU = append(h.Specs.GPU, gpuSpec)
	h.UpdatedAt = time.Now()
}

// UpdateStatus 更新硬件状态
func (h *Hardware) UpdateStatus(status string) {
	h.Status = status
	h.UpdatedAt = time.Now()
}

// UpdateGPUStatus 更新GPU状态
func (h *Hardware) UpdateGPUStatus(gpuID, status string) {
	for i, gpu := range h.Specs.GPU {
		if gpu.ID == gpuID {
			h.Specs.GPU[i].Status = status
			h.UpdatedAt = time.Now()
			break
		}
	}
}
