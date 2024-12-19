# Use a Go base image
FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine3.20 AS builder

RUN apk add --no-cache \
    build-base \
    musl-dev \
    linux-headers

# Set the working directory inside the container
WORKDIR /app

ARG TARGETOS TARGETARCH

# Build the binary for Linux ARM64
RUN --mount=type=bind,source=.,target=/app,readonly \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=1 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /tmp/wastetags ./cmd/wastetags/*.go

# Create a minimal final image
FROM busybox:1.37.0 AS export
COPY --from=builder /tmp/wastetags /wastetags
RUN chmod +x /wastetags