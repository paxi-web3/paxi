#!/usr/bin/env bash
set -euo pipefail

# ====== Config ======
SNAP_RPC="${SNAP_RPC:-http://rpc.paxi.info/}"
SNAP_RPC2="${SNAP_RPC2:-http://rpc.paxi.info/}"
CONFIG_FILE="${CONFIG_FILE:-$HOME/paxid/paxi/config/config.toml}"

# ====== Preflight checks ======
need() { command -v "$1" >/dev/null 2>&1 || { echo "ERROR: '$1' not found. Please install it."; exit 1; }; }
need curl
need jq
need sed

fetch_json() {
  # $1 = url
  # Return JSON; if not JSON or HTTP failure, throw error
  local url="$1"
  local resp
  resp=$(curl -fsSL -H 'Accept: application/json' --connect-timeout 10 --max-time 20 "$url") || {
    echo "ERROR: curl failed for $url"
    exit 1
  }
  # Verify it's JSON
  echo "$resp" | jq -e . >/dev/null 2>&1 || {
    echo "ERROR: response from $url is not valid JSON (possibly HTML/503/redirect)."
    echo "Tip: open $url in a browser to see what it returns."
    exit 1
  }
  echo "$resp"
}

# ====== Get latest height from /status (more reliable) ======
STATUS_JSON=$(fetch_json "${SNAP_RPC%/}/status")
LATEST_HEIGHT=$(echo "$STATUS_JSON" | jq -r '.result.sync_info.latest_block_height // empty')
if [[ -z "${LATEST_HEIGHT:-}" || ! "$LATEST_HEIGHT" =~ ^[0-9]+$ ]]; then
  echo "ERROR: failed to parse latest_block_height from /status"
  exit 1
fi

# Use LATEST_HEIGHT - 100 as trust height, but minimum is 2
if (( LATEST_HEIGHT > 150 )); then
  BLOCK_HEIGHT=$((LATEST_HEIGHT - 100))
else
  BLOCK_HEIGHT=$((LATEST_HEIGHT > 2 ? LATEST_HEIGHT - 1 : 2))
fi

# ====== Get trust hash for that height ======
BLOCK_JSON=$(fetch_json "${SNAP_RPC%/}/block?height=${BLOCK_HEIGHT}")
TRUST_HASH=$(echo "$BLOCK_JSON" | jq -r '.result.block_id.hash // empty')

if [[ -z "${TRUST_HASH:-}" || ! "$TRUST_HASH" =~ ^[A-F0-9]{64}$ ]]; then
  echo "ERROR: failed to get a valid trust_hash for height ${BLOCK_HEIGHT}."
  exit 1
fi

echo "Latest:       $LATEST_HEIGHT"
echo "Trust Height: $BLOCK_HEIGHT"
echo "Trust Hash:   $TRUST_HASH"

# ====== Verify config path ======
if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "ERROR: config not found at: $CONFIG_FILE"
  echo "Please set CONFIG_FILE env or create the file. Common locations:"
  echo "  - $HOME/paxid/paxi/config/config.toml"
  echo "  - $HOME/.paxi/config/config.toml"
  exit 1
fi

# ====== Patch config.toml ======
# Backup
cp -a "$CONFIG_FILE" "${CONFIG_FILE}.bak"

# Modify (allow flexible spacing)
sed -i -E \
  -e "s|^([[:space:]]*enable[[:space:]]*=[[:space:]]*).*$|\1true|g" \
  -e "s|^([[:space:]]*rpc_servers[[:space:]]*=[[:space:]]*).*$|\1\"${SNAP_RPC%/},${SNAP_RPC2%/}\"|g" \
  -e "s|^([[:space:]]*trust_height[[:space:]]*=[[:space:]]*).*$|\1${BLOCK_HEIGHT}|g" \
  -e "s|^([[:space:]]*trust_hash[[:space:]]*=[[:space:]]*).*$|\1\"${TRUST_HASH}\"|g" \
  -e "s|^([[:space:]]*seeds[[:space:]]*=[[:space:]]*).*$|\1\"\"|g" \
  "$CONFIG_FILE"

echo "Updated: $CONFIG_FILE (backup: ${CONFIG_FILE}.bak)"

# ====== Tips ======
cat <<EOF
Next:
  1) Make sure the paxi service is stopped and then restarted (to avoid race conditions):
     sudo systemctl restart paxi   # If your service name is different, change it accordingly
  2) Check logs to confirm state sync is working:
     journalctl -u paxi -f
EOF
