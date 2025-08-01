#!/bin/bash
set -e

### === 需求參數 ===
REQUIRED_CPU=8
REQUIRED_RAM_GB=30
REQUIRED_DISK_GB=1000

### ===  取得 CPU cores ===
CPU_CORES=$(nproc)

### ===  取得實體記憶體 (轉 GB) ===
TOTAL_MEM=$(free -g | awk '/^Mem:/{print $2}')

### ===  取得可用磁碟空間 (轉 GB)，以 / 根目錄為例 ===
DISK_SPACE=$(df -BG / | awk 'NR==2 {gsub("G","",$4); print $4}')

### ===  檢查 CPU ===
if [ "$CPU_CORES" -lt "$REQUIRED_CPU" ]; then
  echo "❌ CPU cores不足: 需要至少 ${REQUIRED_CPU} cores，當前 ${CPU_CORES} cores"
  exit 1
fi

### ===  檢查 RAM ===
if [ "$TOTAL_MEM" -lt "$REQUIRED_RAM_GB" ]; then
  echo "❌ 記憶體不足: 需要至少 ${REQUIRED_RAM_GB} GB，當前 ${TOTAL_MEM} GB"
  exit 1
fi

### ===  檢查磁碟空間 ===
if [ "$DISK_SPACE" -lt "$REQUIRED_DISK_GB" ]; then
  echo "❌ 磁碟空間不足: 需要至少 ${REQUIRED_DISK_GB} GB，當前 ${DISK_SPACE} GB"
  exit 1
fi

echo "============================================================"
echo "🚨  PAXI 驗證人節點安裝警告"
echo "============================================================"
echo ""
echo "🛑 最強烈警告："
echo "❗ 若超過 1/3 驗證人節點掉線，整個區塊鏈將會停擺。"
echo "❗ 請務必備份整個 paxi 資料夾，尤其是節點私鑰（node_key.json、priv_validator_key.json、keyring），"
echo "   一旦電腦故障，才能及時修復並取回質押收益與驗證人身份。"
echo ""
echo "⚠️ 請注意:"
echo "   一旦你質押並成為驗證人，系統會自動監控你的上線狀態。"
echo ""
echo "❗ 如果你無故離線（斷線或關機），系統將視為懲罰性行為，"
echo "   並自動扣除你的一部分質押金（Slashing 機制）。"
echo ""
echo "✅ 正確離線方法:"
echo "   請先使用解除質押命令（Undelegate）退出驗證人角色後，再關閉節點。"
echo ""
echo "🚫 直接關機或停止節點會造成懲罰風險。請務必確認！"
echo ""
echo "============================================================"
read -p "你已了解以上風險，是否繼續安裝？ (y/N): " confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
  echo "❌ 已取消安裝。請再次閱讀說明後再啟動此腳本。"
  exit 1
fi

echo "📝 請輸入你要為節點設定的名稱（moniker）:"
read -p "節點名稱: " NODE_MONIKER

if [[ -z "$NODE_MONIKER" ]]; then
  echo "❌ 節點名稱不能為空，請重新執行腳本。"
  exit 1
fi

echo "✅ 節點名稱設定為: $NODE_MONIKER"

echo "📝 請輸入你要為你的錢包設定的名稱（key name）:"
read -p "請輸入你的錢包名稱（key name）: " KEY_NAME
if [[ -z "$KEY_NAME" ]]; then
  echo "❌ 錢包名稱不能為空，請重新執行腳本。"
  exit 1
fi
echo "✅ 錢包名稱設定為: $KEY_NAME"

read -p "請輸入您的緊急聯絡電子郵件: " SECURITY_CONTACT
if [[ -z "$SECURITY_CONTACT" ]]; then
  echo "❌ 緊急聯絡信箱不能為空，請重新執行本腳本。"
  exit 1
fi
read -p "請輸入您的網站或聯絡頁面（可以是 X / Facebook / Telegram / WhatsApp / Discor / Github 等）: " WEBSITE


