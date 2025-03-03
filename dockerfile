# Use official Golang image with correct version
FROM golang:1.23 AS builder  

# Set working directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./  
RUN go mod download  

# Copy the source code
COPY . .

# Build the application
RUN go build -o healthsync ./cmd/main.go  

# Use a lightweight image for the final container
FROM alpine:3.18  

WORKDIR /root/

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates libc6-compat  

# Copy the compiled binary from the builder stage
COPY --from=builder /app/healthsync .

# Ensure the binary has execution permissions
RUN chmod +x ./healthsync  

# Expose application port
EXPOSE 8080  

# Set entrypoint
ENTRYPOINT ["/root/healthsync"]
