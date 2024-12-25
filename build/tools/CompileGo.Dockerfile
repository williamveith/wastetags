# Use a Go base image
FROM --platform=$BUILDPLATFORM golang@sha256:9dd2625a1ff2859b8d8b01d8f7822c0f528942fe56cfe7a1e7c38d3b8d72d679 AS builder

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
FROM busybox@sha256:2919d0172f7524b2d8df9e50066a682669e6d170ac0f6a49676d54358fe970b5 AS export
COPY --from=builder /tmp/wastetags /wastetags