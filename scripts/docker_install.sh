#!/bin/bash
set -e

### === Required parameters ===
REQUIRED_CPU=8
REQUIRED_RAM_GB=30
REQUIRED_DISK_GB=500

### === Get CPU cores ===
CPU_CORES=$(nproc)

### === Get total memory (in GB) ===
TOTAL_MEM=$(free -g | awk '/^Mem:/{print $2}')

### === Get available disk space (in GB), using root directory / ===
DISK_SPACE=$(df -BG / | awk 'NR==2 {gsub("G","",$4); print $4}')

### === Check CPU ===
if [ "$CPU_CORES" -lt "$REQUIRED_CPU" ]; then
  echo "❌ Insufficient CPU cores: Required at least ${REQUIRED_CPU} cores, but found ${CPU_CORES} cores"
  exit 1
fi

### === Check RAM ===
if [ "$TOTAL_MEM" -lt "$REQUIRED_RAM_GB" ]; then
  echo "❌ Insufficient memory: Required at least ${REQUIRED_RAM_GB} GB, but found ${TOTAL_MEM} GB"
  exit 1
fi

### === Check disk space ===
if [ "$DISK_SPACE" -lt "$REQUIRED_DISK_GB" ]; then
  echo "❌ Insufficient disk space: Required at least ${REQUIRED_DISK_GB} GB, but found ${DISK_SPACE} GB"
  exit 1
fi

echo "============================================================"
echo "🚨  PAXI Validator Node Installation Warning"
echo "============================================================"
echo ""
echo "🛑 CRITICAL WARNING:"
echo "❗ If more than 1/3 of validator nodes go offline, the entire blockchain will halt."
echo "❗ You must back up the entire paxi folder, especially your node's private keys (node_key.json, priv_validator_key.json, mnemonic phrase)."
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

read -p "Enter your emergency contact email: " SECURITY_CONTACT
if [[ -z "$SECURITY_CONTACT" ]]; then
  echo "❌ Emergency contact cannot be empty. Please rerun the script."
  exit 1
fi
read -p "Enter your website or contact page (can be X, Facebook, Telegram, WhatsApp, Discord, Github, etc.): " WEBSITE


PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
PERSISTENT_PEERS="9e64baa45042ae29d999f2677084c9579972722c@139.99.69.74:26656,a5d9c4fce0fe7a877595cf335135f7591dba6bf7@144.217.68.77:26656,d411fc096e0d946bbd2bdea34f0f9da928c1a714@139.99.68.32:26656,a15ff764aa634bbbca93b352d559f1b74e78af31@139.99.68.235:26656,4f26d3ecfb1aa81f2304850af20e74462ed1d341@66.70.181.55:26656,57b44498315f013558e342827f352db790d5d90c@142.44.142.121:26656,3f348c056aadccef0a5c6ea254b0a5a0096c1c7a@51.75.54.185:26656"
RPC_URL="http://rpc.paxi.info"
SNAPSHOT_URL="http://rpc.paxi.info"
SNAPSHOT_DOWNLOAD_HOST="http://snapshot.paxi.info"
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
    ca-certificates curl gnupg lsb-release git make unzip jq

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

  if [ "$EUID" -ne 0 ]; then
    echo "⚠️ You may need to log out and back in (or run 'newgrp docker') to apply Docker permissions, then run this script again."
    exit 1
  fi
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
  git checkout $PAXI_TAG
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
  $BINARY_NAME init "$NODE_MONIKER" --chain-id $CHAIN_ID
sudo chown -R $(whoami) $HOME/paxid
fi 

curl -s $GENESIS_URL | jq -r .result.genesis > ./paxi/config/genesis.json

### === Set state sync ===
BLOCK_OFFSET=1000
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
    s|^discovery_time *=.*|discovery_time = \"60s\"|g; \
  }" "$CONFIG"
fi 

