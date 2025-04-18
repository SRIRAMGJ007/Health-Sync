# Use Go 1.23 image
FROM golang:1.23

# Set the working directory
WORKDIR /Health-Sync

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

WORKDIR /Health-Sync/cmd

# Build the Go binary
RUN go build -o server .

# Expose ports
EXPOSE 8080 8443

# Run the server
CMD ["./server"]
