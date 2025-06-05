#!/bin/bash
set -e

# 當你啟動完 Paxi 節點（執行 paxid start）後，務必要再次執行一次 WASM 智能合約同步指令（bash sync_wasm_zh.sh）。
# 這個步驟會從其他節點下載完整的智能合約文件。
# 否則若缺失了某些 WASM 文件，將會導致共識錯誤，進而使你的驗證人節點被系統懲罰（例如被扣除一部分質押資產）。

sudo apt-get update && sudo apt-get install -y \
    curl unzip jq

SNAPSHOT_DOWNLOAD_HOST="http://mainnet-snapshot.paxinet.io"
PAXI_DATA_PATH="$HOME/paxid/paxi"

### === Download wasm snapshot ===
WASM_SNAPSHOT_URL=$(curl -s "$SNAPSHOT_DOWNLOAD_HOST/utils/latest_wasm_snapshot" | jq -r .url)
curl -f -o wasm_snapshot.zip "$WASM_SNAPSHOT_URL"
if [ $? -ne 0 ]; then
  echo "❌ 無法下載 wasm snapshot。請手動下載並解壓到 $PAXI_DATA_PATH/wasm/wasm/state/wasm。"
else
  mkdir -p "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  sudo chmod 777 "$PAXI_DATA_PATH/wasm/wasm/state/wasm"	
  unzip wasm_snapshot.zip -d "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  rm wasm_snapshot.zip
  echo "✅ Wasm snapshot 已下載並解壓到 $PAXI_DATA_PATH/wasm/wasm/state/wasm。"
fi