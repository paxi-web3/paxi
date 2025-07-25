## PRC-20 Token Standard (Code ID: 1)

**Source code**: [https://github.com/paxi-web3/cw-plus](https://github.com/paxi-web3/cw-plus)

### Overview

PRC-20 is a CW20-compatible token with an added “freeze account” extension. The **minter** (contract admin) can freeze and unfreeze individual accounts, preventing them from sending or receiving tokens.

---

### InstantiateMsg

```jsonc
{
  "name": "MyToken",
  "symbol": "MTK",
  "decimals": 6,
  "initial_balances": [
    { "address": "paxi1...", "amount": "1000000" }
  ],
  "mint": {
    "minter": "paxi1..."
  },
  "marketing": { /* optional metadata */ }
}
```

---

### ExecuteMsg

All CW20 messages plus two extensions:

* **Transfer / TransferFrom / Approve / IncreaseAllowance / DecreaseAllowance**
* **Mint** (minter only)
* **Burn**
* **FreezeAccount**

  ```json
  { "freeze": { "address": "<target>" } }
  ```
* **UnfreezeAccount**

  ```json
  { "unfreeze": { "address": "<target>" } }
  ```

*Only the designated `minter` may call `freeze` or `unfreeze`.*

---

### QueryMsg

* **Balance**

  ```json
  { "balance": { "address": "<addr>" } }
  ```
* **Allowance**

  ```json
  { "allowance": { "owner": "<addr1>", "spender": "<addr2>" } }
  ```
  
---

### Behavior

* When an account is **frozen**, any `Transfer`, `TransferFrom`, or `Receive` attempt involving that account will revert.
* All other CW20 semantics (allowances, decimals, metadata) remain unchanged.

This extension makes PRC-20 ideal for use cases requiring regulatory controls, token vesting, or emergency freezes.
