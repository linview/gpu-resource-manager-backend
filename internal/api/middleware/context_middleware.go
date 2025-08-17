package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// ContextKey 定义Context key类型
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	OperationKey ContextKey = "operation"
)

// ContextMiddleware Context中间件
func ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取请求Context
			ctx := c.Request().Context()

			// 生成请求ID
			requestID := generateRequestID()
			ctx = context.WithValue(ctx, RequestIDKey, requestID)

			// 添加用户ID（如果有的话）
			if userID := getUserID(c); userID != "" {
				ctx = context.WithValue(ctx, UserIDKey, userID)
			}

			// 添加操作类型
			operation := getOperation(c)
			ctx = context.WithValue(ctx, OperationKey, operation)

			// 设置请求超时（60秒）
			timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			// 更新请求Context
			c.SetRequest(c.Request().WithContext(timeoutCtx))

			// 调用下一个处理器
			return next(c)
		}
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// getUserID 从请求中获取用户ID
func getUserID(c echo.Context) string {
	// 这里可以从JWT token、header或其他地方获取用户ID
	// 暂时返回空字符串，后续可以扩展
	return ""
}

// getOperation 获取操作类型
func getOperation(c echo.Context) string {
	// 根据HTTP方法和路径确定操作类型
	method := c.Request().Method
	path := c.Path()

	switch {
	case method == "GET" && path == "/api/v1/gpus":
		return "gpu_list"
	case method == "GET" && path == "/api/v1/gpus/:id":
		return "gpu_get"
	case method == "PUT" && path == "/api/v1/gpus/:id":
		return "gpu_update"
	case method == "GET" && path == "/api/v1/gpus/:id/status":
		return "gpu_status"
	case method == "GET" && path == "/api/v1/gpus/:id/metrics":
		return "gpu_metrics"
	case method == "PUT" && path == "/api/v1/gpus/:id/config":
		return "gpu_config_update"
	case method == "GET" && path == "/api/v1/gpus/:id/config":
		return "gpu_config_get"
	case method == "GET" && path == "/api/v1/gpus/:id/config/versions":
		return "gpu_config_versions"
	case method == "GET" && path == "/api/v1/allocations":
		return "allocation_list"
	case method == "POST" && path == "/api/v1/allocations":
		return "allocation_create"
	case method == "GET" && path == "/api/v1/allocations/:id":
		return "allocation_get"
	case method == "PUT" && path == "/api/v1/allocations/:id":
		return "allocation_update"
	case method == "DELETE" && path == "/api/v1/allocations/:id":
		return "allocation_delete"
	case method == "POST" && path == "/api/v1/allocations/:id/start":
		return "allocation_start"
	case method == "POST" && path == "/api/v1/allocations/:id/stop":
		return "allocation_stop"
	case method == "GET" && path == "/api/v1/allocations/:id/status":
		return "allocation_status"
	case method == "GET" && path == "/api/v1/servers":
		return "server_list"
	case method == "POST" && path == "/api/v1/servers":
		return "server_create"
	case method == "GET" && path == "/api/v1/servers/:id":
		return "server_get"
	case method == "PUT" && path == "/api/v1/servers/:id":
		return "server_update"
	case method == "DELETE" && path == "/api/v1/servers/:id":
		return "server_delete"
	case method == "POST" && path == "/api/v1/servers/:id/power":
		return "server_power_control"
	case method == "GET" && path == "/api/v1/servers/:id/status":
		return "server_status"
	case method == "PUT" && path == "/api/v1/servers/:id/bios":
		return "server_bios_config"
	case method == "POST" && path == "/api/v1/servers/:id/firmware":
		return "server_firmware_upgrade"
	case method == "GET" && path == "/api/v1/servers/:id/gpus":
		return "server_gpus_get"
	case method == "POST" && path == "/api/v1/servers/:id/gpus":
		return "server_gpu_add"
	case method == "GET" && path == "/api/v1/servers/:id/config":
		return "server_config_get"
	case method == "PUT" && path == "/api/v1/servers/:id/config":
		return "server_config_update"
	case method == "GET" && path == "/api/v1/servers/:id/config/versions":
		return "server_config_versions"
	case method == "GET" && path == "/api/v1/events":
		return "event_list"
	case method == "POST" && path == "/api/v1/events":
		return "event_create"
	case method == "GET" && path == "/api/v1/events/:id":
		return "event_get"
	case method == "GET" && path == "/api/v1/alerts":
		return "alert_list"
	case method == "POST" && path == "/api/v1/alerts":
		return "alert_create"
	case method == "PUT" && path == "/api/v1/alerts/:id":
		return "alert_update"
	case method == "GET" && path == "/api/v1/alerts/correlation":
		return "alert_correlation"
	default:
		return "unknown_operation"
	}
}

// GetRequestID 从Context中获取请求ID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserID 从Context中获取用户ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetOperation 从Context中获取操作类型
func GetOperation(ctx context.Context) string {
	if operation, ok := ctx.Value(OperationKey).(string); ok {
		return operation
	}
	return ""
}
