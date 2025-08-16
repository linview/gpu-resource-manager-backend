package services

import (
	"fmt"
	"sync"
	"time"
	"mock-tinkerbell/config"
	"mock-tinkerbell/models"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type WorkflowService struct {
	workflows map[string]*models.Workflow
	templates map[string]*models.Template
	mutex     sync.RWMutex
	logger    *logrus.Logger
	config    *config.Config
}

func NewWorkflowService(cfg *config.Config) *WorkflowService {
	service := &WorkflowService{
		workflows: make(map[string]*models.Workflow),
		templates: make(map[string]*models.Template),
		logger:    logrus.New(),
		config:    cfg,
	}
	service.initializeDefaultTemplates()
	return service
}

func (service *WorkflowService) initializeDefaultTemplates() {
	// 创建默认的OS安装模板
	osTemplate := models.NewTemplate("ubuntu-22.04-install", "Ubuntu 22.04 LTS 安装模板", "1.0.0")
	osTemplate.AddTemplateStep(models.TemplateStep{
		Name:   "power-on",
		Action: "power_on",
		Parameters: map[string]interface{}{
			"timeout": 60,
		},
	})
	osTemplate.AddTemplateStep(models.TemplateStep{
		Name:   "install-os",
		Action: "install_os",
		Parameters: map[string]interface{}{
			"os":      "ubuntu-22.04",
			"timeout": 1800,
		},
	})
	osTemplate.AddTemplateStep(models.TemplateStep{
		Name:   "install-gpu-drivers",
		Action: "install_gpu_drivers",
		Parameters: map[string]interface{}{
			"driver_version": "525.85.05",
			"timeout":        600,
		},
	})
	service.templates[osTemplate.ID] = osTemplate

	// 创建清理模板
	cleanupTemplate := models.NewTemplate("server-cleanup", "服务器清理模板", "1.0.0")
	cleanupTemplate.AddTemplateStep(models.TemplateStep{
		Name:   "cleanup-os",
		Action: "cleanup_os",
		Parameters: map[string]interface{}{
			"timeout": 300,
		},
	})
	cleanupTemplate.AddTemplateStep(models.TemplateStep{
		Name:   "power-off",
		Action: "power_off",
		Parameters: map[string]interface{}{
			"timeout": 60,
		},
	})
	service.templates[cleanupTemplate.ID] = cleanupTemplate
}

func (service *WorkflowService) CreateWorkflow(name, hardwareID, templateID string) (*models.Workflow, error) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	template, exists := service.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}

	workflow := models.NewWorkflow(name, hardwareID, templateID)
	
	// 根据模板创建步骤
	for _, templateStep := range template.Steps {
		step := models.WorkflowStep{
			Name:       templateStep.Name,
			Action:     templateStep.Action,
			Status:     "pending",
			Parameters: templateStep.Parameters,
		}
		workflow.AddStep(step)
	}

	service.workflows[workflow.ID] = workflow
	service.logger.Infof("Created workflow: %s (%s)", workflow.Name, workflow.ID)
	
	return workflow, nil
}

func (service *WorkflowService) GetWorkflow(id string) (*models.Workflow, bool) {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	workflow, exists := service.workflows[id]
	return workflow, exists
}

func (service *WorkflowService) GetAllWorkflows() []*models.Workflow {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	workflows := make([]*models.Workflow, 0, len(service.workflows))
	for _, wf := range service.workflows {
		workflows = append(workflows, wf)
	}
	return workflows
}

func (service *WorkflowService) StartWorkflow(id string) error {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	workflow, exists := service.workflows[id]
	if !exists {
		return fmt.Errorf("workflow not found: %s", id)
	}
	
	if workflow.Status != "pending" {
		return fmt.Errorf("workflow is not in pending status: %s", workflow.Status)
	}
	
	workflow.Start()
	service.logger.Infof("Started workflow: %s", id)
	
	// 异步执行工作流
	go service.executeWorkflow(workflow)
	
	return nil
}

func (service *WorkflowService) executeWorkflow(workflow *models.Workflow) {
	service.logger.Infof("Executing workflow: %s", workflow.ID)
	
	for i := 0; i < len(workflow.Steps); i++ {
		step := &workflow.Steps[i]
		
		// 更新步骤状态为运行中
		service.updateStepStatus(workflow.ID, i, "running", "")
		
		// 模拟步骤执行
		success := service.executeStep(step)
		
		if success {
			service.updateStepStatus(workflow.ID, i, "completed", "")
		} else {
			service.updateStepStatus(workflow.ID, i, "failed", "Step execution failed")
			service.failWorkflow(workflow.ID, "Step execution failed")
			return
		}
		
		// 步骤间延迟
		time.Sleep(time.Duration(service.config.Workflow.StepDelay) * time.Millisecond)
	}
	
	// 完成工作流
	service.completeWorkflow(workflow.ID)
}

func (service *WorkflowService) executeStep(step *models.WorkflowStep) bool {
	// 模拟步骤执行时间
	executionTime := time.Duration(rand.Intn(5000)+2000) * time.Millisecond
	time.Sleep(executionTime)
	
	// 模拟成功率（90%成功）
	if rand.Float64() < 0.9 {
		step.Output = map[string]interface{}{
			"status": "success",
			"message": fmt.Sprintf("Step %s completed successfully", step.Name),
		}
		return true
	} else {
		step.Output = map[string]interface{}{
			"status": "failed",
			"message": fmt.Sprintf("Step %s failed", step.Name),
		}
		return false
	}
}

func (service *WorkflowService) updateStepStatus(workflowID string, stepIndex int, status, error string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	workflow, exists := service.workflows[workflowID]
	if !exists {
		return
	}
	
	workflow.UpdateStepStatus(stepIndex, status, error)
}

func (service *WorkflowService) completeWorkflow(workflowID string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	workflow, exists := service.workflows[workflowID]
	if !exists {
		return
	}
	
	workflow.Complete()
	service.logger.Infof("Completed workflow: %s", workflowID)
}

func (service *WorkflowService) failWorkflow(workflowID, error string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	workflow, exists := service.workflows[workflowID]
	if !exists {
		return
	}
	
	workflow.Fail(error)
	service.logger.Errorf("Failed workflow: %s - %s", workflowID, error)
}

func (service *WorkflowService) GetTemplate(id string) (*models.Template, bool) {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	template, exists := service.templates[id]
	return template, exists
}

func (service *WorkflowService) GetAllTemplates() []*models.Template {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	templates := make([]*models.Template, 0, len(service.templates))
	for _, tmpl := range service.templates {
		templates = append(templates, tmpl)
	}
	return templates
}

func (service *WorkflowService) CreateTemplate(name, description, version string) *models.Template {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	template := models.NewTemplate(name, description, version)
	service.templates[template.ID] = template
	service.logger.Infof("Created template: %s (%s)", template.Name, template.ID)
	
	return template
}

func (service *WorkflowService) GetWorkflowStats() map[string]interface{} {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	stats := map[string]interface{}{
		"total_workflows": len(service.workflows),
		"total_templates": len(service.templates),
		"status_counts":   map[string]int{},
	}
	
	for _, workflow := range service.workflows {
		count := stats["status_counts"].(map[string]int)[workflow.Status]
		stats["status_counts"].(map[string]int)[workflow.Status] = count + 1
	}
	
	return stats
}
