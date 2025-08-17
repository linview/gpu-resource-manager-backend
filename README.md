# GPU资源管理系统

基于事件驱动架构的GPU资源管理后台系统，支持裸金属GPU服务器的自动化部署和管理。

## 项目概述

本项目是一个企业级GPU基础设施管理平台，采用**模式A: IaaS作为资源编排层**的架构设计，实现GPU资源的全生命周期管理。

### 核心特性

- 🚀 **事件驱动架构**: 基于NATS JetStream的实时事件处理
- 🔧 **硬件自动化**: 集成Tinkerbell实现裸金属服务器自动化部署
- 📊 **智能监控**: 实时GPU状态监控和告警关联分析
- 🔄 **云原生**: 基于Kubernetes的容器化部署
- 🛡️ **安全可靠**: 支持多种硬件管理协议(IPMI/Redfish/gNMI)
- ⚡ **高性能**: 基于Echo v4.11.4的高性能Web框架

### 技术栈

- **主要语言**: Golang 1.21+ (核心微服务)
- **辅助语言**: Python (工具开发)
- **Web框架**: Echo v4.11.4 (MVP选定)
- **容器编排**: Kubernetes (k0s)
- **硬件管理**: Tinkerbell v0.12.2
- **操作系统**: Talos Linux
- **配置管理**: Cloud-init
- **数据存储**: PostgreSQL + Redis + InfluxDB
- **消息队列**: NATS JetStream (MVP) / Kafka (生产环境)
- **API接口**: REST (MVP) + gRPC/GraphQL (P1)
- **测试框架**: Testcontainers-Go + k8s

## 架构设计

### 架构模式

系统采用**IaaS作为资源编排层**的架构设计：

```
┌─────────────────────────────────────────────────────────────┐
│                   应用层 (Application Layer)                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │  用户界面   │ │  API网关    │ │  IaaS服务   │           │
│  │             │ │             │ │             │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  编排层 (Orchestration Layer)               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Kubernetes  │ │ Tinkerbell  │ │   CRD资源   │           │
│  │   (k0s)     │ │             │ │             │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  硬件层 (Hardware Layer)                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   IPMI      │ │  Redfish    │ │   gNMI      │           │
│  │   协议      │ │   协议      │ │   协议      │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│               基础设施层 (Infrastructure Layer)             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Talos Linux │ │    k0s      │ │ Cloud-init  │           │
│  │             │ │   集群      │ │             │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### 事件驱动架构

系统采用事件驱动架构，通过NATS JetStream实现组件间的松耦合通信：

- **硬件发现事件**: 自动发现和注册GPU服务器
- **状态变化事件**: 实时监控GPU状态变化
- **告警事件**: 智能告警关联分析
- **配置变更事件**: 自动化配置管理

## 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- Kubernetes集群 (可选)
- 硬件服务器 (支持IPMI/Redfish)

### 本地开发

1. **克隆项目**
```bash
git clone https://github.com/linview/gpu-resource-manager-backend.git
cd gpu-resource-manager-backend
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件配置数据库等连接信息
```

4. **启动依赖服务**
```bash
cd deployments/docker
docker-compose up -d postgres redis nats
```

5. **运行应用**
```bash
go run cmd/server/main.go
```

### Docker部署

1. **构建镜像**
```bash
docker build -t palebluedot-backend:latest .
```

2. **启动完整服务栈**
```bash
cd deployments/docker
docker-compose up -d
```

### Kubernetes部署

1. **部署Tinkerbell**
```bash
kubectl apply -f deployments/tinkerbell/tinkerbell-deployment.yaml
```

2. **部署GPU管理系统**
```bash
kubectl apply -f deployments/kubernetes/k0s-cluster.yaml
```

## API文档

### 核心API接口

#### GPU管理
- `GET /api/v1/gpus` - 获取GPU列表
- `GET /api/v1/gpus/{id}` - 获取GPU详情
- `PUT /api/v1/gpus/{id}` - 更新GPU信息
- `GET /api/v1/gpus/{id}/status` - 获取GPU状态
- `GET /api/v1/gpus/{id}/metrics` - 获取GPU指标
- `POST /api/v1/gpus/{id}/config` - 更新GPU配置
- `GET /api/v1/gpus/{id}/config` - 获取GPU配置
- `GET /api/v1/gpus/{id}/config/versions` - 获取配置版本

#### 服务器管理
- `GET /api/v1/servers` - 获取服务器列表
- `POST /api/v1/servers` - 创建服务器
- `GET /api/v1/servers/{id}` - 获取服务器详情
- `PUT /api/v1/servers/{id}` - 更新服务器信息
- `DELETE /api/v1/servers/{id}` - 删除服务器
- `POST /api/v1/servers/{id}/power` - 电源控制
- `GET /api/v1/servers/{id}/status` - 获取服务器状态
- `PUT /api/v1/servers/{id}/bios` - 配置BIOS
- `POST /api/v1/servers/{id}/firmware` - 升级固件
- `GET /api/v1/servers/{id}/gpus` - 获取服务器GPU
- `POST /api/v1/servers/{id}/gpus` - 添加GPU到服务器
- `GET /api/v1/servers/{id}/config` - 获取服务器配置
- `PUT /api/v1/servers/{id}/config` - 更新服务器配置
- `GET /api/v1/servers/{id}/config/versions` - 获取配置版本

#### 资源分配
- `GET /api/v1/allocations` - 获取分配列表
- `POST /api/v1/allocations` - 创建资源分配
- `GET /api/v1/allocations/{id}` - 获取分配详情
- `PUT /api/v1/allocations/{id}` - 更新分配
- `DELETE /api/v1/allocations/{id}` - 释放资源分配
- `POST /api/v1/allocations/{id}/start` - 启动分配
- `POST /api/v1/allocations/{id}/stop` - 停止分配
- `GET /api/v1/allocations/{id}/status` - 获取分配状态

#### 工作流管理
- `GET /api/v1/workflows` - 获取工作流列表
- `POST /api/v1/workflows` - 创建工作流
- `GET /api/v1/workflows/{id}` - 获取工作流详情
- `DELETE /api/v1/workflows/{id}` - 删除工作流
- `POST /api/v1/workflows/deploy` - 创建部署工作流
- `POST /api/v1/workflows/cleanup` - 创建清理工作流

#### 事件管理
- `GET /api/v1/events` - 获取事件列表
- `POST /api/v1/events` - 创建事件
- `GET /api/v1/events/{id}` - 获取事件详情

#### 告警管理
- `GET /api/v1/alerts` - 获取告警列表
- `POST /api/v1/alerts` - 创建告警
- `PUT /api/v1/alerts/{id}` - 更新告警
- `GET /api/v1/alerts/correlation` - 获取告警关联分析

### 示例请求

#### 创建GPU分配
```bash
curl -X POST http://localhost:8080/api/v1/allocations \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "server_id": "server456",
    "purpose": "AI训练"
  }'
