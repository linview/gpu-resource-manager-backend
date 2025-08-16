package api

import (
	"net/http"
	"mock-tinkerbell/services"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	workflowService *services.WorkflowService
}

func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		HardwareID  string `json:"hardware_id"`
		TemplateID  string `json:"template_id"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	workflow, err := h.workflowService.CreateWorkflow(req.Name, req.HardwareID, req.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, workflow)
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	id := c.Param("id")
	workflow, exists := h.workflowService.GetWorkflow(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}
	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) GetAllWorkflows(c *gin.Context) {
	workflows := h.workflowService.GetAllWorkflows()
	c.JSON(http.StatusOK, gin.H{
		"data": workflows,
		"total": len(workflows),
	})
}

func (h *WorkflowHandler) StartWorkflow(c *gin.Context) {
	id := c.Param("id")
	err := h.workflowService.StartWorkflow(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workflow started successfully"})
}

func (h *WorkflowHandler) GetTemplates(c *gin.Context) {
	templates := h.workflowService.GetAllTemplates()
	c.JSON(http.StatusOK, gin.H{
		"data": templates,
		"total": len(templates),
	})
}

func (h *WorkflowHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, exists := h.workflowService.GetTemplate(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *WorkflowHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	template := h.workflowService.CreateTemplate(req.Name, req.Description, req.Version)
	c.JSON(http.StatusCreated, template)
}

func (h *WorkflowHandler) GetWorkflowStats(c *gin.Context) {
	stats := h.workflowService.GetWorkflowStats()
	c.JSON(http.StatusOK, stats)
}
