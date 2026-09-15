# Octopus Local Development Setup / 本地开发环境搭建

This document provides complete instructions for running the Octopus workflow platform locally for development and testing.

本文档提供完整的 Octopus 工作流平台本地开发和测试说明。

---

## Quick Start / 快速开始

### One-Shot Script / 一键启动脚本

```bash
# Start all services in background
DETACH=true ./scripts/dev-up.sh all

# Or start in foreground (for debugging)
./scripts/dev-up.sh all
```

### Manual Setup / 手动设置

#### 1. Start Infrastructure / 启动基础设施

**Option A: Docker Compose (Recommended)**

```bash
docker compose up -d
```

**Option B: Local Services (Ubuntu/Debian)**

```bash
# Install MySQL and Redis
sudo apt-get update && sudo apt-get install -y mysql-server redis-server

# Start services
sudo service mysql start
sudo service redis-server start

# Configure database
sudo mysql -e "
  CREATE DATABASE IF NOT EXISTS octopus;
  CREATE USER IF NOT EXISTS 'octopus'@'localhost' IDENTIFIED BY 'octopus_pwd';
  GRANT ALL PRIVILEGES ON octopus.* TO 'octopus'@'localhost';
  ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'octopus_root_pwd';
  FLUSH PRIVILEGES;
"

# Apply schema
mysql -u root -poctopus_root_pwd octopus < docker/mysql/init/00_schema.sql
```

#### 2. Start Backend / 启动后端

```bash
cd backend
go mod tidy
go run workflow.go -f etc/workflow-api.yaml
```

#### 3. Start Worker / 启动 Worker

```bash
cd worker
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python3 main.py
```

---

## Service Ports & Endpoints / 服务端口

| Service | Port | Health Check |
|---------|------|--------------|
| MySQL | 3306 | `mysqladmin ping -h localhost -u root -poctopus_root_pwd` |
| Redis | 6379 | `redis-cli ping` |
| Backend API | 8888 | `curl http://localhost:8888/api/v1/health` |
| Worker | N/A | Check `/tmp/octopus-worker.log` |

---

## Configuration / 配置

### Backend (`backend/etc/workflow-api.yaml`)

```yaml
Name: workflow-api
Host: 0.0.0.0
Port: 8888

MySQL:
  DataSource: root:octopus_root_pwd@tcp(localhost:3306)/octopus?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Host: localhost:6379
  Type: node
```

### Worker Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | `redis://localhost:6379` | Redis connection URL |
| `TASK_QUEUE` | `octopus:tasks` | Redis queue for tasks |
| `RESULT_QUEUE` | `octopus:results` | Redis queue for results |
| `WORKER_ID` | `worker-{pid}` | Worker identifier |

---

## End-to-End Smoke Test / 端到端测试

### 1. Health Check / 健康检查

```bash
curl http://localhost:8888/api/v1/health
# Expected: {"status":"ok","version":"v1.0.0"}
```

### 2. Create Workflow / 创建工作流

```bash
curl -X POST http://localhost:8888/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "E2E Test Workflow",
    "description": "End-to-end test workflow",
    "nodes": [
      {"id": "node_1", "type": "script", "name": "Script Node", "position": {"x": 100, "y": 100}},
      {"id": "node_2", "type": "http", "name": "HTTP Node", "position": {"x": 300, "y": 100}}
    ],
    "edges": [
      {"id": "edge_1", "source": "node_1", "target": "node_2", "outlet": "success"}
    ],
    "entry_node_id": "node_1"
  }'
# Expected: {"id":1,"version":1,"status":"draft"}
```

### 3. Enable Workflow / 启用工作流

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/enable
# Expected: {"status":"enabled","version":1}
```

### 4. Trigger Run / 触发运行

```bash
curl -X POST http://localhost:8888/api/v1/workflows/1/runs
# Expected: {"run_id":1}
```

### 5. Check Run Status / 检查运行状态

```bash
curl http://localhost:8888/api/v1/runs/1
# Expected: "status": "succeeded" (after worker processes)
```

### 6. Check Run Steps / 检查运行步骤

```bash
curl http://localhost:8888/api/v1/runs/1/steps
# Expected: Both steps with "status": "succeeded"
```

### 7. SSE Events Test / SSE 事件测试

```bash
# Terminal 1: Subscribe to events (use timeout to auto-close)
timeout 10 curl -N http://localhost:8888/api/v1/runs/1/events

