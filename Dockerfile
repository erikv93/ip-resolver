# Use a small official Go image
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .


# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Create minimal image from scratch
FROM scratch

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 80

# Run the app
ENTRYPOINT ["./server"]