# Selectel CRaaS Web Interface

A modern, full-stack web interface for managing Selectel Container Registry (CRaaS) repositories and images.

## Features

-   **Project & Registry Management**: Browse projects and registries within your Selectel account.
-   **Repository Insights**: View repositories and their details.
-   **Image Management**: List images with detailed metadata (tags, size, creation date).
-   **Bulk Cleanup**: Select multiple images to delete at once.
-   **Garbage Collection Control**: Option to trigger Garbage Collection (GC) immediately upon deletion.
-   **Clipboard Support**: Easily copy image digests.
-   **Responsive UI**: Built with Vue 3 and modern SCSS for a clean dark-mode experience.

## Tech Stack

-   **Backend**: Go 1.24+ (Chi router, Selectel SDK, Slog logging).
-   **Frontend**: Vue 3, TypeScript, Pinia (Setup Stores), Vite.

## Prerequisites

-   Go 1.24+
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
    # Install dependencies
    go mod download
    # Run the server
    go run cmd/server/main.go
    ```
    The server will start on port 8080 (or as configured). It supports graceful shutdown (Ctrl+C).

2.  **Frontend**:
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

-   `backend/`: Go backend source code.
    -   `cmd/server`: Entry point.
    -   `internal/auth`: Selectel Keystone authentication.
    -   `internal/craas`: CRaaS service integration (refactored into modular services).
    -   `internal/api`: REST API handlers.
-   `frontend/`: Vue frontend source code.
    -   `src/stores`: Pinia stores for state management (using Setup Store syntax).
    -   `src/views`: Vue components for pages.
    -   `src/components`: Reusable UI components (e.g., ToastNotification).

## Key Improvements

-   **Concurrency**: Image tag verification uses `errgroup` for efficient parallel fetching.
-   **Architecture**: Modular backend service structure (Images, Registries, Repositories separated).
-   **UX**: Real-time feedback with Toast notifications for success/error states.
-   **Robustness**: Graceful shutdown and improved error handling throughout the stack.
