#!/bin/bash
set -e

echo "============================================================"
echo "🚨  PAXI Validator Node Installation Warning"
echo "============================================================"
echo ""
echo "🛑 CRITICAL WARNING:"
echo "❗ If more than 1/3 of validator nodes go offline, the entire blockchain will halt."
echo "❗ You must back up the entire paxi folder, especially your node's private keys (node_key.json, priv_validator_key.json, keyring)."
echo "   If your machine fails, this is the only way to restore your validator and reclaim your staking rewards."
echo ""
echo "⚠️ Please note:"
echo "   Once you stake and become a validator, the system will"
echo "   automatically monitor your online status."
echo ""
echo "❗ If you go offline without undelegating (e.g. shutdown or disconnect),"
echo "   the system will treat it as a slashing offense and"
echo "   automatically deduct a portion of your staked tokens."
echo ""
echo "✅ Correct way to go offline:"
echo "   Use the Undelegate command to leave the validator role"
echo "   before stopping or shutting down your node."
echo ""
echo "🚫 Shutting down your node directly may result in slashing penalties."
echo "   Please make sure you understand!"
echo ""
echo "============================================================"
read -p "Do you understand the above risks and wish to continue? (y/N): " confirm

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
PAXI_TAG="latest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
SEEDS="key@mainnet-seed.paxinet.io:26656"
PERSISTENT_PEERS="key@mainnet-node-1.paxinet.io:26656"
RPC_URL="http://mainnet-rpc.paxinet.io:26657"
SNAPSHOT_URL="http://mainnet-snapshot.paxinet.io:26657"
GENESIS_URL="$RPC_URL/genesis?"
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
fi 

curl -s $GENESIS_URL | jq -r .result.genesis > ./paxi/config/genesis.json

### === Set state sync ===
BLOCK_OFFSET=1000
LATEST_HEIGHT=$(curl -s "$SNAPSHOT_URL/block" | jq -r .result.block.header.height)
TRUST_HEIGHT=$((LATEST_HEIGHT - BLOCK_OFFSET))
TRUST_HASH=$(curl -s "$SNAPSHOT_URL/block?height=$TRUST_HEIGHT" | jq -r .result.block_id.hash)

if ! [[ "$LATEST_HEIGHT" =~ ^[0-9]+$ ]]; then
  echo "❌ Failed to retrieve trust height or hash. Please check the RPC URL."
  exit 1
fi

if [ "$LATEST_HEIGHT" -gt "$BLOCK_OFFSET" ]; then
  if [[ -z "$TRUST_HEIGHT" || -z "$TRUST_HASH" || "$TRUST_HASH" == "null" ]]; then
    echo "❌ Failed to retrieve trust height or hash. Please check the RPC URL."
    exit 1
  fi

  sed -i "/^\[statesync\]/,/^\[/{                               
    s|^enable *=.*|enable = true|g
    s|^rpc_servers *=.*|rpc_servers = \"$SNAPSHOT_URL,$SNAPSHOT_URL\"|g
    s|^trust_height *=.*|trust_height = $TRUST_HEIGHT|g
    s|^trust_hash *=.*|trust_hash = \"$TRUST_HASH\"|g
    s|^trust_period *=.*|trust_period = \"168h\"|g
  }" "$CONFIG"
fi

### === Configure seeds and peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CONFIG
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === Disable unnecessary ports for security ===
sed -i '/^\[rpc\]/,/^\[/s|^\s*laddr\s*=.*|laddr = "tcp://0.0.0.0:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === Create wallet (if not exists) ===
if ! $BINARY_NAME keys show $KEY_NAME &>/dev/null; then
  echo ""
  echo "After creating the wallet, please write down the mnemonic phrase by hand"
  echo "to recover your wallet in case of loss."
  $BINARY_NAME keys add $KEY_NAME 
fi

### === Display address ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a )
echo ""
echo "Your wallet address is: $ADDR"
echo "Please send tokens to this address and then run the following command to become a validator:"

### === Display create-validator command ===
COUNTRY_CODE=$(curl -s http://ip-api.com/json | jq -r .countryCode)
VAL_PUBKEY=$($BINARY_NAME tendermint show-validator)
echo "You can modify validator.json at: $PAXI_DATA_PATH/validator.json"
echo "Please add your country code at the end of the 'details' parameter, e.g., [US]. This helps us collect node location data and display it on the official website."
echo "Generating validator.json..."
cat <<EOF > $PAXI_DATA_PATH/validator.json
{
  "pubkey": $VAL_PUBKEY,
  "amount": "1000000000$DENOM",
  "moniker": "$NODE_MONIKER",
  "identity": "",
  "website": "",
  "security": "",
  "details": "PAXI validator node [$COUNTRY_CODE]",
  "commission-rate": "0.1",
  "commission-max-rate": "0.2",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF
echo ""
echo "Command to become a validator (copy and paste to run):"
echo "cd $PAXI_PATH && $BINARY_NAME tx staking create-validator $PAXI_DATA_PATH/validator.json \\"
echo "  --from $KEY_NAME \\"
echo "  --fees 10000$DENOM"

### === Common commands ===
echo ""
echo "Before starting the node, remember to set 'your public IP:26656' to the 'external_address' parameter of paxi/config/config.toml, otherwise others will not be able to connect to your node"
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
