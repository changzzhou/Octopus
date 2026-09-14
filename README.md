# Octopus - 通用工作流平台 / Universal Workflow Platform

[English](#english) | [中文](#中文)

---

## English

Octopus is a universal workflow automation platform designed for building, executing, and monitoring templated DAG-based workflows.

### Architecture Overview

- **Frontend**: Vue 3 + TypeScript + Tailwind CSS (Vue Flow for DAG visualization - planned)
- **Backend**: Go with go-zero framework (goctl Spec-First development)
- **Worker**: Python-based task executor with unified multi-language protocol
- **Storage**: MySQL 8 + Redis 7
- **Communication**: SSE for realtime updates

### Quick Start

#### Prerequisites

- Docker & Docker Compose
- Go 1.22+
- Node.js 18+ / pnpm
- Python 3.11+

#### Start Infrastructure

```bash
docker-compose up -d
```

#### Run Backend

```bash
cd backend
go mod tidy
go run workflow.go
```

#### Run Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

#### Run Worker

```bash
cd worker
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

### Project Structure

```
/
├── README.md                 # This file
├── .gitignore
├── docker-compose.yml        # MySQL 8 + Redis 7
├── docker/mysql/init/        # Database initialization scripts
├── docs/ARCHITECTURE.md      # Architecture decisions
├── frontend/                 # Vue 3 + TypeScript + Tailwind
├── backend/                  # Go go-zero API server
└── worker/                   # Python task worker
```

### Design Principles

- **Templated DAG**: Simplified DAG workflow definitions with template support
- **Spec-First API**: Backend APIs defined in `.api` files, generated via goctl
- **Unified Worker Protocol**: Language-agnostic task execution protocol (Python default)
- **Edge Outcomes**: DAG edges support `success` and `failure` outlets

### API Verification

After starting the infrastructure and backend, you can verify the API with curl:

#### Health Check

```bash
curl http://localhost:8888/api/v1/health
# {"status":"ok","version":"v1.0.0"}
```

#### Create Workflow

```bash
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First Workflow",
    "description": "A simple test workflow",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "Start", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "Call API", "position": {"x": 300, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"}
    ],
    "entry_node_id": "node_1"
  }'
# {"id":1,"version":1,"status":"draft"}
```

#### List Workflows

```bash
curl "http://localhost:8888/api/v1/workflows?page=1&page_size=10"
# {"total":1,"workflows":[...]}

# Filter by status
curl "http://localhost:8888/api/v1/workflows?status=draft"
```

#### Get Workflow Detail (with nodes/edges)

```bash
curl http://localhost:8888/api/v1/workflows/1
# Returns full workflow including nodes and edges
```

#### Update/Save Workflow

```bash
curl -X PUT http://localhost:8888/api/v1/workflows/1 \
  -H "Content-Type: application/json" \
  -d '{
    "version": 1,
    "name": "Updated Workflow",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "Start", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "Call API", "position": {"x": 300, "y": 100}},
      {"id": "node_3", "type": "human", "name": "Review", "position": {"x": 500, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"},
      {"id": "edge_2", "source": "node_2", "target": "node_3", "outlet": "success"}
    ]
  }'
# {"version":2}
```

#### Validate Workflow

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/validate
# {"valid":true,"errors":[]} or {"valid":false,"errors":[...]}
```

#### Enable Workflow (draft -> enabled)

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# {"status":"enabled","version":2}
```

#### Disable Workflow (enabled -> disabled)

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/disable
# {"status":"disabled"}
```

#### Delete Workflow (draft only)

```bash
curl -X DELETE http://localhost:8888/api/v1/workflows/1
# {} (empty response on success)
# Note: Can only delete workflows in "draft" status
```

---

## 中文

Octopus 是一个通用工作流自动化平台，用于构建、执行和监控基于模板化 DAG 的工作流。

### 架构概览

- **前端**: Vue 3 + TypeScript + Tailwind CSS（Vue Flow DAG 可视化 - 计划中）
- **后端**: Go + go-zero 框架（goctl Spec-First 开发）
- **Worker**: Python 任务执行器，支持统一的多语言协议
- **存储**: MySQL 8 + Redis 7
- **通信**: SSE 实时更新

### 快速开始

#### 前置条件

- Docker & Docker Compose
- Go 1.22+
- Node.js 18+ / pnpm
- Python 3.11+

#### 启动基础设施

```bash
docker-compose up -d
```

#### 运行后端

```bash
cd backend
go mod tidy
go run workflow.go
```

#### 运行前端

```bash
cd frontend
pnpm install
pnpm dev
```

#### 运行 Worker

```bash
cd worker
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

### 项目结构

```
/
├── README.md                 # 本文件
├── .gitignore
├── docker-compose.yml        # MySQL 8 + Redis 7
├── docker/mysql/init/        # 数据库初始化脚本
├── docs/ARCHITECTURE.md      # 架构决策文档
├── frontend/                 # Vue 3 + TypeScript + Tailwind
├── backend/                  # Go go-zero API 服务
└── worker/                   # Python 任务 Worker
```

### 设计原则

- **模板化 DAG**: 简化的 DAG 工作流定义，支持模板
- **Spec-First API**: 后端 API 在 `.api` 文件中定义，通过 goctl 生成
- **统一 Worker 协议**: 语言无关的任务执行协议（默认 Python）
- **边缘结果**: DAG 边支持 `success` 和 `failure` 出口

### API 验证

启动基础设施和后端后，可以使用 curl 验证 API：

#### 健康检查

```bash
curl http://localhost:8888/api/v1/health
# {"status":"ok","version":"v1.0.0"}
```

#### 创建工作流

```bash
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的第一个工作流",
    "description": "一个简单的测试工作流",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "开始", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "调用API", "position": {"x": 300, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"}
    ],
    "entry_node_id": "node_1"
  }'
# {"id":1,"version":1,"status":"draft"}
```

#### 列出工作流

```bash
curl "http://localhost:8888/api/v1/workflows?page=1&page_size=10"
# {"total":1,"workflows":[...]}

# 按状态过滤
curl "http://localhost:8888/api/v1/workflows?status=draft"
```

#### 获取工作流详情（包含 nodes/edges）

```bash
curl http://localhost:8888/api/v1/workflows/1
# 返回完整的工作流，包括 nodes 和 edges
```

#### 更新/保存工作流

```bash
curl -X PUT http://localhost:8888/api/v1/workflows/1 \
  -H "Content-Type: application/json" \
  -d '{
    "version": 1,
    "name": "更新后的工作流",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "开始", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "调用API", "position": {"x": 300, "y": 100}},
      {"id": "node_3", "type": "human", "name": "审核", "position": {"x": 500, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"},
      {"id": "edge_2", "source": "node_2", "target": "node_3", "outlet": "success"}
    ]
  }'
# {"version":2}
```

#### 验证工作流

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/validate
# {"valid":true,"errors":[]} 或 {"valid":false,"errors":[...]}
```

#### 启用工作流（draft -> enabled）

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# {"status":"enabled","version":2}
```

#### 禁用工作流（enabled -> disabled）

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/disable
# {"status":"disabled"}
```

#### 删除工作流（仅限草稿）

```bash
curl -X DELETE http://localhost:8888/api/v1/workflows/1
# {} (成功时返回空响应)
# 注意：只能删除状态为 "draft" 的工作流
```

---

## License

MIT
