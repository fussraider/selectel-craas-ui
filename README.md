# Selectel CRaaS Web Interface

A modern, full-stack web interface for managing Selectel Container Registry (CRaaS) repositories and images.

## Features

- **Project & Registry Management**: Browse projects and registries within your Selectel account.
- **Repository Insights**: View repositories and their details.
- **Image Management**: List images with detailed metadata (tags, size, creation date).
- **Bulk Cleanup**: Select multiple images to delete at once.
- **Garbage Collection Control**: Option to trigger Garbage Collection (GC) immediately upon deletion.
- **Destructive Action Guards**:
    - Confirmation modals for deleting registries and repositories require typing the resource name for verification.
    - Single image deletion now supports the "Run GC" option via a unified confirmation dialog.
- **Configuration Control**: Environment-based feature flags to disable destructive actions (registry, repository, or
  image deletion).
- **Optimistic UI Updates**: Immediate feedback on deletion actions without waiting for full list re-fetching.
- **Responsive UI**: Built with Vue 3 and modern SCSS for a clean dark-mode experience, with responsive sidebar and
  tooltips for disabled actions.

## Tech Stack

- **Backend**: Go 1.24+ (Chi router, Selectel SDK, Slog logging).
- **Frontend**: Vue 3, TypeScript, Pinia (Setup Stores), Vite.
- **Testing**: Playwright for frontend verification.

## Prerequisites

- Go 1.24+
- Node.js 18+
- A Selectel account with API credentials.

## Configuration

The application is configured via environment variables. You can set these directly or use a `.env` file in the
`backend/` directory.

### Core Configuration

| Variable                | Description                      | Default    |
|:------------------------|:---------------------------------|:-----------|
| `WEB_PORT`              | Port for the backend server      | `8080`     |
| `SELECTEL_USERNAME`     | Selectel Username                | (Required) |
| `SELECTEL_ACCOUNT_ID`   | Selectel Account ID              | (Required) |
| `SELECTEL_PASSWORD`     | Selectel Password                | (Required) |
| `SELECTEL_PROJECT_NAME` | Selectel Project Name (Required) | (Required) |

### Feature Flags (Destructive Actions)

Control which delete operations are permitted. If disabled, the corresponding buttons in the UI will be inactive with a
tooltip explaining the restriction.

| Variable                   | Description                               | Default |
|:---------------------------|:------------------------------------------|:--------|
| `ENABLE_DELETE_REGISTRY`   | Allow deletion of entire registries       | `false` |
| `ENABLE_DELETE_REPOSITORY` | Allow deletion of repositories            | `false` |
| `ENABLE_DELETE_IMAGE`      | Allow deletion of images (single or bulk) | `false` |

### Logging

| Variable     | Description                                       | Default |
|:-------------|:--------------------------------------------------|:--------|
| `LOG_LEVEL`  | Logging level ( `debug`, `info`, `warn`, `error`) | `INFO`  |
| `LOG_FORMAT` | Log format (`text`, `json`)                       | `TEXT`  |

## Running the Application

### Development Mode

You can run backend and frontend separately for development.

1. **Backend**:
   ```bash
   cd backend
   # Install dependencies
   go mod download
   # Run the server
   go run cmd/server/main.go
   ```
   The server will start on port 8080 (or as configured). It supports graceful shutdown (Ctrl+C).

2. **Frontend**:
   ```bash
   cd frontend
   # Install dependencies
   npm install
   # Start dev server
   npm run dev
   ```
   The frontend will start on http://localhost:5173 and proxy `/api` requests to `http://localhost:8080`.

### Production Build

You can use the `Makefile` in the root directory:

```bash
make all
```

To run the built artifacts:

```bash
make run-backend
# In a separate terminal
make run-frontend
```

## Project Structure

- `backend/`: Go backend source code.
    - `cmd/server`: Entry point.
    - `internal/auth`: Selectel Keystone authentication.
    - `internal/config`: Configuration loading and feature flags.
    - `internal/craas`: CRaaS service integration (modularized services).
    - `internal/api`: REST API handlers (split by domain: projects, registries, repositories, images).
- `frontend/`: Vue frontend source code.
    - `src/api`: Centralized Axios client.
    - `src/stores`: Pinia stores for state management (Registry, Config).
    - `src/views`: Vue components for pages.
    - `src/components`: Reusable UI components (ConfirmModal, ToastNotification).
    - `src/types`: Centralized TypeScript interfaces.

## Key Improvements

- **Refactored Architecture**: Backend handlers are now domain-specific, reducing code duplication via shared middleware
  for token management. Frontend uses a centralized API client and type definitions.
- **Concurrency**: Image tag verification uses `errgroup` for efficient parallel fetching.
- **User Safety**: "Destructive Action Guards" require explicit confirmation (typing the resource name) for critical
  deletions.
- **Optimized UX**: Single image deletion now uses the bulk cleanup endpoint to allow optional Garbage Collection
  control, and UI updates are optimistic.
