package models

import (
	"time"
)

// Server 服务器模型
type Server struct {
	ID                   string    `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	Model                string    `json:"model" db:"model"`
	Status               string    `json:"status" db:"status"`
	IPMIIP               string    `json:"ipmi_ip" db:"ipmi_ip"`
	RedfishURL           string    `json:"redfish_url" db:"redfish_url"`
	ManagementIP         string    `json:"management_ip" db:"management_ip"`
	SerialNumber         string    `json:"serial_number" db:"serial_number"`
	TinkerbellHardwareID string    `json:"tinkerbell_hardware_id" db:"tinkerbell_hardware_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// ServerStatus 服务器状态枚举
const (
	ServerStatusDiscovered   = "discovered"
	ServerStatusProvisioning = "provisioning"
	ServerStatusReady        = "ready"
	ServerStatusError        = "error"
)

// ServerConfig 服务器配置模型
type ServerConfig struct {
	ID            string            `json:"id" db:"id"`
	ServerID      string            `json:"server_id" db:"server_id"`
	Version       string            `json:"version" db:"version"`
	BIOSConfig    map[string]string `json:"bios_config" db:"bios_config"`
	RAIDConfig    map[string]string `json:"raid_config" db:"raid_config"`
	NetworkConfig map[string]string `json:"network_config" db:"network_config"`
	CreatedBy     string            `json:"created_by" db:"created_by"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
}

// HardwareCRD Tinkerbell硬件CRD模型
type HardwareCRD struct {
	ID            string                 `json:"id" db:"id"`
	TinkerbellID  string                 `json:"tinkerbell_id" db:"tinkerbell_id"`
	ServerID      string                 `json:"server_id" db:"server_id"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	Status        string                 `json:"status" db:"status"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}

// TemplateCRD Tinkerbell模板CRD模型
type TemplateCRD struct {
	ID          string    `json:"id" db:"id"`
	TinkerbellID string   `json:"tinkerbell_id" db:"tinkerbell_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Content     string    `json:"content" db:"content"`
	Version     string    `json:"version" db:"version"`
	OSType      string    `json:"os_type" db:"os_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WorkflowCRD Tinkerbell工作流CRD模型
type WorkflowCRD struct {
	ID            string                 `json:"id" db:"id"`
	TinkerbellID  string                 `json:"tinkerbell_id" db:"tinkerbell_id"`
	HardwareID    string                 `json:"hardware_id" db:"hardware_id"`
	TemplateID    string                 `json:"template_id" db:"template_id"`
	Status        string                 `json:"status" db:"status"`
	Steps         map[string]interface{} `json:"steps" db:"steps"`
	StartedAt     *time.Time             `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at" db:"completed_at"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// WorkflowStatus 工作流状态枚举
const (
	WorkflowStatusPending   = "pending"
	WorkflowStatusRunning   = "running"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusFailed    = "failed"
)
