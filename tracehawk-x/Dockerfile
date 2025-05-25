# TraceHawk X Dockerfile
# Multi-stage build for minimal final image

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    nmap \
    curl \
    bind-tools \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 tracehawk && \
    adduser -D -s /bin/sh -u 1000 -G tracehawk tracehawk

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/build/tracehawkx /usr/local/bin/tracehawkx

# Create directories for output
RUN mkdir -p /app/results /app/reports && \
    chown -R tracehawk:tracehawk /app

# Switch to non-root user
USER tracehawk

# Set default command
ENTRYPOINT ["tracehawkx"]
CMD ["--help"]

# Labels
LABEL maintainer="TraceHawk Team"
LABEL description="Autonomous Terminal-Native Recon & Exploitation Agent"
LABEL version="1.0.0" 