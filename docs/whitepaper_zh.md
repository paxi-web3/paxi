# Paxi 白皮書

## 簡介

**Paxi** 是一個新世代的區塊鏈協議，其核心理念是簡單、快速與徹底去中心化。它建立在 Cosmos SDK 與 CometBFT 共識引擎之上，並針對效能與安全性進行了深度優化，使用 Go 語言實作，致力於提供高效、低手續費、且對使用者友善的區塊鏈基礎設施，對所有人開放。

Paxi 也支援 CosmWasm，一個基於 Wasm 的智慧合約平台，使開發者能以 Rust 語言撰寫安全且高效的合約，非常適合用於 DeFi、DAO 與跨鏈應用。

Paxi 堅信「少即是多」的理念——每一個功能都以實用與效率為目的而設計。透過專注於必要的核心功能，並避免不必要的複雜性，Paxi 為開發者、驗證者與終端使用者提供流暢無縫的使用體驗。

---

## 願景

Paxi 的願景是建立一個**快速、安全且真正去中心化**的區塊鏈生態系統，不僅僅屬於少數權勢者或大型機構，而是屬於**每一位用戶、開發者與參與者**的開放網路。在這裡，任何人都能透過質押參與共識、透過 DAO 治理提案參與決策，並透過開發貢獻創新，推動整個生態系的演進。

我們相信，唯有社群共治的力量，才能實現真正自由且可持續的未來。Paxi 將持續推動透明、公平且開放的治理機制，讓每一位參與者都能在這個平台上發聲、共創、共建。

---

## 核心理念

- **簡單性**：極簡設計、清晰 API、無冗餘邏輯
- **高速性**：優化效能與吞吐量，支持上千筆交易/秒
- **去中心化**：開放驗證節點、低門檻參與、社群治理
- **安全性**：基於 BFT 共識與簡潔可審計代碼
- **可接近性**：對所有開發者友好，包含無/低程式碼用戶

---

## 技術架構

- **Cosmos SDK**：模組化的區塊鏈應用框架
- **CometBFT**：BFT 共識引擎，提供快速最終性與安全的網路
- **CosmWasm**：一個基於 WebAssembly（Wasm）的智能合約平台，允許開發者使用 Rust 語言構建安全且高效的合約
- **Go 語言**：以開發快速、語法簡潔、併發效能著稱，適合構建高效鏈基礎設施

---

## 驗證者參與

Paxi 通過降低技術與經濟門檻，讓任何人都能成為驗證節點，推動更開放、民主的共識機制。

- 無需高額投入，**最低僅需質押 1,000 PAXI 即可申請成為驗證者**  
- 鼓勵地理與人口多樣性  
- 內建獎勵與穩定性誘因  
- 用戶也可選擇將代幣質押至其他驗證者節點，成為 **委託人 (delegator)** 並獲得質押收益  

為了進一步提升**去中心化程度、公平性與網路安全性**，Paxi 實施了一套自定的驗證者選擇機制：

- 每隔固定區塊數，自動刷新驗證者集合  
- **50% 的驗證者由最高投票權重（Top N）選出**  
- **剩餘 50% 則從其他候選者中，依其質押比例進行加權隨機抽選**

這種混合模型在確保高質量驗證的同時，也給予中小型節點公平的參與機會，進而強化網絡的抗審查性與韌性。

---

## 通脹機制

Paxi 的原生代幣初始發行量為 **1 億枚**，採用逐年遞減的通脹模型，以激勵質押參與、保障網絡安全並維持供應穩定。

### 年度通脹規則
- **第 1 年**：通脹率最高為 **8%**
- **第 2 年**：通脹率固定為 **4%**
- **第 3 年起**：通脹率不超過 **2%**

### 區塊獎勵分配原則
- 所有新鑄造代幣（區塊獎勵）**95% 分配給質押參與者 5% 分配給DAO社區** 
- 包括驗證人（Validators）與委託人（Delegators）
- 系統會自動根據供需情況進行監控。當手續費超過一定的門檻時，**系統將有機率地銷毀部分代幣**，以防止過度通脹，並降低代幣集中於少數人手中的風險

### 注意事項：低質押率的獎勵集中現象
以第一年為例：
- 年通膨：8% → 鑄造 **800 萬枚**
- 若僅有 **8%（800 萬枚）參與質押**
- 則這 **800 萬枚獎勵** 將被這少數質押者 **獨占分配**

### 這種設計能形成有效的激勵機制：
- 鼓勵更多人質押參與  
- 提高網絡安全性與共識活躍度  
- 同時控制長期供應與避免價值過度稀釋

---

## 代幣分配

| 類別                        | 比例   | 解鎖方式                                        | 資金用途與說明 |
|-----------------------------|--------|------------------------------------------------|----------------|
| **創始團隊與顧問**           | 15%    | 上線時釋放 3%，其餘分 6 次，每 4 個月釋放一階段 | 激勵核心成員與顧問長期參與與貢獻，防止短期拋售，確保團隊與項目成功高度綁定。 |
| **Paxi 基金會**             | 10%    | 上線時釋放 4%，其餘分 6 次，每 4 個月釋放一階段 | 用於鏈的維護、法務、品牌行銷與全球推廣活動，保障生態穩定長遠發展。 |
| **Paxi DAO**               | 5%     | 上線時全部釋放                                | 完全由社群治理投票決定使用方向，用於支持開發、工具建設與去中心化治理推進。 |
| **私募與戰略投資**           | 15%    | 上線時釋放 3%，其餘分 6 次，每 4 個月釋放一階段 | 吸引長期合作夥伴，推動基礎設施、技術協作、市場擴張與交易所合作。 |
| **公開發售**                | 45%    | 全額釋放或依策略分階段解鎖                          | 提升流通性與用戶參與度，資金用於掛牌、流動性建設、市場行銷與生態孵化支持。 |
| **用戶激勵與推廣**           | 10%    | 任務下載、邀請推薦、社群活動等動態發放              | 促進用戶成長與社群活躍，透過任務、推薦、公會獎勵激勵自然擴張與品牌滲透。 |

