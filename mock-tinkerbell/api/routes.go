package api

import (
	"mock-tinkerbell/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(hardwareService *services.HardwareService, workflowService *services.WorkflowService) *gin.Engine {
	router := gin.Default()
	
	// 创建处理器
	hardwareHandler := NewHardwareHandler(hardwareService)
	workflowHandler := NewWorkflowHandler(workflowService)
	
	// API版本
	v1 := router.Group("/api/v1")
	
	// 硬件管理API
	hardware := v1.Group("/hardware")
	{
		hardware.GET("", hardwareHandler.GetHardware)
		hardware.GET("/stats", hardwareHandler.GetHardwareStats)
		hardware.GET("/discover", hardwareHandler.SimulateDiscovery)
		hardware.GET("/:id", hardwareHandler.GetHardwareByID)
		hardware.PUT("/:id/status", hardwareHandler.UpdateHardwareStatus)
	}
	
	// 工作流管理API
	workflows := v1.Group("/workflows")
	{
		workflows.GET("", workflowHandler.GetAllWorkflows)
		workflows.GET("/stats", workflowHandler.GetWorkflowStats)
		workflows.POST("", workflowHandler.CreateWorkflow)
		workflows.GET("/:id", workflowHandler.GetWorkflow)
		workflows.POST("/:id/start", workflowHandler.StartWorkflow)
	}
	
	// 模板管理API
	templates := v1.Group("/templates")
	{
		templates.GET("", workflowHandler.GetTemplates)
		templates.POST("", workflowHandler.CreateTemplate)
		templates.GET("/:id", workflowHandler.GetTemplate)
	}
	
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "mock-tinkerbell",
			"version": "1.0.0",
		})
	})
	
	// API文档
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Mock Tinkerbell API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"hardware": gin.H{
					"GET /api/v1/hardware": "获取所有硬件",
					"GET /api/v1/hardware/:id": "获取指定硬件",
					"PUT /api/v1/hardware/:id/status": "更新硬件状态",
					"GET /api/v1/hardware/stats": "获取硬件统计",
					"GET /api/v1/hardware/discover": "模拟硬件发现",
				},
				"workflows": gin.H{
					"GET /api/v1/workflows": "获取所有工作流",
					"POST /api/v1/workflows": "创建工作流",
					"GET /api/v1/workflows/:id": "获取指定工作流",
					"POST /api/v1/workflows/:id/start": "启动工作流",
					"GET /api/v1/workflows/stats": "获取工作流统计",
				},
				"templates": gin.H{
					"GET /api/v1/templates": "获取所有模板",
					"POST /api/v1/templates": "创建模板",
					"GET /api/v1/templates/:id": "获取指定模板",
				},
			},
		})
	})
	
	return router
}
