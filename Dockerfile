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

# Build and install RocksDB (v9.2.1) from source
WORKDIR /deps
RUN git clone https://github.com/facebook/rocksdb.git && \
    cd rocksdb && \
    git checkout v9.2.1 && \
    mkdir build && cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release \
          -DWITH_TESTS=OFF \
          -DWITH_TOOLS=OFF \
          -DWITH_GFLAGS=ON \
          -DPORTABLE=ON \
          -DCMAKE_INSTALL_PREFIX=/usr/local .. && \
    make -j$(nproc) && \
    make install

# Update dynamic linker with RocksDB path
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/rocksdb.conf && \
    ldconfig

# Build paxid
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Build the binary with required build tags (cosmwasm and rocksdb)
RUN go build -mod=readonly -tags "rocksdb cosmwasm" -o paxid ./cmd/paxid

# Copy wasmvm from the cache
RUN mkdir -p /root/.wasmvm/lib && \
    cp /root/go/pkg/mod/github.com/!cosm!wasm/wasmvm/*/internal/api/libwasmvm.x86_64.so /root/.wasmvm/lib/


RUN find / -name 'libwasmvm*' 2>/dev/null

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

# Copy the rocksdb dynamic lib from the build stage
COPY --from=builder /usr/local/lib/librocksdb.so* /usr/local/lib/
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/rocksdb.conf && \
    ldconfig

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
