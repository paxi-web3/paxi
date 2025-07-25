## PRC‑721 Token Standard (Code ID: 2)

**Source code**: [https://github.com/paxi-web3/cw-plus](https://github.com/paxi-web3/cw-plus)

### Overview

PRC‑721 implements the standard non‑fungible token (NFT) interface from CW721. Each token is unique, with its own owner and optional metadata. No freeze functionality is added—behaves exactly like CW721.

---

### InstantiateMsg

```jsonc
{
  "name": "MyNFTCollection",
  "symbol": "MNFT",
  "minter": "paxi1yourminteraddress"
}
```

* `name` & `symbol` describe the collection.
* `minter` is the only address allowed to mint new NFTs.

---

### ExecuteMsg

All of CW721’s execute messages:

* **Mint** (minter only)

  ```json
  { "mint": { "token_id": "unique_id", "owner": "paxi1...", "token_uri": "https://..." } }
  ```
* **Transfer**

  ```json
  { "transfer_nft": { "recipient": "paxi1...", "token_id": "unique_id" } }
  ```
* **Send** (with callback)

  ```json
  { "send_nft": { "contract": "receiver_addr", "token_id": "unique_id", "msg": "<base64>" } }
  ```
* **Approve** / **Revoke** single token

  ```json
  { "approve":  { "spender": "paxi1...", "token_id": "unique_id" } }
  { "revoke":   { "spender": "paxi1...", "token_id": "unique_id" } }
  ```
* **ApproveAll** / **RevokeAll** (operator approvals)

  ```json
  { "approve_all": { "operator": "paxi1...", "expires": null } }
  { "revoke_all":  { "operator": "paxi1..." } }
  ```

---

### QueryMsg

* **OwnerOf**

  ```json
  { "owner_of": { "token_id": "unique_id" } }
  ```
* **Tokens** (by owner)

  ```json
  { "tokens": { "owner": "paxi1...", "start_after": null, "limit": 50 } }
  ```
* **AllTokens** (collection)

  ```json
  { "all_tokens": { "start_after": null, "limit": 50 } }
  ```
* **Approvals** / **AllApprovals**

  ```json
  { "approvals":     { "token_id": "unique_id" } }
  { "all_approvals": { "owner": "paxi1...", "start_after": null, "limit": 50 } }
  ```
* **TokenInfo** (metadata)

  ```json
  { "nft_info": { "token_id": "unique_id" } }
  ```

---

All other CW721 conventions (pagination, expiration, hooks) apply unchanged. This makes PRC‑721 a drop‑in replacement for any CW721‑based NFT application.
