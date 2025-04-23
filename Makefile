# Makefile for building and managing the Paxi blockchain (Cosmos SDK + CometBFT)

APP_NAME = paxid
VERSION ?= v1.0.0
DOCKER_IMAGE = paxi-chain/node
BUILD_TAGS = "rocksdb"
CGO_ENABLED=1
CGO_CFLAGS="-I/usr/local/include" 
CGO_LDFLAGS="-L/usr/local/lib -lrocksdb"

.PHONY: all build install clean proto config run test

all: build

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -mod=readonly -tags $(BUILD_TAGS) -o build/$(APP_NAME) ./cmd/$(APP_NAME)

install:
	@echo "📦 Installing $(APP_NAME) to GOPATH/bin..."
	go install -mod=readonly -tags $(BUILD_TAGS) ./cmd/$(APP_NAME)

clean:
	@echo "🧹 Cleaning build files..."
	rm -rf build

run:
	@echo "🚀 Running Paxi node..."
	build/$(APP_NAME) start

version:
	@echo "📄 Version: $(VERSION)"

proto:
	@echo "📚 Generating protobuf files..."
	buf generate
	echo "✅ Protobufs generated."

test:
	@echo "🧪 Running unit tests..."
	go test ./... -v

docker:
	@echo "🐳 Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE):latest .
	echo "✅ Docker image built: $(DOCKER_IMAGE):latest"


PROTOC_GEN_GOGO := $(shell which protoc-gen-gogofaster)

proto:
	@echo "Generating gogofaster protos..."
	protoc \
	-I=proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/cosmos-sdk)/proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/gogoproto) \
	--gogofaster_out=plugins=grpc,paths=source_relative:./ \
	./proto/x/custommint/types/genesis.proto