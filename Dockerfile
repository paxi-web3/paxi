# Stage 1: Use Debian for building
FROM debian:bullseye AS builder

# Install required build tools and dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    cmake \
    libsnappy-dev \
    zlib1g-dev \
    libbz2-dev \
    liblz4-dev \
    libzstd-dev \
    wget \
    curl \
    pkg-config \
    ca-certificates \
    libgflags-dev

# Install Go
ENV GOLANG_VERSION=1.24.2
RUN curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go

ENV PATH="/usr/local/go/bin:${PATH}"

# Build paxid
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Build the binary with required build tags (cosmwasm and pebbledb)
RUN go build -mod=readonly -tags "pebbledb cosmwasm" -o paxid ./cmd/paxid

# Copy wasmvm from the cache
RUN mkdir -p /root/.wasmvm/lib && \
    cp /root/go/pkg/mod/github.com/!cosm!wasm/wasmvm/*/internal/api/libwasmvm.x86_64.so /root/.wasmvm/lib/

# Stage 2: runtime image
FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libstdc++6 \
        libsnappy-dev \
        zlib1g-dev \
        libbz2-dev \
        libgflags-dev \
        liblz4-dev \
        libzstd-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/paxid .

# Copy the wasm dynamic lib from the build stage
COPY --from=builder /root/.wasmvm/lib/libwasmvm* /usr/local/lib/
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/wasmvm.conf && ldconfig

# Expose typical Cosmos SDK ports:
# - 26656: p2p
# - 26657: RPC
# - 1317: REST API (legacy)
# - 9090: gRPC
EXPOSE 26656 26657 1317 9090

# Default command to run your node
CMD ["./paxid", "start"]
