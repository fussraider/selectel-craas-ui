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
- **Linting**: Oxlint and ESLint for frontend code quality.

## Prerequisites

- Go 1.24+
- Node.js 18+
- A Selectel account with API credentials.

## Configuration

The application is configured via environment variables. You can set these directly or use a `.env` file in the
`backend/` directory.

### Selectel Authorization

To connect to the Selectel Container Registry, you need credentials for a user with access to your project. It is **strongly recommended** to create a dedicated Service User for this purpose.

1.  Log in to the **Selectel Control Panel**.
2.  Navigate to **Identity & Access Management** > **User Management**.
3.  Create a new **Service User** (or use an existing one).
    -   Note the **User Name** (`SELECTEL_USERNAME`) and **Password** (`SELECTEL_PASSWORD`).
    -   Note the **Account ID** (`SELECTEL_ACCOUNT_ID`) from the top-right corner of the panel (numeric ID).
4.  Ensure the user has access to the project containing your registries.
    -   Note the **Project Name** (`SELECTEL_PROJECT_NAME`).

### Core Configuration

| Variable                | Description                      | Default    |
|:------------------------|:---------------------------------|:-----------|
| `WEB_PORT`              | Port for the backend server      | `8080`     |
| `SELECTEL_USERNAME`     | Selectel Username                | (Required) |
| `SELECTEL_ACCOUNT_ID`   | Selectel Account ID              | (Required) |
| `SELECTEL_PASSWORD`     | Selectel Password                | (Required) |
| `SELECTEL_PROJECT_NAME` | Selectel Project Name (Required) | (Required) |
| `SELECTEL_AUTH_URL`     | Selectel Auth URL                | `https://cloud.api.selcloud.ru/identity/v3/auth/tokens` |
| `SELECTEL_PROJ_URL`     | Selectel Projects URL            | `https://cloud.api.selcloud.ru/identity/v3/auth/projects` |
| `SELECTEL_CRAAS_URL`    | Selectel CRaaS API Endpoint      | `https://cr.selcloud.ru/api/v1` |
| `CORS_ALLOWED_ORIGIN`   | Allowed Origin for CORS requests | `*`        |

### Web Interface Security

You can protect the web interface with Basic Authentication to restrict access.

| Variable        | Description                              | Default               |
|:----------------|:-----------------------------------------|:----------------------|
| `AUTH_ENABLED`  | Enable Basic Auth for the web interface  | `false`               |
| `AUTH_LOGIN`    | Username for web interface login         | (Required if enabled) |
| `AUTH_PASSWORD` | Password for web interface login         | (Required if enabled) |
| `JWT_SECRET`    | Secret key for signing JWT tokens        | (Required if enabled) |

### Feature Flags & Safety

Control which operations are permitted and protect critical resources.

| Variable                   | Description                                         | Default |
|:---------------------------|:----------------------------------------------------|:--------|
| `ENABLE_DELETE_REGISTRY`   | Allow deletion of entire registries                 | `false` |
| `ENABLE_DELETE_REPOSITORY`    | Allow deletion of repositories                                   | `false` |
| `ENABLE_DELETE_IMAGE`         | Allow deletion of images (single or bulk)                        | `false` |
| `ENABLE_MISSING_TAGS_CHECK`   | Check and fetch missing tags for images (performance heavy)      | `false` |
| `PROTECTED_TAGS`              | Comma-separated list of tags that cannot be deleted              | (empty) |

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

To run lint checks:

```bash
make lint-frontend
```

To run the built artifacts:

```bash
make run-backend
# In a separate terminal
make run-frontend
```

## Docker Deployment

This repository includes a `docker-compose.example.yml` file to demonstrate how to deploy the application using Docker
Compose.

### Frontend Configuration

The frontend image supports runtime configuration for the backend API URL. This allows you to build the image once and
deploy it to different environments without rebuilding.

Set the `API_BASE_URL` environment variable for the frontend container to point to your backend API.

**Note:** The `API_BASE_URL` must be accessible from the user's browser (e.g., `https://api.example.com` or `http://localhost:8080` for local development). Do not use internal Docker network aliases (like `http://backend:8080`) as the frontend runs in the client's browser.

### Example Usage

1. Copy the example file:
   ```bash
   cp docker-compose.example.yml docker-compose.yml
   ```

2. Edit `docker-compose.yml` to set your Selectel credentials and preferred configuration.

3. Run the stack:
   ```bash
   docker-compose up -d
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
