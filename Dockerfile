# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/server

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install certificates (for HTTPS support)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose port (เปลี่ยนตามที่ Go app ใช้)
EXPOSE 8080

# Run the binary
CMD ["./main"]
