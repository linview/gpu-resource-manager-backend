package models

import (
	"time"
	"github.com/google/uuid"
)

// Workflow 工作流模型
type Workflow struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	HardwareID  string            `json:"hardware_id"`
	TemplateID  string            `json:"template_id"`
	Status      string            `json:"status"` // pending, running, completed, failed
	Steps       []WorkflowStep    `json:"steps"`
	CurrentStep int               `json:"current_step"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Action      string                 `json:"action"`
	Status      string                 `json:"status"` // pending, running, completed, failed
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Duration    int                    `json:"duration_ms,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
}

// Template 工作流模板
type Template struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Steps       []TemplateStep    `json:"steps"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TemplateStep 模板步骤
type TemplateStep struct {
	Name       string                 `json:"name"`
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Condition  string                 `json:"condition,omitempty"`
	Timeout    int                    `json:"timeout_seconds,omitempty"`
}

// NewWorkflow 创建新的工作流
func NewWorkflow(name, hardwareID, templateID string) *Workflow {
	return &Workflow{
		ID:         uuid.New().String(),
		Name:       name,
		HardwareID: hardwareID,
		TemplateID: templateID,
		Status:     "pending",
		CurrentStep: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Metadata:   make(map[string]string),
	}
}

// AddStep 添加工作流步骤
func (w *Workflow) AddStep(step WorkflowStep) {
	step.ID = uuid.New().String()
	step.Status = "pending"
	w.Steps = append(w.Steps, step)
	w.UpdatedAt = time.Now()
}

// Start 开始工作流
func (w *Workflow) Start() {
	w.Status = "running"
	now := time.Now()
	w.StartedAt = &now
	w.UpdatedAt = now
}

// Complete 完成工作流
func (w *Workflow) Complete() {
	w.Status = "completed"
	now := time.Now()
	w.CompletedAt = &now
	w.UpdatedAt = now
}

// Fail 工作流失败
func (w *Workflow) Fail(error string) {
	w.Status = "failed"
	now := time.Now()
	w.CompletedAt = &now
	w.UpdatedAt = now
	w.Metadata["error"] = error
}

// UpdateStepStatus 更新步骤状态
func (w *Workflow) UpdateStepStatus(stepIndex int, status string, error string) {
	if stepIndex >= 0 && stepIndex < len(w.Steps) {
		step := &w.Steps[stepIndex]
		step.Status = status
		
		if status == "running" && step.StartedAt == nil {
			now := time.Now()
			step.StartedAt = &now
		}
		
		if status == "completed" || status == "failed" {
			now := time.Now()
			step.CompletedAt = &now
			if step.StartedAt != nil {
				step.Duration = int(now.Sub(*step.StartedAt).Milliseconds())
			}
		}
		
		if error != "" {
			step.Error = error
		}
		
		w.UpdatedAt = time.Now()
	}
}

// GetCurrentStep 获取当前步骤
func (w *Workflow) GetCurrentStep() *WorkflowStep {
	if w.CurrentStep >= 0 && w.CurrentStep < len(w.Steps) {
		return &w.Steps[w.CurrentStep]
	}
	return nil
}

// NextStep 移动到下一步
func (w *Workflow) NextStep() {
	w.CurrentStep++
	w.UpdatedAt = time.Now()
}

// IsCompleted 检查工作流是否完成
func (w *Workflow) IsCompleted() bool {
	return w.Status == "completed" || w.Status == "failed"
}

// NewTemplate 创建新的模板
func NewTemplate(name, description, version string) *Template {
	return &Template{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Version:     version,
		Steps:       []TemplateStep{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}
}

// AddTemplateStep 添加模板步骤
func (t *Template) AddTemplateStep(step TemplateStep) {
	t.Steps = append(t.Steps, step)
	t.UpdatedAt = time.Now()
}
