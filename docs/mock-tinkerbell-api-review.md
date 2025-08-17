# Mock Tinkerbell API 一致性审查报告

## 📋 审查概述

本报告对比了当前 Mock Tinkerbell 实现与实际 Tinkerbell 项目的 API 一致性，并提供了修正建议。

## 🔍 发现的问题

### 1. **API 协议不匹配**

| 方面 | 实际 Tinkerbell | 当前 Mock 实现 | 问题等级 |
|------|----------------|----------------|----------|
| 主要协议 | **gRPC** (端口 42113) | **REST** (端口 8080) | 🔴 严重 |
| 辅助协议 | HTTP (端口 42114) | 无 | 🟡 中等 |
| 数据格式 | Protocol Buffers | JSON | 🔴 严重 |

### 2. **核心概念差异**

| 概念 | 实际 Tinkerbell | 当前 Mock 实现 | 差异说明 |
|------|----------------|----------------|----------|
| Hardware | 简单的硬件记录 | 复杂的 GPU 服务器模型 | 数据结构不匹配 |
| Template | 包含 Action 列表 | 简化的步骤定义 | 缺少 Action 概念 |
| Workflow | 基于 Template 创建 | 直接创建工作流 | 流程不一致 |
| Action | 独立的工作单元 | 内嵌在工作流中 | 架构差异 |

### 3. **服务架构差异**

**实际 Tinkerbell 架构：**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Tink Server   │    │   Tink Worker   │    │   Kubernetes    │
│  (gRPC + HTTP)  │◄──►│     (gRPC)      │    │    Backend      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      Smee       │    │     Hegel       │    │      PBnJ       │
│   (DHCP)        │    │  (Metadata)     │    │     (BMC)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**当前 Mock 实现：**
```
┌─────────────────┐
│  Mock Server    │
│   (REST only)   │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  In-Memory      │
│   Storage       │
└─────────────────┘
```

## 🛠️ 修正方案

### 1. **添加 gRPC 支持**

已创建 `protos/workflow.proto` 文件，定义了符合 Tinkerbell 的 protobuf 接口：

```protobuf
service WorkflowService {
  rpc CreateWorkflow(CreateWorkflowRequest) returns (CreateWorkflowResponse);
  rpc GetWorkflow(GetWorkflowRequest) returns (GetWorkflowResponse);
  rpc ListWorkflows(ListWorkflowsRequest) returns (ListWorkflowsResponse);
  rpc DeleteWorkflow(DeleteWorkflowRequest) returns (DeleteWorkflowResponse);
  rpc StartWorkflow(StartWorkflowRequest) returns (StartWorkflowResponse);
  rpc StopWorkflow(StopWorkflowRequest) returns (StopWorkflowResponse);
}

service HardwareService {
  rpc CreateHardware(CreateHardwareRequest) returns (CreateHardwareResponse);
  rpc GetHardware(GetHardwareRequest) returns (GetHardwareResponse);
  rpc ListHardware(ListHardwareRequest) returns (ListHardwareResponse);
  rpc UpdateHardware(UpdateHardwareRequest) returns (UpdateHardwareResponse);
  rpc DeleteHardware(DeleteHardwareRequest) returns (DeleteHardwareResponse);
}

service TemplateService {
  rpc CreateTemplate(CreateTemplateRequest) returns (CreateTemplateResponse);
  rpc GetTemplate(GetTemplateRequest) returns (GetTemplateResponse);
  rpc ListTemplates(ListTemplatesRequest) returns (ListTemplatesResponse);
  rpc UpdateTemplate(UpdateTemplateRequest) returns (UpdateTemplateResponse);
  rpc DeleteTemplate(DeleteTemplateRequest) returns (DeleteTemplateResponse);
}
```

### 2. **修正数据结构**

**Hardware 数据结构：**
```go
// 实际 Tinkerbell 格式
type Hardware struct {
    ID                string    `json:"id"`
    NetworkInterfaces string    `json:"network_interfaces"` // JSON 字符串
    Metadata          string    `json:"metadata"`           // JSON 字符串
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

// 当前 Mock 实现（需要修正）
type Hardware struct {
    ID           string        `json:"id"`
    Name         string        `json:"name"`
    Type         string        `json:"type"`
    Status       string        `json:"status"`
    IPAddress    string        `json:"ip_address"`
    MACAddress   string        `json:"mac_address"`
    SerialNumber string        `json:"serial_number"`
    Model        string        `json:"model"`
    Specs        HardwareSpecs `json:"specs"` // 这个字段在实际 Tinkerbell 中不存在
    CreatedAt    time.Time     `json:"created_at"`
    UpdatedAt    time.Time     `json:"updated_at"`
}
```

**Workflow 数据结构：**
```go
// 实际 Tinkerbell 格式
type Workflow struct {
    ID         string        `json:"id"`
    TemplateID string        `json:"template_id"`
    HardwareID string        `json:"hardware_id"`
    State      WorkflowState `json:"state"`
    Actions    []Action      `json:"actions"`
    CreatedAt  time.Time     `json:"created_at"`
    UpdatedAt  time.Time     `json:"updated_at"`
}

type Action struct {
    ID          string      `json:"id"`
    Name        string      `json:"name"`
    Image       string      `json:"image"`
    Timeout     int32       `json:"timeout"`
    Environment []string    `json:"environment"`
    Volumes     []string    `json:"volumes"`
    State       ActionState `json:"state"`
    StartedAt   *time.Time  `json:"started_at,omitempty"`
    CompletedAt *time.Time  `json:"completed_at,omitempty"`
}
```

### 3. **端口配置修正**

更新配置文件以匹配实际 Tinkerbell：

```yaml
server:
  port: 8080          # HTTP API (辅助)
  host: "0.0.0.0"
  grpc_port: 42113    # gRPC API (主要)
  grpc_host: "0.0.0.0"
```

## 📊 兼容性评估

### 当前兼容性状态

| 功能模块 | 兼容性 | 说明 |
|----------|--------|------|
| 硬件管理 | 🟡 部分兼容 | 数据结构需要调整 |
| 工作流管理 | 🔴 不兼容 | 缺少 Action 概念 |
| 模板管理 | 🟡 部分兼容 | 需要添加 Action 支持 |
| API 协议 | 🔴 不兼容 | 需要添加 gRPC 支持 |

### 修正优先级

1. **🔴 高优先级**：添加 gRPC 支持
2. **🔴 高优先级**：修正数据结构
3. **🟡 中优先级**：添加 Action 概念
4. **🟢 低优先级**：保持 REST API 作为辅助

## 🚀 实施建议

### 阶段一：基础兼容性（1-2周）
1. 添加 gRPC 服务器
2. 实现基本的 protobuf 接口
3. 修正 Hardware 数据结构

### 阶段二：功能完善（2-3周）
1. 实现 Action 概念
2. 完善 Workflow 和 Template 结构
3. 添加工作流执行引擎

### 阶段三：高级功能（3-4周）
1. 实现 Kubernetes 后端模拟
2. 添加事件流支持
3. 完善监控和日志

## 📝 总结

当前的 Mock Tinkerbell 实现虽然提供了基本的模拟功能，但与实际 Tinkerbell API 存在显著差异。主要问题包括：

1. **协议不匹配**：缺少 gRPC 支持
2. **数据结构差异**：Hardware 和 Workflow 模型不匹配
3. **概念缺失**：缺少 Action 等核心概念

建议按照上述修正方案进行改进，以确保 Mock 服务能够更好地模拟真实的 Tinkerbell 环境，为您的 IaaS 系统开发提供更准确的测试环境。
