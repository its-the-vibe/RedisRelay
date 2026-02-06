# Build stage
FROM golang:1.26rc3-alpine AS builder

WORKDIR /build

# Copy all source files
COPY . .

# Download dependencies and build
# Note: If you encounter TLS errors during build, ensure your Docker environment has proper CA certificates
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o redisrelay .

# Runtime stage
FROM scratch

# Copy CA certificates for HTTPS connections if needed
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder
COPY --from=builder /build/redisrelay /redisrelay

# Set the entrypoint
ENTRYPOINT ["/redisrelay"]
