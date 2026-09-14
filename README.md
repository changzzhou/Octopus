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

---

## License

MIT
