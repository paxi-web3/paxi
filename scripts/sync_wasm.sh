#!/bin/bash
set -e

# After starting the Paxi node (paxid start), make sure to run the WASM contract synchronization script again (bash sync_wasm.sh).
# This step downloads the full set of smart contract files from other peers.
# Failing to do so may cause consensus failures due to missing WASM files, which can result in your validator being slashed (i.e., a portion of your stake may be deducted).

sudo apt-get update && sudo apt-get install -y \
    curl unzip jq

SNAPSHOT_DOWNLOAD_HOST="http://snapshot.paxi.info"
PAXI_DATA_PATH="$HOME/paxid/paxi"

### === Download wasm snapshot ===
WASM_SNAPSHOT_URL=$(curl -s "$SNAPSHOT_DOWNLOAD_HOST/utils/latest_wasm_snapshot" | jq -r .url)
curl -f -o wasm_snapshot.zip "$WASM_SNAPSHOT_URL"
if [ $? -ne 0 ]; then
  echo "❌ Failed to download wasm snapshot. Please download it and unzip it to $PAXI_DATA_PATH/wasm/wasm/state/wasm."
else
  mkdir -p "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  sudo chmod 777 "$PAXI_DATA_PATH/wasm/wasm/state/wasm"	
  unzip wasm_snapshot.zip -d "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  rm wasm_snapshot.zip
  echo "✅ Wasm snapshot downloaded and extracted to $PAXI_DATA_PATH/wasm/wasm/state/wasm."
fi
