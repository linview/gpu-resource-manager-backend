package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"gpu-management/internal/api/handlers"
	"gpu-management/internal/services/event"
)

// Setup 设置路由
func Setup(e *echo.Echo, eventBus event.EventBus) {
	// 创建处理器
	gpuHandler := handlers.NewGPUHandler(eventBus)
	serverHandler := handlers.NewServerHandler(eventBus)
	allocationHandler := handlers.NewAllocationHandler(eventBus)
	eventHandler := handlers.NewEventHandler(eventBus)

	// API v1 路由组
	v1 := e.Group("/api/v1")

	// GPU管理路由
	gpus := v1.Group("/gpus")
	gpus.GET("", gpuHandler.List)
	gpus.GET("/:id", gpuHandler.Get)
	gpus.PUT("/:id", gpuHandler.Update)
	gpus.GET("/:id/status", gpuHandler.GetStatus)
	gpus.GET("/:id/metrics", gpuHandler.GetMetrics)
	gpus.POST("/:id/config", gpuHandler.UpdateConfig)
	gpus.GET("/:id/config", gpuHandler.GetConfig)
	gpus.GET("/:id/config/versions", gpuHandler.GetConfigVersions)

	// 服务器管理路由
	servers := v1.Group("/servers")
	servers.GET("", serverHandler.List)
	servers.POST("", serverHandler.Create)
	servers.GET("/:id", serverHandler.Get)
	servers.PUT("/:id", serverHandler.Update)
	servers.DELETE("/:id", serverHandler.Delete)
	servers.POST("/:id/power", serverHandler.PowerControl)
	servers.GET("/:id/status", serverHandler.GetStatus)
	servers.POST("/:id/bios", serverHandler.ConfigureBIOS)
	servers.POST("/:id/firmware", serverHandler.UpgradeFirmware)
	servers.GET("/:id/gpus", serverHandler.GetGPUs)
	servers.POST("/:id/gpus", serverHandler.AddGPU)
	servers.GET("/:id/config", serverHandler.GetConfig)
	servers.PUT("/:id/config", serverHandler.UpdateConfig)
	servers.GET("/:id/config/versions", serverHandler.GetConfigVersions)

	// 资源分配路由
	allocations := v1.Group("/allocations")
	allocations.GET("", allocationHandler.List)
	allocations.POST("", allocationHandler.Create)
	allocations.GET("/:id", allocationHandler.Get)
	allocations.PUT("/:id", allocationHandler.Update)
	allocations.DELETE("/:id", allocationHandler.Delete)
	allocations.POST("/:id/start", allocationHandler.Start)
	allocations.POST("/:id/stop", allocationHandler.Stop)
	allocations.GET("/:id/status", allocationHandler.GetStatus)

	// 工作流路由
	workflows := v1.Group("/workflows")
	workflows.GET("", allocationHandler.ListWorkflows)
	workflows.POST("", allocationHandler.CreateWorkflow)
	workflows.GET("/:id", allocationHandler.GetWorkflow)
	workflows.DELETE("/:id", allocationHandler.DeleteWorkflow)
	workflows.POST("/deploy", allocationHandler.CreateDeployWorkflow)
	workflows.POST("/cleanup", allocationHandler.CreateCleanupWorkflow)

	// 事件路由
	events := v1.Group("/events")
	events.GET("", eventHandler.List)
	events.POST("", eventHandler.Create)
	events.GET("/:id", eventHandler.Get)

	// 告警路由
	alerts := v1.Group("/alerts")
	alerts.GET("", eventHandler.ListAlerts)
	alerts.POST("", eventHandler.CreateAlert)
	alerts.PUT("/:id", eventHandler.UpdateAlert)
	alerts.GET("/:id/correlation", eventHandler.GetAlertCorrelation)

	// 健康检查
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
}
