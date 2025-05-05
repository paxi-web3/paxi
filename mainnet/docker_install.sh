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

PAXI_REPO="https://github.com/paxi-web3/paxi"
PAXI_TAG="v1.0.1"
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
DOCKER_IMAGE="paxi-node"
DOCKER_PAXI_DATA_PATH="/root/paxi"

### === 安裝依賴 ===
echo ""
sudo apt-get update
sudo apt-get install -y \
    ca-certificates curl gnupg lsb-release git make

### === 安裝 Docker ===
if ! command -v docker &> /dev/null; then
  echo "安裝 Docker 中..."
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

  # 啟動 Docker 並允許非 root 用戶執行
  sudo systemctl enable docker
  sudo systemctl start docker
  sudo usermod -aG docker $USER

  echo "⚠️ 你可能需要重新登入，讓 docker 權限生效（或執行 newgrp docker）"
else
  echo "✅ Docker 已安裝"
fi


### === 安裝 Paxi ===
if [ ! -d "paxi" ]; then
  git clone $PAXI_REPO
  cd paxi
  git checkout $PAXI_TAG
  make docker
else
  cd paxi
  make docker
fi

if ![ -d "$HOME/paxid" ]; then
mkdir "$HOME/paxid" 
fi
cd $HOME/paxid 

### === 初始化節點 ===
if ! [ -f ./paxi/config/genesis.json ]; then
docker run --rm \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  $BINARY_NAME init $NODE_MONIKER --chain-id $CHAIN_ID
curl -L $GENESIS_URL > ./paxi/config/genesis.json
fi 

### === 設定種子與peers ===
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CONFIG
sed -i "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CONFIG

### === 關閉不必要端口，強化安全性 ===
sed -i 's|^laddr *=.*|laddr = "tcp://127.0.0.1:26657"|' $CONFIG
sed -i 's|^prometheus *=.*|prometheus = false|' $CONFIG
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[api\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^enable *=.*|enable = false|' $(grep -l "\[grpc-web\]" $APP_CONFIG -A 3 | tail -n 1)
sed -i 's|^address *=.*|address = "127.0.0.1:9090"|' $(grep -l "\[grpc\]" $APP_CONFIG -A 3 | tail -n 1)

### === 建立錢包（如不存在）===
# 檢查 key 是否已存在（使用 docker run 執行 paxid keys show）
if ! docker run --rm \
  -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
  $DOCKER_IMAGE \
  -it \ 
  ./paxid keys show $KEY_NAME --keyring-backend os &>/dev/null; then

  echo ""
  echo "錢包創建完成後，請用手寫的方式記下以下的所有助記詞，以便遺失時恢復你的錢包"

  docker run --rm \
    -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \
    -it \
    $DOCKER_IMAGE \
    ./paxid keys add $KEY_NAME --keyring-backend os
fi

### === 顯示地址 ===
ADDR=$($BINARY_NAME keys show $KEY_NAME -a --keyring-backend os)
echo ""
echo "你的地址為: $ADDR"
echo "請向此地址轉入代幣後執行以下指令進行質押:"

### === 顯示 create-validator 指令 ===
VAL_PUBKEY=$($BINARY_NAME tendermint show-validator)
echo "你可以從 $PAXI_PATH/validator.json 修改參數"
echo "正在產生 validator.json..."
cat <<EOF > validator.json
{
  "pubkey": $VAL_PUBKEY,
  "amount": "1000000000$DENOM",
  "moniker": "$NODE_MONIKER",
  "identity": "",
  "website": "",
  "security": "",
  "details": "PAXI validator initialized by install.sh",
  "commission-rate": "0.1",
  "commission-max-rate": "0.2",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF
echo ""
echo "成為驗證人指令（複製貼上執行）:"
echo ""
echo "docker run --rm -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH -it $DOCKER_IMAGE \\"
echo "$BINARY_NAME tx staking create-validator ./paxi/validator.json \\"
echo "  --from $KEY_NAME --keyring-backend os \\"
echo "  --fees 10000$DENOM"

### === 常用指令 ===
echo ""
echo "啓動節點(這節點會在後台運行，在電腦啓動後它也會自動啓動，除非你手動關停):"
echo "docker run -d --name paxi-node-1 --restart unless-stopped -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH \\"
echo "-p 26656:26656 -p 26657:26657 -p 1317:1317 -p 9090:9090 \\"
echo "paxi-node \\"
echo "$BINARY_NAME start"
echo ""
echo "查看錢包餘額指令:"
echo "docker run --rm -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH -it $DOCKER_IMAGE \\"
echo "$BINARY_NAME query bank balances <你的地址/錢包名稱>"
echo ""
echo "查看你的質押收益:"
echo "docker run --rm -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH -it $DOCKER_IMAGE \\"
echo "$BINARY_NAME query distribution rewards <你的地址/錢包名稱>"
echo ""
echo "查看你的驗證人地址指令:"
echo "docker run --rm -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH -it $DOCKER_IMAGE \\"
echo "$BINARY_NAME tendermint show-validator"
echo ""
echo "查看你的驗證人收益:"
echo "docker run --rm -v $PAXI_DATA_PATH:$DOCKER_PAXI_DATA_PATH -it $DOCKER_IMAGE \\"
echo "$BINARY_NAME query distribution validator-outstanding-rewards <你的驗證人地址>"