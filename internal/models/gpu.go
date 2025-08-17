package models

import (
	"time"
)

// GPU GPU资源模型
type GPU struct {
	ID            string    `json:"id" db:"id"`
	ServerID      string    `json:"server_id" db:"server_id"`
	Model         string    `json:"model" db:"model"`
	Status        string    `json:"status" db:"status"`
	MemoryGB      int       `json:"memory_gb" db:"memory_gb"`
	SerialNumber  string    `json:"serial_number" db:"serial_number"`
	DriverVersion string    `json:"driver_version" db:"driver_version"`
	CUDAVersion   string    `json:"cuda_version" db:"cuda_version"`
	GNMIDeviceID  string    `json:"gnmi_device_id" db:"gnmi_device_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// GPUStatus GPU状态枚举
const (
	GPUStatusAvailable = "available"
	GPUStatusAllocated = "allocated"
	GPUStatusInUse     = "in_use"
	GPUStatusError     = "error"
)

// GPUConfig GPU配置模型
type GPUConfig struct {
	ID           string            `json:"id" db:"id"`
	GPUID        string            `json:"gpu_id" db:"gpu_id"`
	Version      string            `json:"version" db:"version"`
	PowerLimit   int               `json:"power_limit" db:"power_limit"`
	MemoryClock  int               `json:"memory_clock" db:"memory_clock"`
	DriverConfig map[string]string `json:"driver_config" db:"driver_config"`
	CreatedBy    string            `json:"created_by" db:"created_by"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
}

// GPUUtilization GPU利用率模型
type GPUUtilization struct {
	GPUID                    string  `json:"gpu_id"`
	GPUUsagePercent          float64 `json:"gpu_usage_percent"`
	MemoryUsagePercent       float64 `json:"memory_usage_percent"`
	ComputeUsagePercent      float64 `json:"compute_usage_percent"`
	MemoryBandwidthUsage     float64 `json:"memory_bandwidth_usage"`
	Temperature              float64 `json:"temperature"`
	PowerConsumption         float64 `json:"power_consumption"`
	Timestamp                time.Time `json:"timestamp"`
}

// GPUEvent GPU事件模型
type GPUEvent struct {
	ID          string            `json:"id"`
	GPUID       string            `json:"gpu_id"`
	EventType   string            `json:"event_type"`
	Severity    string            `json:"severity"`
	Message     string            `json:"message"`
	Metadata    map[string]string `json:"metadata"`
	Timestamp   time.Time         `json:"timestamp"`
}

// GPUEventType GPU事件类型枚举
const (
	GPUEventTypeStatusChanged = "status_changed"
	GPUEventTypeError         = "error"
	GPUEventTypeAllocated     = "allocated"
	GPUEventTypeReleased      = "released"
)

// GPUSeverity GPU事件严重程度枚举
const (
	GPUSeverityInfo    = "info"
	GPUSeverityWarning = "warning"
	GPUSeverityError   = "error"
)
