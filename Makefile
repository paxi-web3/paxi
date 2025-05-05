# Makefile for building and managing the Paxi blockchain (Cosmos SDK + CometBFT)

APP_NAME = paxid
VERSION ?= v1.0.2
DOCKER_IMAGE = paxi-node
BUILD_TAGS = "rocksdb cosmwasm"
CGO_ENABLED=1
CGO_CFLAGS="-I/usr/local/include" 
CGO_LDFLAGS="-L/usr/local/lib -lrocksdb"

.PHONY: all build install clean proto config run test

all: build

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -mod=readonly -tags $(BUILD_TAGS) -o build/$(APP_NAME) ./cmd/$(APP_NAME)

install:
	@echo "📦 Installing $(APP_NAME) to $$HOME/paxid..."
	go build -mod=readonly -tags $(BUILD_TAGS) -o $$HOME/paxid/$(APP_NAME) ./cmd/$(APP_NAME)

clean:
	@echo "🧹 Cleaning build files..."
	rm -rf build

run:
	@echo "🚀 Running Paxi node..."
	build/$(APP_NAME) start

version:
	@echo "📄 Version: $(VERSION)"

test:
	@echo "🧪 Running unit tests..."
	go test ./... -v

docker:
	@echo docker run -d --name paxi-node"🐳 Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE):latest . --progress=plain
	echo "✅ Docker image built: $(DOCKER_IMAGE):latest"


PROTOC_GEN_GOGO := $(shell which protoc-gen-gogofaster)

proto:
	@echo "Generating gogofaster protos..."
	protoc \
	-I=proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/cosmos-sdk)/proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/cosmos-proto)/proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/gogoproto) \
	--gogofaster_out=plugins=grpc,paths=source_relative:./ \
	./proto/x/custommint/types/tx.proto \
	./proto/x/custommint/types/query.proto \

	protoc \
	-I=proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/cosmos-sdk)/proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/cosmos-proto)/proto \
	-I=$(shell go list -m -f '{{ .Dir }}' github.com/cosmos/gogoproto) \
	--gogofaster_out=plugins=grpc,paths=source_relative:./ \
	./proto/x/paxi/types/query.proto \
	./proto/x/paxi/types/tx.proto