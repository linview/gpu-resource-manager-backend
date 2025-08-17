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

### 技术栈

- **后端**: Golang + Echo v4.11.4框架
- **容器编排**: Kubernetes (k0s)
- **硬件管理**: Tinkerbell v0.12.2
- **操作系统**: Talos Linux
- **消息队列**: NATS JetStream (MVP) / Kafka (生产环境)
- **数据库**: PostgreSQL + Redis + InfluxDB
- **监控**: Prometheus + Grafana

## 架构设计

### 整体架构

```
User → IaaS → K8s → Tinkerbell → Hardware
```

- **应用层**: 用户界面、API网关、IaaS服务
- **编排层**: Kubernetes集群、Tinkerbell、CRD资源
- **硬件层**: IPMI、Redfish、gNMI协议
- **基础设施层**: Talos Linux、k0s集群、Cloud-init

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
git clone https://github.com/your-org/palebluedot-backend.git
cd palebluedot-backend
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

#### 服务器管理
- `GET /api/v1/servers` - 获取服务器列表
- `POST /api/v1/servers` - 创建服务器
- `POST /api/v1/servers/{id}/power` - 电源控制
- `GET /api/v1/servers/{id}/status` - 获取服务器状态

#### 资源分配
- `GET /api/v1/allocations` - 获取分配列表
- `POST /api/v1/allocations` - 创建资源分配
- `DELETE /api/v1/allocations/{id}` - 释放资源分配

#### 工作流管理
- `GET /api/v1/workflows` - 获取工作流列表
- `POST /api/v1/workflows/deploy` - 创建部署工作流
- `POST /api/v1/workflows/cleanup` - 创建清理工作流

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
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   ├── services/        # 业务服务
│   └── repository/      # 数据访问层
├── pkg/                 # 公共包
├── deployments/         # 部署配置
├── docs/               # 文档
└── tests/              # 测试
```

### 开发规范

- 遵循Go官方代码规范
- 使用结构化日志记录
- 编写单元测试和集成测试
- 使用语义化版本控制

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

### 生产环境架构

```
┌─────────────────────────────────────────────────────────────┐
│                   负载均衡器 (HAProxy/Nginx)                  │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Kubernetes集群 (k0s)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ GPU管理API  │ │ Tinkerbell  │ │ 监控服务    │           │
│  │ (多副本)    │ │ (单副本)    │ │ (多副本)    │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    数据层                                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ PostgreSQL  │ │ Redis       │ │ InfluxDB    │           │
│  │ (主从)      │ │ (集群)      │ │ (集群)      │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### 高可用配置

- **API服务**: 多副本部署，支持自动扩缩容
- **数据库**: PostgreSQL主从复制
- **缓存**: Redis集群模式
- **消息队列**: NATS集群模式
- **监控**: Prometheus联邦模式

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

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码变更
4. 创建Pull Request

## 许可证

本项目采用Apache 2.0许可证，详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目主页: https://github.com/your-org/palebluedot-backend
- 问题反馈: https://github.com/your-org/palebluedot-backend/issues
- 文档: https://github.com/your-org/palebluedot-backend/docs
