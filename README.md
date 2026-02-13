# Selectel CRaaS Web Interface

A full-stack web interface for managing Selectel Container Registry (CRaaS) repositories and images.

## Features

-   List Selectel projects.
-   View registries within a project.
-   List repositories in a registry.
-   Browse images and tags.
-   Delete registries, repositories, and images (with confirmation).

## Tech Stack

-   **Backend**: Go (Chi router, Selectel SDK).
-   **Frontend**: Vue 3, TypeScript, Pinia, Vite.

## Prerequisites

-   Go 1.20+
-   Node.js 18+
-   A Selectel account with API credentials.

## Configuration

1.  Copy `.env.example` to `backend/.env` (or set environment variables directly).
2.  Fill in your Selectel credentials:

```bash
WEB_PORT=8080
SELECTEL_USERNAME=your_username
SELECTEL_ACCOUNT_ID=your_account_id
SELECTEL_PASSWORD=your_password
```

## Running the Application

### Development Mode

You can run backend and frontend separately for development.

1.  **Backend**:
    ```bash
    cd backend
    go run cmd/server/main.go
    ```
    The server will start on port 8080 (or as configured).

2.  **Frontend**:
    ```bash
    cd frontend
    npm install
    npm run dev
    ```
    The frontend will start on http://localhost:5173 and proxy API requests to `http://localhost:8080`.

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

## Testing

-   **Backend**: `make test-backend`
-   **Frontend**: `make test-frontend`

## Project Structure

-   `backend/`: Go backend source code.
    -   `internal/auth`: Selectel Keystone authentication.
    -   `internal/craas`: CRaaS service integration.
    -   `internal/api`: REST API handlers.
-   `frontend/`: Vue frontend source code.
    -   `src/stores`: Pinia stores for state management.
    -   `src/views`: Vue components for pages.