```

#### 获取GPU状态
```bash
curl http://localhost:8080/api/v1/gpus/gpu-001/status
```

#### 更新GPU配置
```bash
curl -X POST http://localhost:8080/api/v1/gpus/gpu-001/config \
  -H "Content-Type: application/json" \
  -d '{
    "power_limit": 200,
    "memory_clock": 8000
  }'
```

### API响应格式

所有API响应都包含请求追踪信息：

```json
{
  "data": [...],
  "total": 1,
  "request_id": "38c23617-7afc-2547-8aaf-be8b6e6c1e67",
  "operation": "gpu_list"
}
```

## 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | 8080 | 服务端口 |
| `DB_HOST` | localhost | 数据库主机 |
| `DB_PORT` | 5432 | 数据库端口 |
| `DB_USER` | postgres | 数据库用户 |
| `DB_PASSWORD` | password | 数据库密码 |
| `REDIS_HOST` | localhost | Redis主机 |
| `NATS_URL` | nats://localhost:4222 | NATS连接地址 |
| `TINKERBELL_URL` | http://localhost:50061 | Tinkerbell API地址 |

### Tinkerbell配置

Tinkerbell配置文件位于 `deployments/tinkerbell/tinkerbell-config.yaml`，包含：

- 硬件管理协议配置 (IPMI/Redfish/gNMI)
- 操作系统模板 (Talos Linux/Ubuntu)
- GPU驱动安装脚本
- 网络和存储配置

## 监控和告警

### 监控指标

- **GPU指标**: 利用率、温度、功耗、内存使用
- **系统指标**: CPU、内存、网络、磁盘
- **业务指标**: 分配成功率、响应时间、错误率

### 告警规则

- GPU温度超过85°C
- GPU利用率超过95%
- 内存使用率超过90%
- 系统错误率超过5%

### 访问监控界面

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093

## 开发指南

### 项目结构

```
palebluedot-backend/
├── cmd/server/           # 应用入口
├── internal/             # 内部包
│   ├── api/             # API层
│   │   ├── handlers/    # 处理器
│   │   ├── middleware/  # 中间件
│   │   └── routes/      # 路由配置
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   ├── services/        # 业务服务
│   │   └── event/       # 事件服务
│   └── repository/      # 数据访问层
├── pkg/                 # 公共包
│   └── logger/          # 日志组件
├── deployments/         # 部署配置
├── docs/               # 文档
│   ├── analysis/       # 技术分析
│   ├── design/         # 架构设计
│   └── review/         # 代码评审
└── tests/              # 测试
```

### 开发规范

- 遵循Go官方代码规范
- 使用结构化日志记录
- 编写单元测试和集成测试
- 使用语义化版本控制
- 遵循Go语言Context最佳实践

### 测试

```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test -tags=integration ./...