---

## Paxi DAO（去中心化自治組織）

DAO（Decentralized Autonomous Organization），即去中心化自治組織，是透過智能合約與區塊鏈技術實現的無中心化管理的組織結構。不同於傳統組織，DAO 沒有單一領導或中央權威，成員透過治理代幣投票，共同決定組織的運營、資金使用、規則制定與未來發展方向。

### Paxi DAO 的運作原理
- **智能合約**：DAO 的所有規則與流程均寫入區塊鏈上的智能合約，公開透明且不可篡改。
- **治理代幣**：成員持有治理代幣即可獲得投票權利，投票權重通常與持有代幣數量成正比。
- **提案與投票**：成員可以提出提案（如資金分配、參數修改、專案支持），經由社群投票通過後，智能合約自動執行相應操作。

### Paxi DAO 的主要功能
- **治理決策**：透過社群投票決定區塊鏈的核心參數設定（如通膨率、交易手續費、質押 規則）。
- **資金管理**：管理社區資金池，決定如何分配資源（如資助項目開發、社區推廣活動）。
- **協作平台**：提供透明、公平且可追蹤的合作環境，促進成員間的信任與合作。

### 應用範例
- **鏈上參數變更**：如投票調整通膨比例、降低治理門檻。
- **軟體升級提案**：經由 DAO 決策區塊鏈節點軟體版本升級，以確保鏈的穩定性與功能提升。
- **社區資金運用**：如投票決定是否資助開發新的 dApp 或贊助生態建設專案。
- **權限變更管理**：如授予或移除某智能合約治理權限，確保去中心化且透明的權限管理。

透過 DAO，Paxi 可真正實現去中心化自治，讓整個生態系統的發展方向完全交由社區共同決定，提升生態系統的可持續性與社區參與感。

---

## 應用場景

Paxi 適合多種實際應用場景，包括但不限於：

- **DeFi**：低手續費、高速的去中心化金融基礎設施
- **GameFi**：可擴展的鏈上遊戲與 NFT 應用平台
- **社交與身份應用**：去信任化的社群平台、認證系統等

---

## AMM 系統：PAXI ↔ PRC20 兌換機制

Paxi 內建一套高效、模組化的 AMM（自動做市商）系統，允許用戶與開發者自由在 PAXI 與 PRC20 代幣之間進行兌換，無需依賴中心化交易所。該系統參考 Uniswap V2 設計原則，並進行針對性優化，更適合鏈上即時交易與低成本交換。

### 系統特性

- **原生模組實作**：AMM 系統並非以智能合約構建，而是作為 Paxi 原生模組開發，提供極致效能與低 gas 成本。
- **穩定交換曲線**：使用恆定乘積模型（x \* y = k），確保每筆交易價格透明、可預測。
- **開放式流動性提供**：任何用戶均可在任意 PRC20 合約與 PAXI 間創建池，並提供雙邊資產以獲取流動性獎勵。
- **即時定價**：根據池內資產比例即時報價，無需外部預言機。
- **低滑價與手續費**：針對小型交易優化曲線，內建 Swap 手續費為 0.4%（可由 DAO 調整）。

### 用戶功能

- **交換（Swap）**：用戶可在任意 PRC20 ↔ PAXI 間交換，操作簡單、即時成交。
- **提供流動性（Provide Liquidity）**：使用者提供等價的 PAXI 與 PRC20，即可鑄造 LP 代幣，並享有交易費收益。
- **收益提領（Claim Rewards）**：每位 LP 可按其持有比例提領累積的 Swap 交易手續費獎勵。
- **撤回流動性（Withdraw）**：隨時提取本金與已累積的收益，無需綁定鎖倉。

### 系統設計目標

- 降低用戶參與 DeFi 門檻
- 支持開發者打造無中介的金融應用（DEX、收益農場、穩定幣協議等）
- 作為 DAO、GameFi 與 NFT 流動性支持的核心模組

### 示例用途

- 發行 PRC20 代幣後即可建立流動池，提供公開兌換渠道
- DAO 資金池可使用 AMM 池自動兌換項目所需代幣
- 遊戲內代幣經濟透過 AMM 系統實現流通與定價

---

## 開發者體驗

Paxi 致力於降低區塊鏈應用開發的門檻，將提供：

- 專屬智能合約 IDE（整合開發環境）
- 直觀的 SDK 與 API
- 無碼／低碼智能合約功能
- 完善的文件與教學資源

即使沒有程式基礎的使用者，也能在 Paxi 上部署去中心化應用（dApp）。

---

## 結語

Paxi 不只是另一個區塊鏈，而是一種有原則的方法，專門解決當前第一層網絡所面臨的複雜性與低效率問題。透過堅持簡潔、高效與去中心化的核心價值，Paxi 致力於重新定義真正以用戶為核心的區塊鏈應有的樣貌。

我們誠摯邀請開發者、驗證者、創作者與夢想家加入 Paxi 的行動。

**少即是多，精簡而成。就在 Paxi 上構建未來。**