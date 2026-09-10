# Build stage
FROM golang:1.24-alpine AS builder

ARG VERSION=dev

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Template extension manifest with build version.
# Talos reads metadata.version from this file, not the image tag.
RUN sed -i "s/__VERSION__/$(echo ${VERSION} | sed 's/^v//')/" manifest.yaml

# Build the binary using Makefile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build UPX_FLAGS= VERSION=${VERSION}

# Extension stage - Talos system extension format
FROM scratch

# Copy binary to Talos extension location (under /rootfs/)
COPY --from=builder /app/bin/kommodity-attestation-extension \
    /rootfs/usr/local/lib/containers/kommodity-attestation/kommodity-attestation-extension

# Copy service definition (under /rootfs/)
COPY kommodity-attestation.yaml \
    /rootfs/usr/local/etc/containers/kommodity-attestation.yaml

# Copy extension manifest (at root, not under /rootfs/)
# Uses the templated manifest from the builder stage.
COPY --from=builder /app/manifest.yaml /manifest.yaml