# 运行性能测试
go test -bench=. ./...
```

## 部署架构

### Docker Compose部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                   反向代理层 (Nginx)                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   HTTP      │ │   HTTPS     │ │   负载均衡   │           │
│  │   (80)      │ │   (443)     │ │             │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  应用服务层                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ GPU管理API  │ │ Tinkerbell  │ │ Alertmanager│           │
│  │ (单实例)    │ │ (单实例)    │ │ (单实例)    │           │
│  │   :8080     │ │  :50061     │ │   :9093     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  数据存储层                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ PostgreSQL  │ │ Redis       │ │ InfluxDB    │           │
│  │ (单实例)    │ │ (单实例)    │ │ (单实例)    │           │
│  │   :5432     │ │   :6379     │ │   :8086     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  消息和监控层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ NATS        │ │ Prometheus  │ │ Grafana     │           │
│  │ (单实例)    │ │ (单实例)    │ │ (单实例)    │           │
│  │ :4222/8222  │ │   :9090     │ │   :3000     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### 部署配置

- **API服务**: 单实例部署，端口8080
- **数据库**: PostgreSQL单实例，端口5432
- **缓存**: Redis单实例，端口6379
- **消息队列**: NATS单实例，端口4222/8222
- **监控**: Prometheus + Grafana + Alertmanager
- **反向代理**: Nginx，端口80/443
- **硬件管理**: Tinkerbell，端口50061
- **时序数据库**: InfluxDB，端口8086

## 故障排除

### 常见问题

1. **Tinkerbell连接失败**
   - 检查网络连通性
   - 验证认证信息
   - 查看Tinkerbell日志

2. **GPU发现失败**
   - 检查硬件管理协议配置
   - 验证网络连接
   - 查看硬件日志

3. **事件处理异常**
   - 检查NATS连接状态
   - 验证事件格式
   - 查看事件处理日志

### 日志查看

```bash
# 查看应用日志
docker logs gpu-management-api

# 查看Tinkerbell日志
docker logs gpu-management-tinkerbell

# 查看数据库日志
docker logs gpu-management-postgres
```

## 最近更新

### v2.0.0 (2025-08-17)
- ✅ **Echo框架选定**: 使用Echo v4.11.4作为Web框架
- ✅ **事件总线优化**: EventBus支持Context参数
- ✅ **中间件增强**: 新增Context中间件，支持请求追踪
- ✅ **错误处理完善**: 区分Context错误和业务错误
- ✅ **超时控制**: 多层次超时控制机制
- ✅ **代码质量提升**: 符合Go语言Context最佳实践

### 技术改进
- **性能优化**: Echo框架零内存分配路由
- **可维护性**: 统一的Context使用模式
- **可扩展性**: 灵活的中间件架构
- **可靠性**: 完整的错误处理和资源管理

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码变更（遵循Go语言最佳实践）
4. 创建Pull Request

## 许可证

本项目采用Apache 2.0许可证，详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目主页: https://github.com/linview/gpu-resource-manager-backend
- 问题反馈: https://github.com/linview/gpu-resource-manager-backend/issues
- 文档: https://github.com/linview/gpu-resource-manager-backend/docs
- 架构设计: [docs/design/architecture-design.md](docs/design/architecture-design.md)
- Context重构: [docs/review/context-refactoring-summary.md](docs/review/context-refactoring-summary.md)
