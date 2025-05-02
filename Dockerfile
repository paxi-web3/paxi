# Stage 1: Use Debian for building
FROM debian:bullseye as builder

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

# Stage 2: runtime image
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libstdc++6 && apt-get clean

# Create working directory
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/paxid .

# Expose typical Cosmos SDK ports:
# - 26656: p2p
# - 26657: RPC
# - 1317: REST API (legacy)
# - 9090: gRPC
EXPOSE 26656 26657 1317 9090

# Default command to run your node
CMD ["./paxid", "start"]
