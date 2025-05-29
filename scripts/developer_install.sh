#!/bin/bash
set -e

GOLANG_VERSION=1.24.2
ROCKSDB_VERSION=v9.2.1 
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
CHAIN_ID="my-testnet"
NODE_MONIKER="test-node-1"
BINARY_NAME="./paxid"
PAXI_PATH="$HOME/paxid"

### === Install dependencies ===
echo ""
sudo apt update
sudo apt-get update && sudo apt-get install -y \
    build-essential git cmake \
    libsnappy-dev zlib1g-dev libbz2-dev \
    liblz4-dev libzstd-dev wget curl pkg-config \
    ca-certificates 

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
make install
cd $PAXI_PATH
fi

### === Initialize node ===
if ! [ -f ./paxi/config/genesis.json ]; then
$BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
fi 

### === Common commands ===
echo ""
echo "This is your root directory:"
echo "cd $PAXI_PATH"
echo ""
echo "Initialize the node configuration (set your node name and chain ID):"
echo "$BINARY_NAME init your_node_name --chain-id my-testnet"
echo ""
echo "Create a new account (make sure to save your mnemonic safely):"
echo "$BINARY_NAME keys add your_account_name"
echo ""
echo "Allocate 1,000,000 PAXI tokens to your account in the genesis file:"
echo "$BINARY_NAME genesis add-genesis-account your_account_name 1000000000000upaxi"
echo ""
echo "Generate a genesis transaction by staking 900,000 PAXI:"
echo "$BINARY_NAME genesis gentx your_account_name 900000000000upaxi"
echo ""
echo "Aggregate all genesis transactions into the genesis file:"
echo "$BINARY_NAME genesis collect-gentxs"
echo ""
echo "Start the blockchain node:"
echo "$BINARY_NAME start"
