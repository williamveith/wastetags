# Build Stage
FROM golang:1.23.2-alpine3.20 AS builder

# Install build dependencies
RUN apk add --no-cache build-base sqlite-dev

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the working directory
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

# Build the Go binary
RUN go build -ldflags="-s -w" -o wastetags ./cmd/wastetags/main.go

# Final Stage
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache libc6-compat sqlite-libs

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/wastetags .

# Expose the application port
EXPOSE 8080

# Set the entrypoint to the compiled binary
ENTRYPOINT ["./wastetags"]
