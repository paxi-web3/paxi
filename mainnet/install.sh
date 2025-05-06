#!/bin/bash
set -e

echo "============================================================"
echo "🚨  PAXI Validator Node Installation Warning"
echo "============================================================"
echo ""
echo "⚠️ Please note:"
echo "   Once you stake and become a validator, the system will"
echo "   automatically monitor your online status."
echo ""
echo "❗ If you go offline without reason (disconnect or shut down),"
echo "   the system will consider it as malicious behavior and"
echo "   automatically slash a portion of your staked tokens."
echo ""
echo "✅ Proper way to go offline:"
echo "   Please use the Undelegate command to exit the validator role"
echo "   before shutting down your node."
echo ""
echo "🚫 Shutting down the node directly without undelegating risks slashing penalties."
echo "   Please make sure you understand!"
echo ""
echo "============================================================"
read -p "Do you understand the risks above and want to continue? (y/N): " confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
  echo "❌ Installation cancelled. Please read the instructions again before proceeding."
  exit 1
fi

echo "📝 Please enter the name for your node (moniker):"
read -p "Node name: " NODE_MONIKER

if [[ -z "$NODE_MONIKER" ]]; then
  echo "❌ Node name cannot be empty. Please rerun the script."
  exit 1
fi

echo "✅ Node name set to: $NODE_MONIKER"

echo "📝 Please enter the name for your wallet (key name):"
read -p "Wallet name (key name): " KEY_NAME
if [[ -z "$KEY_NAME" ]]; then
  echo "❌ Wallet name cannot be empty. Please rerun the script."
  exit 1
fi
echo "✅ Wallet name set to: $KEY_NAME"

GOLANG_VERSION=1.24.2
ROCKSDB_VERSION=v9.2.1 
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="lastest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
GENESIS_URL="https://raw.githubusercontent.com/paxi-web3/mainnet/genesis.json"
SEEDS="mainner-seed-1.paxi.io:26656"
PERSISTENT_PEERS="key@mainnet-node-1.paxi.io:26656"
CONFIG="./paxi/config/config.toml"
APP_CONFIG="./paxi/config/app.toml"
PAXI_PATH="$HOME/paxid"
PAXI_DATA_PATH="$HOME/paxid/paxi"
DENOM="upaxi"

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

### === Install RocksDB ===
if ! [ -f /usr/local/lib/librocksdb.so ]; then
    echo "Installing RocksDB..."

    if ! [ -d ./rocksdb ]; then
      git clone https://github.com/facebook/rocksdb.git
    fi
    
    cd rocksdb
    git checkout $ROCKSDB_VERSION
    make -j$(nproc) shared_lib
    sudo make install-shared INSTALL_PATH=/usr/local
    sudo echo "/usr/local/lib" | sudo tee /etc/ld.so.conf.d/rocksdb.conf
    sudo ldconfig && cd ..
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
curl -L $GENESIS_URL > ./paxi/config/genesis.json
fi 

### === Configure seeds and peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CONFIG
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === Disable unnecessary ports for security ===
sed -i 's|^laddr *=.*|laddr = "tcp://127.0.0.1:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === Create wallet (if not exists) ===
if ! $BINARY_NAME keys show $KEY_NAME --keyring-backend os &>/dev/null; then
  echo ""
  echo "After creating the wallet, please write down the mnemonic phrase by hand"
  echo "to recover your wallet in case of loss."
  $BINARY_NAME keys add $KEY_NAME --keyring-backend os
fi

### === Display address ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a --keyring-backend os)
echo ""
echo "Your wallet address is: $ADDR"
echo "Please send tokens to this address and then run the following command to become a validator:"

### === Display create-validator command ===
VAL_PUBKEY=$($BINARY_NAME tendermint show-validator)
echo "You can modify validator.json at: $PAXI_DATA_PATH/validator.json"
echo "Generating validator.json..."
cat <<EOF > $PAXI_DATA_PATH/validator.json
{
  "pubkey": $VAL_PUBKEY,
  "amount": "1000000000$DENOM",
  "moniker": "$NODE_MONIKER",
  "identity": "",
  "website": "",
  "security": "",
  "details": "PAXI validator initialized by install.sh",
  "commission-rate": "0.1",
  "commission-max-rate": "0.2",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF
echo ""
echo "Command to become a validator (copy and paste to run):"
echo "cd $PAXI_PATH && $BINARY_NAME tx staking create-validator ./paxi/validator.json \\"
echo "  --from $KEY_NAME --keyring-backend os \\"
echo "  --fees 10000$DENOM"

### === Common commands ===
echo ""
echo "Start the node:"
echo "$BINARY_NAME start"
echo ""
echo "Check wallet balance:"
echo "$BINARY_NAME query bank balances <your address or wallet name>"
echo ""
echo "Check your staking rewards:"
echo "$BINARY_NAME query distribution rewards <your address or wallet name>"
echo ""
echo "Check your validator public key:"
echo "$BINARY_NAME tendermint show-validator"
echo ""
echo "Check your validator rewards:"
echo "$BINARY_NAME query distribution validator-outstanding-rewards <your validator address>"
