package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"mock-tinkerbell/api"
	"mock-tinkerbell/config"
	"mock-tinkerbell/services"
)

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	
	// 加载配置
	config.LoadConfig()
	cfg := &config.AppConfig
	
	fmt.Println("🚀 启动 Mock Tinkerbell 服务...")
	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("   - 服务地址: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("   - 模拟延迟: %dms\n", cfg.Hardware.MockDelay)
	fmt.Printf("   - 工作流延迟: %dms\n", cfg.Workflow.StepDelay)
	
	// 创建服务实例
	hardwareService := services.NewHardwareService(cfg)
	workflowService := services.NewWorkflowService(cfg)
	
	// 设置路由
	router := api.SetupRoutes(hardwareService, workflowService)
	
	// 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("🌐 服务启动在: http://%s\n", addr)
	fmt.Printf("📖 API文档: http://%s\n", addr)
	fmt.Printf("💚 健康检查: http://%s/health\n", addr)
	
	log.Fatal(router.Run(addr))
}
