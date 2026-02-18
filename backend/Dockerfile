# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

# Install necessary tools
RUN apk add --no-cache git upx

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
# -ldflags="-s -w" strips debug information for smaller size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Compress binary with UPX
RUN upx --best --lzma /app/server

# Stage 2: Create minimal image
FROM scratch

# Copy SSL certificates for HTTPS requests (needed for Selectel API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compressed binary
COPY --from=builder /app/server /server

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
