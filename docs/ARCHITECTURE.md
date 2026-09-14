# Octopus Architecture / 架构文档

## Locked Decisions / 锁定决策

This document records the architectural decisions that have been locked for the Octopus workflow platform.

本文档记录 Octopus 工作流平台已锁定的架构决策。

---

### 1. Workflow Design / 工作流设计

- **Templated Simplified DAG**: Workflows are defined as simplified DAG structures with template support
- **Vue Flow**: DAG visualization will use Vue Flow library (planned for FE-2+)
- **n8n P0 UX Patterns**: UI/UX follows n8n's proven patterns (reference only, no source code usage)

**模板化简化 DAG**：工作流定义为支持模板的简化 DAG 结构
**Vue Flow**：DAG 可视化将使用 Vue Flow 库（计划于 FE-2+）
**n8n P0 UX 模式**：UI/UX 遵循 n8n 的成熟模式（仅参考，不使用源代码）

---

### 2. Execution Engine / 执行引擎

- **Go State Machine**: Workflow execution is driven by a Go-based state machine
- **Unified Multi-Language Worker Protocol**: Task execution uses a language-agnostic protocol
- **Python Default Worker**: Primary worker implementation in Python

**Go 状态机**：工作流执行由 Go 状态机驱动
**统一多语言 Worker 协议**：任务执行使用语言无关协议
**Python 默认 Worker**：主要 Worker 实现使用 Python

---

### 3. Storage / 存储

- **MySQL 8**: Primary relational database for workflow definitions and execution state
- **Redis 7**: Caching, task queues, and realtime coordination

**MySQL 8**：主关系数据库，存储工作流定义和执行状态
**Redis 7**：缓存、任务队列和实时协调

---

### 4. Communication / 通信

- **SSE (Server-Sent Events)**: Realtime updates from backend to frontend
- **Edge Outlets**: DAG edges support `success` and `failure` outcomes for conditional routing

**SSE（服务器发送事件）**：后端到前端的实时更新
**边缘出口**：DAG 边支持 `success` 和 `failure` 结果，用于条件路由

---

### 5. Development Conventions / 开发约定

- **go-zero Framework**: Backend uses go-zero microservice framework
- **goctl Spec-First**: API definitions in `.api` files, code generated via `goctl api go`
- **zeromicro/ai-context Conventions**: Follow go-zero's AI-friendly project conventions
- **Model Generation**: Database models via `goctl model mysql ddl` (without `-c` flag)

**go-zero 框架**：后端使用 go-zero 微服务框架
**goctl Spec-First**：API 定义在 `.api` 文件中，通过 `goctl api go` 生成代码
**zeromicro/ai-context 约定**：遵循 go-zero 的 AI 友好项目约定
**模型生成**：数据库模型通过 `goctl model mysql ddl` 生成（不使用 `-c` 标志）

---

### 6. Task Protocol / 任务协议

The unified task protocol enables multi-language worker implementations:

```json
{
  "task_id": "uuid",
  "task_type": "http_request|script|...",
  "payload": {},
  "context": {
    "workflow_id": "uuid",
    "execution_id": "uuid",
    "node_id": "string"
  }
}
```

Response:

```json
{
  "task_id": "uuid",
  "status": "success|failure|pending",
  "result": {},
  "error": null
}
```

统一任务协议支持多语言 Worker 实现（协议结构见上）

---

## Roadmap / 路线图

### Phase 1: Foundation / 基础
- [x] Project scaffold
- [ ] BE-1: Workflow Definition CRUD
- [ ] FE-1: Workflow List Page

### Phase 2: Execution / 执行
- [ ] BE-2: Execution Engine
- [ ] FE-2: DAG Canvas (Vue Flow)
- [ ] W-1: Python Worker

### Phase 3: Advanced / 高级
- [ ] Template System
- [ ] Multi-language Workers
- [ ] Monitoring & Observability
