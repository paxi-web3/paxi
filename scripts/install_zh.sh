#!/bin/bash
set -e

echo "============================================================"
echo "🚨  PAXI 驗證人節點安裝警告"
echo "============================================================"
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

GOLANG_VERSION=1.24.2
ROCKSDB_VERSION=v9.2.1 
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="latest-main"
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

### === 安裝依賴 ===
echo ""
sudo apt update
sudo apt-get update && sudo apt-get install -y \
    build-essential git cmake \
    libsnappy-dev zlib1g-dev libbz2-dev \
    liblz4-dev libzstd-dev wget curl pkg-config \
    ca-certificates 

### === 安裝 Go ===
if ! command -v go &> /dev/null; then
    echo "正在安裝 Go..."
    curl -LO https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    sudo ln -s /usr/local/go/bin/go /usr/bin/go
fi

### === 安裝 Rocksdb ===
if ! [ -f /usr/local/lib/librocksdb.so ]; then
    echo "正在安裝 Rocksdb..."
    
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
make install
cd $PAXI_PATH
fi

### === 初始化節點 ===
if ! [ -f ./paxi/config/genesis.json ]; then
$BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
curl -L $GENESIS_URL > ./paxi/config/genesis.json
fi 

### === 設定種子與peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CONFIG
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === 關閉不必要端口，強化安全性 ===
sed -i '/^\[rpc\]/,/^\[/s|^\s*laddr\s*=.*|laddr = "tcp://0.0.0.0:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === 建立錢包（如不存在）===
if ! $BINARY_NAME keys show $KEY_NAME --keyring-backend os &>/dev/null; then
  echo ""
  echo "錢包創建完成後，請用手寫的方式記下以下的所有助記詞，以便遺失時恢復你的錢包"
  $BINARY_NAME keys add $KEY_NAME --keyring-backend os
fi

### === 顯示地址 ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a --keyring-backend os)
echo ""
echo "你的地址為: $ADDR"
echo "請向此地址轉入代幣後執行以下指令進行質押:"

### === 顯示 create-validator 指令 ===
COUNTRY_CODE=$(curl -s http://ip-api.com/json | jq .countryCode)
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
echo "成為驗證人指令（複製貼上執行）:"
echo "cd $PAXI_PATH && $BINARY_NAME tx staking create-validator ./paxi/validator.json \\"
echo "  --from $KEY_NAME --keyring-backend os \\"
echo "  --fees 10000$DENOM"

### === 常用指令 ===
echo ""
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