# Terminal 2: Trigger a new run
curl -X POST http://localhost:8888/api/v1/workflows/1/runs
```

**Expected SSE Output:**

```
event: step.status_changed
id: <uuid>
data: {"event_id":"...","event_type":"step.status_changed","run_id":2,"workflow_id":1,...}

event: run.status_changed
id: <uuid>
data: {"event_id":"...","event_type":"run.status_changed",...,"payload":{"from_status":"running","to_status":"succeeded"}}

event: run.terminal
id: <uuid>
data: {"event_id":"...","event_type":"run.terminal",...,"payload":{"final_status":"succeeded"}}
```

---

## Known Issues & Pitfalls / 已知问题

### 1. SSE Connection Error on Client Disconnect

**症状**: SSE 客户端断开时后端日志显示 `context deadline exceeded` 500 错误

**原因**: go-zero 的日志处理器在 SSE 长连接上会记录超时错误

**影响**: 仅影响日志美观，功能正常

**解决方案**: 可以通过自定义日志中间件过滤 SSE 路由的超时日志

### 2. MySQL Socket Permission (Linux Local Install)

**症状**: `Can't connect to local MySQL server through socket`

**解决方案**:
```bash
sudo chmod 755 /var/run/mysqld
sudo chown mysql:mysql /var/run/mysqld
sudo service mysql restart
```

### 3. Worker Standalone Mode

**症状**: Worker 日志显示 "Running in standalone mode"

**原因**: 无法连接到 Redis

**解决方案**: 确保 Redis 正在运行且 `REDIS_URL` 配置正确

### 4. Docker Compose Volume Permissions

**症状**: MySQL 无法启动，权限错误

**解决方案**:
```bash
docker compose down -v
docker compose up -d
```

---

## Troubleshooting / 故障排除

### Check Service Logs / 检查服务日志

```bash
# Backend logs
tail -f /tmp/octopus-backend.log

# Worker logs
tail -f /tmp/octopus-worker.log

# Docker MySQL logs
docker logs octopus-mysql

# Docker Redis logs
docker logs octopus-redis
```

### Reset Database / 重置数据库

```bash
# Docker
docker compose down -v
docker compose up -d

# Local
mysql -u root -poctopus_root_pwd -e "DROP DATABASE octopus; CREATE DATABASE octopus;"
mysql -u root -poctopus_root_pwd octopus < docker/mysql/init/00_schema.sql
```

### Verify Redis Queue / 验证 Redis 队列

```bash
# Check pending tasks
redis-cli LLEN octopus:tasks

# Check pending results
redis-cli LLEN octopus:results

# Monitor queue activity
redis-cli MONITOR
```

---

## API Reference Summary / API 参考

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/workflows` | Create workflow |
| GET | `/api/v1/workflows` | List workflows |
| GET | `/api/v1/workflows/:id` | Get workflow |
| PUT | `/api/v1/workflows/:id` | Update workflow |
| DELETE | `/api/v1/workflows/:id` | Delete workflow (draft only) |
| POST | `/api/v1/workflows/:id/validate` | Validate workflow |
| POST | `/api/v1/workflows/:id/enable` | Enable workflow |
| POST | `/api/v1/workflows/:id/disable` | Disable workflow |
| POST | `/api/v1/workflows/:id/runs` | Trigger run |
| GET | `/api/v1/runs/:runId` | Get run details |
| GET | `/api/v1/runs/:runId/steps` | Get run steps |
| GET | `/api/v1/runs/:runId/events` | SSE event stream |

---

## Verified Environment / 验证环境

This setup has been tested on:

- **OS**: Ubuntu 24.04 / Docker Desktop
- **Go**: 1.22+
- **Python**: 3.11+
- **MySQL**: 8.0
- **Redis**: 7.x

**Verification Date**: 2026-09-15
