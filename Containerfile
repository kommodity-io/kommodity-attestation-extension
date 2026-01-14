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

# Build the binary using Makefile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build UPX_FLAGS= VERSION=${VERSION}

# Extension stage - Talos system extension format
FROM scratch

# Copy binary to Talos extension location
COPY --from=builder /app/bin/kommodity-attestation-extension \
    /usr/local/lib/containers/kommodity-attestation/kommodity-attestation-extension

# Copy service definition
COPY kommodity-attestation.yaml \
    /usr/local/etc/containers/kommodity-attestation.yaml

# Copy extension manifest
COPY manifest.yaml /manifest.yaml
