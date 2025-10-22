#!/bin/bash
set -e

NODE_MONIKER="my-paxi-node"
GOLANG_VERSION=1.24.2
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
PERSISTENT_PEERS="d411fc096e0d946bbd2bdea34f0f9da928c1a714@139.99.68.32:26656,509b20ca82d34d0aae1751a681ee386659fb71da@66.70.181.55:26656,a325cced9b360c0e5fcbf756e0b1ca139b8f2eef@51.75.54.185:26656,9e64baa45042ae29d999f2677084c9579972722c@139.99.69.74:26656"
RPC_URL="http://rpc.paxi.info"
SNAPSHOT_URL="http://rpc.paxi.info"
SNAPSHOT_DOWNLOAD_HOST="http://snapshot.paxi.info"
GENESIS_URL="$RPC_URL/genesis?"
CONFIG="./paxi/config/config.toml"
APP_CONFIG="./paxi/config/app.toml"
PAXI_PATH="$HOME/paxid"
PAXI_DATA_PATH="$HOME/paxid/paxi"
DENOM="upaxi"

### === Install dependencies ===
echo ""
sudo apt-get update && sudo apt-get install -y \
    build-essential git cmake \
    libsnappy-dev zlib1g-dev libbz2-dev \
    liblz4-dev libzstd-dev wget curl pkg-config \
    ca-certificates unzip jq

### === Install Go ===
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    curl -LO https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo ln -s /usr/local/go/bin/go /usr/bin/go
fi

### === Compile Paxi ===
if ! [ -d ./paxi ]; then
echo "Installing Paxi..."
git clone $PAXI_REPO
cd paxi
git checkout $PAXI_TAG
make install
cd $PAXI_PATH
else
cd paxi
git checkout $PAXI_TAG
make install
cd $PAXI_PATH
fi

### === Copy libwasmvm.x86_64.so to /usr/loca/lib ===
WASM_SO_FILE="$HOME/go/pkg/mod/github.com/!cosm!wasm/wasmvm/v2@v2.2.1/internal/api/libwasmvm.x86_64.so"
if [ -f "$WASM_SO_FILE" ]; then
  sudo cp $WASM_SO_FILE /usr/local/lib/
  sudo ldconfig
fi

### === Initialize node ===
if ! [ -f ./paxi/config/genesis.json ]; then
$BINARY_NAME init "$NODE_MONIKER" --chain-id $CHAIN_ID
fi 

curl -s $GENESIS_URL | jq -r .result.genesis > ./paxi/config/genesis.json

### === Set state sync ===
BLOCK_OFFSET=100
LATEST_HEIGHT=$(curl -s "$RPC_URL/block" | jq -r .result.block.header.height)
TRUST_HEIGHT=$(( ( (LATEST_HEIGHT - BLOCK_OFFSET) / BLOCK_OFFSET ) * BLOCK_OFFSET ))
TRUST_HASH=$(curl -s "$RPC_URL/block?height=$TRUST_HEIGHT" | jq -r .result.block_id.hash)

if ! [[ "$LATEST_HEIGHT" =~ ^[0-9]+$ ]]; then
  echo "❌ Failed to retrieve trust height or hash. Please check the RPC URL."
  exit 1
fi

### === Download wasm snapshot ===
WASM_SNAPSHOT_URL=$(curl -s "$SNAPSHOT_DOWNLOAD_HOST/utils/latest_wasm_snapshot" | jq -r .url)
curl -f -o wasm_snapshot.zip "$WASM_SNAPSHOT_URL"
if [ $? -ne 0 ]; then
  echo "❌ Failed to download wasm snapshot. Please download it and unzip it to $PAXI_DATA_PATH/wasm/wasm/state/wasm."
else
  mkdir -p "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  unzip -o wasm_snapshot.zip -d "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  rm wasm_snapshot.zip
  echo "✅ Wasm snapshot downloaded and extracted to $PAXI_DATA_PATH/wasm/wasm/state/wasm."
fi

### === Detect platform ===
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED="sed -i ''"
else
    SED="sed -i"
fi

if [ "$LATEST_HEIGHT" -gt "$BLOCK_OFFSET" ]; then
  if [[ -z "$TRUST_HEIGHT" || -z "$TRUST_HASH" || "$TRUST_HASH" == "null" ]]; then
    echo "❌ Failed to retrieve trust height or hash. Please check the RPC URL."
    exit 1
  fi

  $SED "/^\[statesync\]/,/^\[/{ \
    s|^enable *=.*|enable = true|g; \
    s|^rpc_servers *=.*|rpc_servers = \"$SNAPSHOT_URL,$SNAPSHOT_URL\"|g; \
    s|^trust_height *=.*|trust_height = $TRUST_HEIGHT|g; \
    s|^trust_hash *=.*|trust_hash = \"$TRUST_HASH\"|g; \
    s|^trust_period *=.*|trust_period = \"168h\"|g; \
    s|^discovery_time *=.*|discovery_time = \"30s\"|g; \
  }" "$CONFIG"
fi

### === Configure seeds and peers ===
$SED "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === Common commands ===
echo ""
echo "This is your root directory:"
echo "cd $PAXI_PATH"
echo ""
echo "Start the blockchain node:"
echo "$BINARY_NAME start"
echo ""
echo "Send token:"
echo "$BINARY_NAME tx bank send <your address or key name> <receiver's address> <amount>$DENOM --fees 10000$DENOM --gas auto"