# Build stage
ARG GO_VERSION=1.25.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/api .

# Cloud Run uses PORT environment variable (default: 8080)
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Run the API
CMD ["./api"]
