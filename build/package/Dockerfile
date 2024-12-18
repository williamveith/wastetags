# Use a Go base image
FROM golang:1.23.2-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project to the container
COPY . .

# Download dependencies
RUN go mod download

# Build the binary for Linux ARM64
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -ldflags="-s -w" -o wastetags ./cmd/wastetags/*.go

# Create a minimal final image
FROM alpine:3.18
COPY --from=builder /app/wastetags /wastetags