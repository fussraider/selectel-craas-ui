# Build stage
FROM golang:1.24-alpine AS builder

# Install git for fetch dependencies
RUN apk add --no-cache git upx

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically
# -ldflags="-w -s" reduces binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /backend ./cmd/server

# Compress binary
RUN upx --best --lzma /backend

# Final stage
FROM scratch

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /backend /backend

# Expose port
EXPOSE 8080

ENTRYPOINT ["/backend"]
