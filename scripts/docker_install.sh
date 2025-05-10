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
  echo "❌ Installation cancelled. Please review the instructions and run this script again."
  exit 1
fi

echo "📝 Please enter your node moniker (name):"
read -p "Node name: " NODE_MONIKER

if [[ -z "$NODE_MONIKER" ]]; then
  echo "❌ Node name cannot be empty. Please rerun the script."
  exit 1
fi

echo "✅ Node name set to: $NODE_MONIKER"

echo "📝 Please enter a name for your wallet (key name):"
read -p "Wallet name (key name): " KEY_NAME
if [[ -z "$KEY_NAME" ]]; then
  echo "❌ Wallet name cannot be empty. Please rerun the script."
  exit 1
fi
echo "✅ Wallet name set to: $KEY_NAME"

PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
SEEDS="mainnet-seed.paxinet.io:26656"
PERSISTENT_PEERS="key@mainnet-node-1.paxinet.io:26656"
RPC_URL="http://mainnet-rpc.paxinet.io:26657"
GENESIS_URL="$RPC_URL/genesis?"
CONFIG="./paxi/config/config.toml"
APP_CONFIG="./paxi/config/app.toml"
PAXI_PATH="$HOME/paxid"
PAXI_DATA_PATH="$HOME/paxid/paxi"
DENOM="upaxi"
DOCKER_IMAGE="paxi-node"
DOCKER_PAXI_DATA_PATH="/root/paxi"

### === Install dependencies ===
echo ""
sudo apt-get update
sudo apt-get install -y \
    ca-certificates curl gnupg lsb-release git make

### === Install Docker ===
if ! command -v docker &> /dev/null; then
  echo "Installing Docker..."
  sudo mkdir -p /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
    sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

  echo \
    "deb [arch=$(dpkg --print-architecture) \
    signed-by=/etc/apt/keyrings/docker.gpg] \
    https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

  sudo apt-get update
  sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

  # Enable Docker for non-root users
  sudo systemctl enable docker
  sudo systemctl start docker
  sudo usermod -aG docker $USER

  echo "⚠️ You may need to log out and back in (or run 'newgrp docker') to apply Docker permissions."
else
  echo "✅ Docker is already installed."
fi

### === Install Paxi ===
if [ ! -d "paxi" ]; then
  git clone $PAXI_REPO
  cd paxi
  git checkout $PAXI_TAG
  make docker
else
  cd paxi
  make docker
fi

if [ ! -d "$HOME/paxid" ]; then
  mkdir "$HOME/paxid" 
fi
cd $HOME/paxid 

### === Initialize node ===
if ! [ -f ./paxi/config/genesis.json ]; then
docker run --rm \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
sudo chown -R $(whoami) $HOME/paxid
fi 

curl -s $GENESIS_URL | jq -r .result.genesis > ./paxi/config/genesis.json

### === Set state sync ===
BLOCK_OFFSET=100
LATEST_HEIGHT=$(curl -s "$RPC_URL/block" | jq -r .result.block.header.height)
TRUST_HEIGHT=$((LATEST_HEIGHT - BLOCK_OFFSET))
TRUST_HASH=$(curl -s "$RPC_URL/block?height=$TRUST_HEIGHT" | jq -r .result.block_id.hash)

if [[ -z "$TRUST_HEIGHT" || -z "$TRUST_HASH" || "$TRUST_HASH" == "null" ]]; then
  echo "❌ Failed to retrieve trust height or hash. Please check the RPC URL."
  exit 1
fi

sed -i "/^\[statesync\]/,/^\[/{                               
  s|^enable *=.*|enable = true|g
  s|^rpc_servers *=.*|rpc_servers = \"$RPC_URL,$RPC_URL\"|g
  s|^trust_height *=.*|trust_height = $TRUST_HEIGHT|g
  s|^trust_hash *=.*|trust_hash = \"$TRUST_HASH\"|g
  s|^trust_period *=.*|trust_period = \"168h\"|g
}" "$CONFIG"

### === Configure seeds and peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CONFIG
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === Disable unused ports for security ===
sed -i '/^\[rpc\]/,/^\[/s|^\s*laddr\s*=.*|laddr = "tcp://0.0.0.0:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === Create wallet (if not exists) ===
if ! [ -f ./paxi/keyring-file/$KEY_NAME.info ]; then
  echo ""
  echo "After the wallet is created, please write down the mnemonic phrase"
  echo "manually to ensure you can recover it if lost."

  docker run --rm -it \
    -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
    $DOCKER_IMAGE \
    $BINARY_NAME keys add $KEY_NAME --keyring-backend file
fi
sudo chown -R $(whoami) $HOME/paxid

### === Show wallet address ===
echo ""
echo "Your wallet address:"
docker run --rm -it \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME keys show $KEY_NAME -a --keyring-backend file
echo "Please fund this address and run the following command to become a validator:"

### === Generate create-validator command ===
COUNTRY_CODE=$(curl -s http://ip-api.com/json | jq -r .countryCode)
VAL_PUBKEY=$(docker run --rm \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME tendermint show-validator)
echo "You may edit validator.json at: $PAXI_DATA_PATH/validator.json"
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

### === Common commands ===
echo ""
echo "Start the node (this node runs in the background and auto-starts on system reboot unless stopped manually):"
echo "docker run -d --name paxi-node-1 --restart unless-stopped -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \\"
echo "-p 26656:26656 -p 26657:26657 -p 1317:1317 -p 9090:9090 \\"
echo "paxi-node \\"
echo "$BINARY_NAME start"
echo ""
echo "For convenience, please enter the container before running the following commands:"
echo "docker run --rm -it --network host -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH $DOCKER_IMAGE bash"
echo ""
echo "Validator creation command (copy & paste to execute):"
echo "$BINARY_NAME tx staking create-validator $DOCKER_PAXI_DATA_PATH/validator.json \\"
echo "  --from $KEY_NAME --keyring-backend file \\"
echo "  --fees 10000$DENOM"
echo ""
echo "To view your wallet address:"
echo "$BINARY_NAME keys show $KEY_NAME --keyring-backend file"
echo ""
echo "Check wallet balance:"
echo "$BINARY_NAME query bank balances <your address or key name> --keyring-backend file"
echo ""
echo "Check your staking rewards:"
echo "$BINARY_NAME query distribution rewards <your address or key name> --keyring-backend file"
echo ""
echo "Check your validator public key:"
echo "$BINARY_NAME tendermint show-validator"
echo ""
echo "Check your validator's outstanding rewards:"
echo "$BINARY_NAME query distribution validator-outstanding-rewards <your validator address>"