GOLANG_VERSION=1.24.2
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
PERSISTENT_PEERS="25aa07b1073c49e298cabe00fcc5ccb0904fb840@51.79.176.58:26656,d411fc096e0d946bbd2bdea34f0f9da928c1a714@139.99.68.32:26656,a15ff764aa634bbbca93b352d559f1b74e78af31@139.99.68.235:26656,4f26d3ecfb1aa81f2304850af20e74462ed1d341@66.70.181.55:26656,57b44498315f013558e342827f352db790d5d90c@142.44.142.121:26656,3f348c056aadccef0a5c6ea254b0a5a0096c1c7a@51.75.54.185:26656"
RPC_URL="http://mainnet-rpc.paxinet.io"
SNAPSHOT_URL="http://mainnet-rpc.paxinet.io"
SNAPSHOT_DOWNLOAD_HOST="http://mainnet-snapshot.paxinet.io"
GENESIS_URL="$RPC_URL/genesis?"
CONFIG="./paxi/config/config.toml"
APP_CONFIG="./paxi/config/app.toml"
PAXI_PATH="$HOME/paxid"
PAXI_DATA_PATH="$HOME/paxid/paxi"
DENOM="upaxi"

### === 安裝依賴 ===
echo ""
sudo apt-get update && sudo apt-get install -y \
    build-essential git cmake \
    libsnappy-dev zlib1g-dev libbz2-dev \
    liblz4-dev libzstd-dev wget curl pkg-config \
    ca-certificates unzip jq

### === 安裝 Go ===
if ! command -v go &> /dev/null; then
    echo "正在安裝 Go..."
    curl -LO https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo ln -s /usr/local/go/bin/go /usr/bin/go
fi

### === 編譯 Paxi ===
if ! [ -d ./paxi ]; then
echo "正在安裝 Paxi..."
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

### === 複製 libwasmvm.x86_64.so 到 /usr/loca/lib ===
WASM_SO_FILE="$HOME/go/pkg/mod/github.com/!cosm!wasm/wasmvm/v2@v2.2.1/internal/api/libwasmvm.x86_64.so"
if [ -f "$WASM_SO_FILE" ]; then
  sudo cp $WASM_SO_FILE /usr/local/lib/
  sudo ldconfig
fi

### === 初始化節點 ===
if ! [ -f ./paxi/config/genesis.json ]; then
$BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
fi 

curl -s $GENESIS_URL | jq -r .result.genesis > ./paxi/config/genesis.json

### === 設置快照同步 ===
BLOCK_OFFSET=1000
LATEST_HEIGHT=$(curl -s "$RPC_URL/block" | jq -r .result.block.header.height)
TRUST_HEIGHT=$(( ( (LATEST_HEIGHT - BLOCK_OFFSET) / BLOCK_OFFSET ) * BLOCK_OFFSET ))
TRUST_HASH=$(curl -s "$RPC_URL/block?height=$TRUST_HEIGHT" | jq -r .result.block_id.hash)

if ! [[ "$LATEST_HEIGHT" =~ ^[0-9]+$ ]]; then
  echo "❌ 無法取得 trust 高度或 hash，請檢查 RPC URL。"
  exit 1
fi

### === 下載 wasm snapshot ===
WASM_SNAPSHOT_URL=$(curl -s "$SNAPSHOT_DOWNLOAD_HOST/utils/latest_wasm_snapshot" | jq -r .url)
curl -f -o wasm_snapshot.zip "$WASM_SNAPSHOT_URL"
if [ $? -ne 0 ]; then
  echo "❌ 無法下載 wasm snapshot。請手動下載並解壓到 $PAXI_DATA_PATH/wasm/wasm/state/wasm。"
else
  mkdir -p "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  unzip wasm_snapshot.zip -d "$PAXI_DATA_PATH/wasm/wasm/state/wasm"
  rm wasm_snapshot.zip
  echo "✅ Wasm snapshot 已下載並解壓到 $PAXI_DATA_PATH/wasm/wasm/state/wasm。"
fi

### === 檢測系統類別 ===
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED="sed -i ''"
else
    SED="sed -i"
fi

if [ "$LATEST_HEIGHT" -gt "$BLOCK_OFFSET" ]; then
  if [[ -z "$TRUST_HEIGHT" || -z "$TRUST_HASH" || "$TRUST_HASH" == "null" ]]; then
    echo "❌ 無法取得 trust 高度或 hash，請檢查 RPC URL。"
    exit 1
  fi

  $SED "/^\[statesync\]/,/^\[/{ \
    s|^enable *=.*|enable = true|g; \
    s|^rpc_servers *=.*|rpc_servers = \"$SNAPSHOT_URL,$SNAPSHOT_URL\"|g; \
    s|^trust_height *=.*|trust_height = $TRUST_HEIGHT|g; \
    s|^trust_hash *=.*|trust_hash = \"$TRUST_HASH\"|g; \
    s|^trust_period *=.*|trust_period = \"168h\"|g; \
  }" "$CONFIG"
