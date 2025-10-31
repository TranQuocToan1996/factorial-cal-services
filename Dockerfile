# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached unless go.mod/go.sum change)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /bin/api ./cmd/api

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=build /bin/api /bin/api
COPY --from=build /app/docs /app/docs
RUN apk add --no-cache ca-certificates
# RUN apk add --no-cache curl
ENTRYPOINT ["/bin/api"]
