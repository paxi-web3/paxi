#!/usr/bin/env bash
set -euo pipefail

############################################
# PAXI CLI installer – paxid only
# - Build & install paxid
# - Let user manually choose mainnet/testnet
# - Write chain-id & node RPC into client.toml
############################################

GO_VERSION="1.24.2"
REPO="https://github.com/paxi-web3/paxi"
TAG="latest-main"
BIN_NAME="$HOME/paxid/paxid"

# Network presets
MAINNET_CHAIN_ID="paxi-mainnet"
MAINNET_RPC="tcp://mainnet-rpc.paxinet.io:80"

TESTNET_CHAIN_ID="paxi-testnet"
TESTNET_RPC="tcp://testnet-rpc.paxinet.io:80"

############################################
choose_network () {
  echo "Select network:"
  echo "  1) mainnet  (${MAINNET_CHAIN_ID})"
  echo "  2) testnet  (${TESTNET_CHAIN_ID})"
  read -rp "Enter 1 or 2 [default:1]: " choice
  if [[ "${choice:-1}" == "2" ]]; then
    SELECTED_NET="testnet"
  else
    SELECTED_NET="mainnet"
  fi
}

set_client_config () {
  local chain_id="$1"
  local rpc="$2"

  # Init node
  [ -s "$HOME/paxid/paxi/config/genesis.json" ] || \
  "$BIN_NAME" init my-node-0 --chain-id my-testnet

  # Use official CLI config commands (Cosmos SDK style)
  "$BIN_NAME" config set client chain-id "$chain_id"
  "$BIN_NAME" config set client node "$rpc"

  # Try to locate client.toml for info
  local default_home="${BIN_NAME}"
  local config_home
  if "$BIN_NAME" config home >/dev/null 2>&1; then
    config_home=$("$BIN_NAME" config home)
  else
    [[ -d "$HOME/.paxi" ]]  && config_home="$HOME/.paxi"
    [[ -d "$HOME/.paxid" ]] && config_home="$HOME/.paxid"
    [[ -z "${config_home:-}" ]] && config_home="$default_home"
  fi

  local client_toml="$config_home/config/client.toml"
  echo "✅ Updated client.toml:"
  echo "   chain-id = $chain_id"
  echo "   node     = $rpc"
  echo "   file     = $client_toml"
}

install_go () {
  if command -v go >/dev/null 2>&1; then
    echo "Go found: $(go version)"
    return
  fi
  echo "Installing Go ${GO_VERSION}..."
  wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
  sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go
  rm "go${GO_VERSION}.linux-amd64.tar.gz"
  echo "Go installed: $(go version)"
}

install_deps () {
  sudo apt-get update
  sudo apt-get install -y build-essential git wget curl jq unzip ca-certificates pkg-config
}

build_paxid () {
  if [[ ! -d paxi ]]; then
    git clone "$REPO"
    cd paxi
  else
    cd paxi
    git fetch --all
  fi
  git checkout "$TAG"
  make install
  cd -

  echo "✅ Installed $(command -v $BIN_NAME)"
  $BIN_NAME version || true
}

############################################
# Main
############################################
install_deps
install_go
build_paxid
choose_network

if [[ "$SELECTED_NET" == "mainnet" ]]; then
  set_client_config "$MAINNET_CHAIN_ID" "$MAINNET_RPC"
else
  set_client_config "$TESTNET_CHAIN_ID" "$TESTNET_RPC"
fi

echo ""
echo "Done! Examples:"
echo "  $BIN_NAME status"
echo "  $BIN_NAME query bank balances <address>"
echo "  $BIN_NAME tx bank send <from> <to> 100000upaxi --fees 10000upaxi --gas auto"
echo ""
echo "Switch network later:"
echo "  $BIN_NAME config set client chain-id $MAINNET_CHAIN_ID   # or $TESTNET_CHAIN_ID"
echo "  $BIN_NAME config set client node $MAINNET_RPC            # or $TESTNET_RPC"