### === Get network information ===
IP_DATA=$(curl -s http://ip-api.com/json)
if [ $? -ne 0 ]; then
  echo "❌ Failed to retrieve country code. Please check your internet connection."
  exit 1
fi
COUNTRY_CODE=$(echo "$IP_DATA" | jq -r .countryCode)
IP_ADDRESS=$(echo "$IP_DATA" | jq -r .query)

### === Configure seeds and peers ===
$SED "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG
# $SED "s|^[[:space:]]*external_address = \".*\"|external_address = \"${IP_ADDRESS}:26656\"|" $CONFIG

### === Disable unnecessary ports for security ===
$SED 's|^prometheus *=.*|prometheus = false|' $CONFIG

### === Create wallet (if not exists) ===
if ! [ -f ./paxi/keyring-file/$KEY_NAME.info ] && ! [ -f ./paxi/$KEY_NAME.info ]; then
  echo ""
  echo "After the wallet is created, please write down the mnemonic phrase"
  echo "manually to ensure you can recover it if lost."

  docker run --rm -it \
    -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
    $DOCKER_IMAGE \
    $BINARY_NAME keys add $KEY_NAME
fi
sudo chown -R $(whoami) $HOME/paxid

### === Show wallet address ===
echo ""
echo "Your wallet address:"
docker run --rm -it \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME keys show $KEY_NAME -a
echo "Please fund this address and run the following command to become a validator:"

### === Generate create-validator command ===
VAL_PUBKEY=$(docker run --rm \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME tendermint show-validator)
echo "You may edit validator.json at: $PAXI_DATA_PATH/validator.json"
echo "Please add your country code at the end of the 'details' parameter, e.g., [US]. This helps us collect node location data and display it on the official website."
echo "Generating validator.json..."
cat <<EOF > $PAXI_DATA_PATH/validator.json
{
  "security": "$SECURITY_CONTACT",
  "moniker": "$NODE_MONIKER",
  "identity": "",
  "website": "$WEBSITE",
  "details": "PAXI validator node [$COUNTRY_CODE]",
  "pubkey": $VAL_PUBKEY,
  "amount": "2000000000$DENOM",
  "commission-rate": "0.25",
  "commission-max-rate": "0.5",
  "commission-max-change-rate": "0.1",
  "min-self-delegation": "1"
}
EOF

### === Common commands ===
echo ""
echo "Start the node (this node runs in the background and auto-starts on system reboot unless stopped manually):"
echo "docker run -d --name paxi-node-1 --restart unless-stopped -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \\"
echo "--network=host \\"
echo "paxi-node \\"
echo "$BINARY_NAME start"
echo ""
echo "For convenience, please enter the container before running the following commands:"
echo "docker run --rm -it --network host -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH $DOCKER_IMAGE bash"
echo ""
echo "Validator creation command (copy & paste to execute):"
echo "$BINARY_NAME tx staking create-validator $DOCKER_PAXI_DATA_PATH/validator.json \\"
echo "  --from $KEY_NAME \\"
echo "  --fees 10000$DENOM"
echo ""
echo "To view your wallet address:"
echo "$BINARY_NAME keys show $KEY_NAME"
echo ""
echo "Check wallet balance:"
echo "$BINARY_NAME query bank balances <your address or key name>"
echo ""
echo "Check your staking rewards:"
echo "$BINARY_NAME query distribution rewards <your address or key name>"
echo ""
echo "Check your validator public key:"
echo "$BINARY_NAME tendermint show-validator"
echo ""
echo "Check your validator's outstanding rewards:"
echo "$BINARY_NAME query distribution validator-outstanding-rewards <your validator address>"
echo ""
echo "Send token:"
echo "$BINARY_NAME tx bank send <your address or key name> <receiver's address> <amount>$DENOM --fees 10000$DENOM --gas auto"
echo ""
echo "============================================================"
echo "❗After starting the Paxi node (paxid start), make sure to run the WASM contract synchronization script again."
echo "❗Failing to do so may cause consensus failures due to missing WASM files, which can result in your validator being slashed (i.e., a portion of your stake may be deducted)."
echo ""
echo "curl -sL https://raw.githubusercontent.com/paxi-web3/paxi/main/scripts/sync_wasm.sh | bash"
echo ""