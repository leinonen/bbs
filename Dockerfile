# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gobbs .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gobbs .
COPY --from=builder /app/config.example.json ./config.json

# Create directory for database
RUN mkdir -p /data

# Initialize database
RUN ./gobbs -init

# Expose SSH port
EXPOSE 2222

# Volume for persistent data
VOLUME ["/data"]

# Run the BBS
CMD ["./gobbs"]