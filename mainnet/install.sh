#!/bin/bash
set -e

echo "============================================================"
echo "🚨  PAXI 驗證人節點安裝警告"
echo "============================================================"
echo ""
echo "⚠️  請注意:"
echo "一旦你質押並成為驗證人，系統會自動監控你的上線狀態。"
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
echo "當節點安裝完畢後，此程序將自動幫你創建本地錢包，屆時請記下所有助記詞"

GOLANG_VERSION=1.24.2
ROCKSDB_VERSION=v9.2.1 
CHAIN_ID="paxi-mainnet"
BINARY_NAME="./paxid"
PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="v1.0.1"
GENESIS_URL="https://raw.githubusercontent.com/paxi-web3/mainnet/genesis.json"
SEEDS="mainner-seed-1.paxi.io:26656"
PERSISTENT_PEERS="key@mainnet-node-1.paxi.io:26656"
CONFIG="./config/config.toml"
APP_CONFIG="./config/app.toml"
PAXI_PATH="$HOME/paxid"
DENOM="upaxi"

### === 安裝依賴 ===
sudo apt update
sudo apt-get update && apt-get install -y \
    build-essential git cmake \
    libsnappy-dev zlib1g-dev libbz2-dev \
    liblz4-dev libzstd-dev wget curl pkg-config \
    ca-certificates 

### === 安裝 Go ===
if ! command -v go &> /dev/null; then
    echo "正在安裝 Go..."
    curl -LO https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go
fi

### === 安裝 Rocksdb ===
if ! [ -f /usr/local/lib/librocksdb.so ]; then
    echo "正在安裝 Rocksdb..."
    git clone https://github.com/facebook/rocksdb.git && cd rocksdb
    git checkout $ROCKSDB_VERSION
    make -j$(nproc) shared_lib
    sudo make install-shared INSTALL_PATH=/usr/local
    echo "/usr/local/lib" | sudo tee /etc/ld.so.conf.d/rocksdb.conf
    sudo ldconfig && cd ..
fi

### === 編譯 Paxi ===
echo "正在安裝 Paxi..."
git clone $PAXI_REPO
cd paxi
git checkout $PAXI_TAG
make install
cd $PAXI_PATH

### === 初始化節點 ===
$BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
curl -L $GENESIS_URL > ./config/genesis.json

### === 設定種子與peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $HOME/.paxi/config/config.toml
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $HOME/.paxi/config/config.toml

### === 關閉不必要端口，強化安全性 ===
sed -i 's|^laddr *=.*|laddr = "tcp://127.0.0.1:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === 建立錢包（如不存在）===
if ! $BINARY_NAME keys show $KEY_NAME --keyring-backend os &>/dev/null; then
  $BINARY_NAME keys add $KEY_NAME --keyring-backend os
fi

### === 顯示地址 ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a --keyring-backend os)
echo ""
echo "你的地址為: $ADDR"
echo "請向此地址轉入代幣後執行以下指令進行質押:"

### === 顯示 create-validator 指令 ===
echo ""
echo "成為驗證人指令（複製貼上執行）:"
echo "$BINARY_NAME tx staking create-validator \\"
echo "  --amount 1000000$DENOM \\"
echo "  --pubkey \$($BINARY_NAME tendermint show-validator) \\"
echo "  --moniker \"$NODE_MONIKER\" \\"
echo "  --chain-id $CHAIN_ID \\"
echo "  --from $KEY_NAME \\"
echo "  --commission-rate 0.10 \\"
echo "  --commission-max-rate 0.20 \\"
echo "  --commission-max-change-rate 0.01 \\"
echo "  --min-self-delegation 1 \\"
echo "  --keyring-backend os \\"
echo "  --fees 10000$DENOM \\"
echo "  -y"

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