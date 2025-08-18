package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"gpu-management/internal/api/middleware"
	"gpu-management/internal/models"
	"gpu-management/internal/services/event"
)

// GPUHandler GPU处理器
type GPUHandler struct {
	eventBus event.EventBus
}

// NewGPUHandler 创建新的GPU处理器
func NewGPUHandler(eventBus event.EventBus) *GPUHandler {
	return &GPUHandler{
		eventBus: eventBus,
	}
}

// List 列出GPU列表
func (h *GPUHandler) List(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	// 添加请求追踪信息
	requestID := middleware.GetRequestID(ctx)
	operation := middleware.GetOperation(ctx)

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	gpus, err := h.listGPUs(businessCtx)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":       gpus,
		"total":      len(gpus),
		"request_id": requestID,
		"operation":  operation,
	})
}

// Get 获取GPU详情
func (h *GPUHandler) Get(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	gpu, err := h.getGPU(businessCtx, id)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, gpu)
}

// Update 更新GPU信息
func (h *GPUHandler) Update(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	var gpu models.GPU
	if err := c.Bind(&gpu); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	err := h.updateGPU(businessCtx, id, &gpu)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "更新超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	gpu.ID = id
	return c.JSON(http.StatusOK, gpu)
}

// GetStatus 获取GPU状态
func (h *GPUHandler) GetStatus(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	status, err := h.getGPUStatus(businessCtx, id)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

// GetMetrics 获取GPU指标
func (h *GPUHandler) GetMetrics(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	metrics, err := h.getGPUMetrics(businessCtx, id)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, metrics)
}

// UpdateConfig 更新GPU配置
func (h *GPUHandler) UpdateConfig(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	var config models.GPUConfig
	if err := c.Bind(&config); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	err := h.updateGPUConfig(businessCtx, id, &config)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "配置更新超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 发布配置更新事件
	event := event.HardwareEvent{
		Event: event.Event{
			Type:   event.EventTypeHardwareProvisioned,
			Source: "gpu-service",
		},
		HardwareID: id,
		Status:     "config_updated",
	}

	// 使用Context发布事件
	if err := h.eventBus.Publish(businessCtx, "gpu.config.updated", event); err != nil {
		// 事件发布失败不影响主流程，只记录日志
		c.Logger().Errorf("Failed to publish event: %v", err)
	}

	config.GPUID = id
	return c.JSON(http.StatusOK, config)
}

// GetConfig 获取GPU配置
func (h *GPUHandler) GetConfig(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	config, err := h.getGPUConfig(businessCtx, id)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, config)
}

// GetConfigVersions 获取GPU配置版本
func (h *GPUHandler) GetConfigVersions(c echo.Context) error {
	// 获取请求Context
	ctx := c.Request().Context()

	id := c.Param("id")

	// 设置业务超时
	businessCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用服务层，传递Context
	versions, err := h.getGPUConfigVersions(businessCtx, id)
	if err != nil {
		// 检查Context相关错误
		if businessCtx.Err() == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusRequestTimeout, "查询超时")
		}
		if businessCtx.Err() == context.Canceled {
			return echo.NewHTTPError(http.StatusRequestTimeout, "请求被取消")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, versions)
}

// 服务层方法实现
func (h *GPUHandler) listGPUs(ctx context.Context) ([]models.GPU, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU列表查询
	gpus := []models.GPU{
		{
			ID:       "gpu-001",
			ServerID: "server-001",
			Model:    "NVIDIA A100",
			Status:   models.GPUStatusAvailable,
			MemoryGB: 80,
		},
	}

	return gpus, nil
}

func (h *GPUHandler) getGPU(ctx context.Context, id string) (*models.GPU, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU详情查询
	gpu := &models.GPU{
		ID:       id,
		ServerID: "server-001",
		Model:    "NVIDIA A100",
		Status:   models.GPUStatusAvailable,
		MemoryGB: 80,
	}

	return gpu, nil
}

func (h *GPUHandler) updateGPU(ctx context.Context, id string, gpu *models.GPU) error {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TODO: 实现GPU信息更新
	return nil
}

func (h *GPUHandler) getGPUStatus(ctx context.Context, id string) (map[string]interface{}, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU状态查询
	status := map[string]interface{}{
		"gpu_id": id,
		"status": models.GPUStatusAvailable,
		"utilization": map[string]interface{}{
			"gpu_usage_percent":     45.2,
			"memory_usage_percent":  60.8,
			"compute_usage_percent": 42.1,
			"temperature":           65.0,
			"power_consumption":     180.5,
		},
	}

	return status, nil
}

func (h *GPUHandler) getGPUMetrics(ctx context.Context, id string) (*models.GPUUtilization, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU指标查询
	metrics := &models.GPUUtilization{
		GPUID:                id,
		GPUUsagePercent:      45.2,
		MemoryUsagePercent:   60.8,
		ComputeUsagePercent:  42.1,
		MemoryBandwidthUsage: 75.3,
		Temperature:          65.0,
		PowerConsumption:     180.5,
	}

	return metrics, nil
}

func (h *GPUHandler) updateGPUConfig(ctx context.Context, id string, config *models.GPUConfig) error {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TODO: 实现GPU配置更新
	return nil
}

func (h *GPUHandler) getGPUConfig(ctx context.Context, id string) (*models.GPUConfig, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU配置查询
	config := &models.GPUConfig{
		GPUID:       id,
		Version:     "v1.0",
		PowerLimit:  200,
		MemoryClock: 8000,
		DriverConfig: map[string]string{
			"driver_version": "470.82.01",
			"cuda_version":   "11.4",
		},
	}

	return config, nil
}

func (h *GPUHandler) getGPUConfigVersions(ctx context.Context, id string) ([]models.GPUConfig, error) {
	// 检查Context状态
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// TODO: 实现GPU配置版本查询
	versions := []models.GPUConfig{
		{
			ID:      "config-001",
			GPUID:   id,
			Version: "v1.0",
		},
		{
			ID:      "config-002",
			GPUID:   id,
			Version: "v1.1",
		},
	}

	return versions, nil
}
