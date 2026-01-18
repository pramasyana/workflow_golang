# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go modules and bash for scripts
RUN apk add --no-cache git bash

# Set GOTOOLCHAIN to auto to download correct toolchain if needed
ENV GOTOOLCHAIN=auto

# Copy go mod and sum files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Tidy dependencies and build (tidy after COPY to ensure go.mod is up to date)
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# Production stage
FROM alpine:3.19

# Install ca-certificates for HTTPS connections and wget for health check
RUN apk --no-cache add ca-certificates wget curl

# Create non-root user for security
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/main /app/

# Copy config file
COPY --from=builder /app/config/config.yaml /app/config/config.yaml

# Create logs directory
RUN mkdir -p /app/logs && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV APP_ENV=production
ENV TZ=UTC

# Expose port
EXPOSE 8080

# Health check using wget
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/main"]

