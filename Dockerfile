# Multi-stage build for smaller image size
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application from cmd/api directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port (adjust to your service port)
EXPOSE 8080

# Run the application
CMD ["./main"]