# Octopus Frontend

Vue 3 + TypeScript + Vite frontend for the Octopus Universal Workflow Platform.

## Quick Start

### Prerequisites
- Node.js 18+
- pnpm (recommended) or npm

### Installation

```bash
cd frontend
pnpm install
```

### Development

**Option 1: With Backend (Real API)**

First, start the backend service (see `/backend/README.md` or root `docker-compose.yml`), then:

```bash
pnpm dev
```

The dev server runs at http://localhost:5173 and proxies `/api` requests to `http://localhost:8888`.

**Option 2: Mock Mode (No Backend)**

For frontend-only development without running the backend:

```bash
# Create .env.local file
echo "VITE_MOCK_API=true" > .env.local
pnpm dev
```

Or enable mock mode at runtime in browser console:
```javascript
window.__OCTOPUS_MOCK_API__ = true
```

Then refresh the page.

### Build

```bash
pnpm build
```

## Project Structure

```
frontend/
├── src/
│   ├── api/           # API client with mock support
│   ├── router/        # Vue Router configuration
│   ├── views/         # Page components
│   │   ├── HomeView.vue           # Landing page
│   │   ├── WorkflowsView.vue      # Workflow list
│   │   ├── WorkflowDetailView.vue # Workflow detail
│   │   ├── WorkflowDesignView.vue # Design-time placeholder
│   │   └── WorkflowExecutionView.vue # Execution-time placeholder
│   ├── App.vue
│   ├── main.ts
│   └── style.css
├── .env.example       # Environment config template
└── vite.config.ts
```

## Routes

| Path | Description |
|------|-------------|
| `/` | Home / Landing page |
| `/workflows` | Workflow list |
| `/workflows/:id` | Workflow detail |
| `/workflows/:id/design` | Design-time canvas (placeholder) |
| `/workflows/:id/executions/:executionId` | Execution-time view (placeholder) |

## API Mode Configuration

### Real Backend Mode (Default)

When `VITE_MOCK_API` is not set or set to `false`, the frontend connects to the real backend API.

**Prerequisites:**
1. Backend running on `http://localhost:8888` (see `/backend/README.md` or `docker-compose.yml`)
2. MySQL and Redis configured for the backend

**To run with real backend:**
```bash
# Ensure no mock override exists
rm -f .env.local
# Or explicitly disable mock
echo "VITE_MOCK_API=false" > .env.local

pnpm dev
```

**Features when connected to real backend:**
- Workflow list, create, update, delete via REST API
- Design canvas save/load with version control
- Trigger run / trial run
- Execution view with **SSE-first** realtime updates (auto fallback to polling if SSE fails)

### Mock Mode (No Backend)

For frontend-only development without running the backend:

- **Environment variable**: Set `VITE_MOCK_API=true` in `.env.local`
- **Runtime toggle**: Set `window.__OCTOPUS_MOCK_API__ = true` in browser console

```bash
echo "VITE_MOCK_API=true" > .env.local
pnpm dev
```

Mock data includes sample workflows with different statuses (draft/enabled/disabled).

### Execution View Realtime Updates

The execution view (`/workflows/:id/executions/:executionId`) uses a **SSE-first with polling fallback** strategy:

1. **SSE Mode**: Connects to `GET /api/v1/runs/:runId/events` for realtime updates
2. **Polling Fallback**: If SSE connection fails, automatically switches to polling every 2s
3. **State Refresh**: On each event, fetches full run/steps state for consistency

The UI shows the current connection mode: "Live (SSE)" or "Live (Polling)".

## Tech Stack

- Vue 3 with Composition API
- TypeScript
- Vite
- Vue Router 4
- Tailwind CSS 4