fi

### === 獲取網絡資料 ===
IP_DATA=$(curl -s http://ip-api.com/json)
if [ $? -ne 0 ]; then
  echo "❌ Failed to retrieve country code. Please check your internet connection."
  exit 1
fi
COUNTRY_CODE=$(echo "$IP_DATA" | jq -r .countryCode)
IP_ADDRESS=$(echo "$IP_DATA" | jq -r .query)

### === 設定種子與peers ===
$SED "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG
$SED "s|^[[:space:]]*external_address = \".*\"|external_address = \"${IP_ADDRESS}:26656\"|" $CONFIG

### === 關閉不必要端口，強化安全性 ===
$SED 's|^prometheus *=.*|prometheus = false|' $CONFIG

### === 建立錢包（如不存在）===
if ! $BINARY_NAME keys show $KEY_NAME &>/dev/null; then
  echo ""
  echo "錢包創建完成後，請用手寫的方式記下以下的所有助記詞，以便遺失時恢復你的錢包"
  $BINARY_NAME keys add $KEY_NAME
fi

### === 顯示地址 ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a)
echo ""
echo "你的地址為: $ADDR"
echo "請向此地址轉入代幣後執行以下指令進行質押:"

### === 顯示 create-validator 指令 ===
VAL_PUBKEY=$($BINARY_NAME tendermint show-validator)
echo "你可以從 $PAXI_DATA_PATH/validator.json 修改參數"
echo "請在 'details' 參數的最後加上你的國家代號，例如: [US]，此舉方便我們收集節點位置數據然後顯示在官網上"
echo "正在產生 validator.json..."
cat <<EOF > $PAXI_DATA_PATH/validator.json
{
  "pubkey": $VAL_PUBKEY,
  "amount": "1000000000$DENOM",
  "moniker": "$NODE_MONIKER",
  "identity": "",
  "website": "$WEBSITE",
  "security-contact": "$SECURITY_CONTACT",
  "details": "PAXI validator node [$COUNTRY_CODE]",
  "commission-rate": "0.25",
  "commission-max-rate": "0.5",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF
echo ""
echo "成為驗證人指令（複製貼上執行）:"
echo "cd $PAXI_PATH && $BINARY_NAME tx staking create-validator $PAXI_DATA_PATH/validator.json \\"
echo "  --from $KEY_NAME \\"
echo "  --fees 10000$DENOM"

### === 常用指令 ===
echo ""
echo "啓動節點前記得設置'你的公網IP:26656'到 paxi/config/config.toml 的 'external_address' 參數，不然其他人無法連接到你的節點"
echo "啓動節點:"
echo "$BINARY_NAME start"
echo ""
echo "查看錢包餘額指令:"
echo "$BINARY_NAME query bank balances <你的地址/錢包名稱>"
echo ""
echo "查看你的質押收益:"
echo "$BINARY_NAME query distribution rewards <你的地址/錢包名稱>"
echo ""
echo "查看你的驗證人地址指令:"
echo "$BINARY_NAME tendermint show-validator"
echo ""
echo "查看你的驗證人收益:"
echo "$BINARY_NAME query distribution validator-outstanding-rewards <你的驗證人地址>"
echo ""
echo "發送代幣:"
echo "$BINARY_NAME tx bank send <你的錢包名稱> <接收地址> <數量>$DENOM --fees 10000$DENOM --gas auto"
echo ""
echo "============================================================"
echo "❗當你啟動完 Paxi 節點（執行 paxid start）後，務必要再次執行一次 WASM 智能合約同步指令（bash sync_wasm_zh.sh）。"
echo "❗否則若缺失了某些 WASM 文件，將會導致共識錯誤，進而使你的驗證人節點被系統懲罰（例如被扣除一部分質押資產）。"
echo ""
echo "curl -sL https://raw.githubusercontent.com/paxi-web3/paxi/scripts/sync_wasm_zh.sh | bash"
echo ""
