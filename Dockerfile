# Multi-stage Dockerfile for API and Worker

# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# API Image
FROM alpine:3.18 AS api
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/api /app/api
COPY --from=builder /app/docs /app/docs
EXPOSE 8080
ENTRYPOINT ["/app/api"]

# Worker Image
FROM alpine:3.18 AS worker
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/worker /app/worker
ENTRYPOINT ["/app/worker"]
