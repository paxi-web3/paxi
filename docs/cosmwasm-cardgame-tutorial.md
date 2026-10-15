# Building a Card Game on Paxi with CosmWasm

> A worked, end-to-end case study using the **Three Kingdoms (三国) card battler** that ships
> on Paxi mainnet. By the end you will understand the contract patterns and the frontend
> integration needed to ship a real, gas-friendly on-chain game on Paxi.

**Level**: Intermediate · **Time**: ~45 min read · **Chain**: Paxi mainnet (`paxi-1`)

---

## Overview

Paxi is a CosmWasm 2.x chain, so any game you build is a standard CosmWasm contract plus a
web frontend that talks to it. Two properties make Paxi a good fit for games:

- **Cheap, fast blocks.** Game actions cost a fraction of a cent, so per-turn transactions are
  viable (no need to simulate everything off-chain).
- **Session-key auth.** Instead of popping the wallet on every click, the player signs once,
  and a short-lived **session key** authorizes subsequent game transactions. The experience
  feels like a Web2 game while staying fully on-chain and non-custodial.

This tutorial walks through a production contract — the Three Kingdoms card game — and extracts
the reusable patterns:

| Pattern | What it solves |
| --- | --- |
| `register_session` | One wallet signature → many gas-free game actions via a session key |
| Merged battle API | `SanguoBattlePvp` / `SanguoBattleAi` instead of a dozen one-shot messages |
| Unified fragment economy | One `u128` balance instead of four fragment buckets |
| Tunable difficulty | Operator scales AI strength without redeploying |

The complete, running source is linked in [Resources](#resources).

---

## Prerequisites

| Tool | Version | Note |
| --- | --- | --- |
| Rust | **1.69.0** (pinned) | See the toolchain warning below — this is not optional |
| `wasm32-unknown-unknown` | target for 1.69.0 | `rustup target add wasm32-unknown-unknown` |
| `paxid` | latest | Paxi CLI (Cosmos SDK based) |
| Node.js | 18+ | Frontend wallet + session signing |

> ### ⚠️ Toolchain pin: use Rust 1.69.0
> The Paxi CosmWasm VM is **0.55**, which only accepts MVP(+mutable-globals) bytecode. From
> Rust 1.70 the prebuilt `wasm32-unknown-unknown` `std` enables `sign-ext` instructions that the
> VM rejects. Even forcing `RUSTFLAGS` off does not help — the residual instructions come from
> `std` itself. Pin the toolchain:
> ```bash
> rustup toolchain install 1.69.0
> rustup override set 1.69.0
> ```
> Also pin the CosmWasm crates so a `cargo update` in CI can't drift to an incompatible line:
> ```toml
> cosmwasm-std    = { version = "=2.0.4", default-features = false, features = ["iterator", "cosmwasm_2_0"] }
> cosmwasm-crypto = "=2.0.4"
> cosmwasm-derive = "=2.0.4"
> cw-storage-plus = "=2.0.0"
> ```

---

## Architecture

The contract is split into a **generic game engine** (`games/`, `assets/`, `session/`) and a
**game-specific domain** (`sanguo/`). The generic layer handles wallets, deposits, sessions and
accounting; the domain layer only implements game rules.

```
contract/src/
├── contract.rs        # instantiate / execute / migrate entry points
├── msg.rs             # all message types (generic + Sanguo domain)
├── session/           # session-key registration & verification
│   ├── register.rs    # register_session (address-binding check)
│   └── revoke.rs
├── assets/            # PRC-20 deposit / withdraw / swap channels
├── games/             # generic engine (roulette, dice, battle scaffolding)
└── sanguo/            # Three Kingdoms rules
    ├── exec.rs        # battle resolution, reward grants
    ├── state.rs       # FRAGMENT_BALANCE, params, card templates
    ├── battle/        # ai_gen, compute, skills, bonds
    └── query.rs       # read endpoints (fragments, stats, shop)
```

The key idea: **auth and accounting are generic; only `sanguo/` is game-specific.** You can fork
this layout for chess, a lottery, or any turn-based game without touching the session or asset
code.

---

## Smart Contract

### 1. Messages

Instantiation is intentionally minimal — the operator supplies admins, a PRC-20 whitelist, and
treasury routing:

```rust
#[cw_serde]
pub struct InstantiateMsg {
    pub admins: Vec<String>,
    pub prc20_whitelist: Vec<String>,
    pub initial_burn_rate: Option<u64>,
    pub swap_module: Option<String>,
    /// multiple treasury addresses, rotated deterministically by block height + time (not random)
    pub treasury_list: Option<Vec<String>>,
    pub treasury_share: Option<u64>,
}
```

For the battle surface, avoid one message per action. Merge them behind two variants with
enums for mode and action. This keeps the `ExecuteMsg` enum small and the frontend mapping
obvious:

```rust
#[cw_serde]
pub enum PvpMode { Single, Royale }

#[cw_serde]
pub enum PvpAction { Create, Accept, Join, Cancel, Claim, Settle }

#[cw_serde]
pub enum AiMode { Casual, Ranked }

#[cw_serde]
pub enum ExecuteMsg {
    // ... admin & generic messages ...
    SanguoBattlePvp { mode: PvpMode, action: PvpAction, /* card_ids, match_id, ... */ },
    SanguoBattleAi  { mode: AiMode, level: u8, card_ids: Vec<u32> },
    // ...
}
```

> **Serde discipline.** `#[cw_serde]` derives `rename_all = "snake_case"` on enums. The variant
> `SanguoBattlePvp` serializes to the JSON key `"sanguo_battle_pvp"`. The frontend must send
> **snake_case** keys — a PascalCase key is rejected as `unknown variant`. Centralize this
> conversion (see [Calling the contract](#3-calling-the-contract)) so call sites never get it wrong.

### 2. Session-key auth (the core pattern)

Players don't want to approve a wallet popup on every card play. The contract issues a
**session key**: an ephemeral secp256k1 keypair generated in the browser. The player signs *once*
with their main wallet to authorize it. Every following game action is signed by the session key.

Registration binds the session key to its derived address and enforces a daily spend cap:

```rust
pub fn register_session(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    session_addr: String,      // bech32 paxi1... — key into the SESSIONS map
    pubkey_hex: String,        // hex secp256k1 pubkey (65 or 33 bytes)
    daily_limit: Uint128,      // loss limit if the session key leaks
) -> Result<Response, ContractError> {
    let user = info.sender;
    let addr = deps.api.addr_validate(&session_addr)?;

    // hard cap: a caller could send u128::MAX to disable the止损线 entirely
    let cap = session_daily_limit_cap(deps.storage);
    if daily_limit > cap {
        return Err(ContractError::SessionDailyLimitTooHigh { requested: daily_limit, cap });
    }

    // 🔑 address-binding check: session_addr MUST equal addr(pubkey)
    // compressed_secp256k1 → SHA-256 → RIPEMD-160 → bech32(hrp="paxi")
    let pubkey_bytes = hex::decode(pubkey_hex.to_lowercase())?;
    if pubkey_bytes.len() != 33 && pubkey_bytes.len() != 65 {
        return Err(ContractError::InvalidSignatureFormat);
    }
    let sha_out = Sha256::digest(&pubkey_bytes);
    let rip_out = Ripemd160::digest(sha_out);
    let expected = bech32::encode("paxi", rip_out.to_vec().to_base32(), Variant::Bech32)?;
    if expected != session_addr {
        return Err(ContractError::SessionAddressMismatch { expected, got: session_addr });
    }
    // store SessionInfo { user, pubkey, nonce, daily_used, expires }
}
```

Three security properties fall out of this:

1. **No impersonation.** `session_addr` is the *derived* address of the pubkey, so a thief
   cannot register a key that points at someone else's account.
2. **Loss limit.** `daily_limit` is capped at a contract-level hard ceiling (default 1000 PAXI);
   renewal must also pass the cap, so a leaked key can't be "topped up" past the limit.
3. **Replay protection.** A `nonce` is stored per session and bumped on every action; the
   frontend signs `nonce` into the message.

### 3. Unified fragment economy

A common early mistake is storing four fragment buckets (common / rare / epic / legend) and
writing reward logic for each. Instead, keep **one** `u128` balance and convert on the way in:

```rust
// FRAG_UNIT: conversion weight per rarity
const FRAG_UNIT: [u128; 4] = [1, 5, 20, 80]; // common, rare, epic, legend

pub fn to_unified(f: &Fragments) -> u128 {
    f.common as u128 * FRAG_UNIT[0]
        + f.rare   as u128 * FRAG_UNIT[1]
        + f.epic   as u128 * FRAG_UNIT[2]
        + f.legend as u128 * FRAG_UNIT[3]
}
```

A stage reward of "10 common + 5 rare" becomes `10*1 + 5*5 = 35` unified fragments in a single
`FRAGMENT_BALANCE` write. Crafting a rare card costs `60` unified fragments — the same balance.
The player never has to "save up rare fragments separately"; one number drives everything.

> The balance query keeps the four legacy fields for backward-compatible UIs, but **all values
> live in the common bucket**. New frontends should read the unified total and stop showing the
> (always-zero) rare/epic/legend buckets.

### 4. Compliance-friendly, tunable design

- **Zero-fee gameplay.** Game actions read no `exec` fee field; rewards are granted from a
  single fragment pool. There is no staking, no betting, no voting-with-tokens in the play loop.
- **Operator-tunable AI.** A single basis-point parameter (e.g. `DIFF_BOOST_BP`) scales AI
  strength per `level`. Difficulty changes are a `set_params` execute — no redeploy.

### 5. Compile & optimize

```bash
cargo build --release --target wasm32-unknown-unknown
```

The release profile is tuned for minimum on-chain size:

```toml
[profile.release]
opt-level = "z"      # smallest, not fastest — games don't need raw speed
lto = "fat"
codegen-units = 1
panic = "abort"
strip = true
```

A 1.2 MiB artifact is normal for a contract this size; the Paxi VM accepts it.

---

## Frontend Integration

### 1. Wallet connect

Use any CosmJS-compatible wallet (Keplr, Leap, or the in-app provider). You only need the
signing client and the contract address.

### 2. Session-key sign

Generate the session keypair with `@noble/secp256k1`, register it once (see
[`register_session`](#2-session-key-auth-the-core-pattern)), then sign every action with it.

```js
import * as secp from '@noble/secp256k1'; // v2.x — sign() is ASYNC, must await

// build the signed message the contract expects
const message = Session.buildMessage({
  gameId: 'sanguo',
  action,            // e.g. "battle_pvp"
  roundId,
  amountPayout: spend,
  nonce: state.sessNonce,
  extra: opts.extra, // "single:create" / "royale:join" / "casual:battle" / ...
});

const sig = await Session.sign(message); // 🔑 await — v2 returns a Promise
if (!/^[0-9a-f]{128}$/.test(sig)) throw new Error('bad signature');
```

Two footguns the Three Kingdoms frontend hit in production — don't repeat them:

- **`@noble/secp256k1` v2 `sign` is async.** Forgetting `await` yields an empty signature and a
  confusing "contract rejected" error *after* the wallet popup.
- **Nonce drift.** If the user switches devices or a previous tx failed, the on-chain nonce
  moves. Re-sync `nonce` from chain and re-sign once before throwing.

### 3. Calling the contract

Wrap every execute in one helper that (a) converts the Rust variant name to snake_case, (b)
injects `session_addr` / `nonce` / `signature`, and (c) bumps the nonce on success:

```js
async function sanguoExec(variant, fields, opts) {
  const action = opts.action;                 // "battle_pvp" | "battle_ai" | ...
  const roundId = await sanguoRoundId(action);
  const message = Session.buildMessage({
    gameId: 'sanguo', action, roundId,
    amountPayout: opts.spend || 0, nonce: state.sessNonce, extra: opts.extra,
  });
  const sig = await Session.sign(message);
  // cw_serde enums serialize as snake_case: SanguoBattlePvp -> "sanguo_battle_pvp"
  const msgKey = variant.replace(/([a-z0-9])([A-Z])/g, '$1_$2').toLowerCase();
  const payload = { [msgKey]: {
    ...fields, round_id: roundId,
    session_addr: state.sessAddr, nonce: state.sessNonce, signature: sig,
  }};
  const hash = await execContract(payload, opts.funds || [], '', opts);
  await waitForTx(hash);
  Session.bumpNonce();
  return { hash, roundId };
}
```

The `extra` field carries the mode/action/Target tuple the contract uses for dispatch and gas
tiering, e.g. `single:create:<match_id>`, `royale:join:<royale_id>`, `casual:battle`.

### 4. Reading the unified balance

```js
const f = await queryContract({ sanguo_fragments: { owner } });
const total = Number(f.common); // V2: everything is in `common`; rare/epic/legend are 0
```

Show the player **one** number. Don't surface the legacy buckets — they are always zero and only
confuse players ("why did my rare fragments stay at 0?").

---

## Deploy & Instantiate

Standard CosmWasm lifecycle on Paxi (`paxid` mirrors `wasmd`):

```bash
# 1. upload the artifact
paxid tx wasm store game_contract.wasm \
  --from <key> --chain-id paxi-1 --gas auto --fees 5000upaxi -y

# 2. instantiate (capture the returned contract address)
paxid tx wasm instantiate <CODE_ID> \
  '{"admins":["paxi1..."],"prc20_whitelist":["paxi1..."],"treasury_list":["paxi1..."]}' \
  --from <key> --chain-id paxi-1 --label "sanguo-cardgame" -y

# 3. migrate to V2 (skip the guard only if legacy battles are already settled)
paxid tx wasm execute <CONTRACT> \
  '{"migrate":{"ignore_legacy_battles":true}}' --from <admin> -y

# 4. push parameters (zero-fee play, difficulty, fragment weights)
paxid tx wasm execute <CONTRACT> '{"set_params":{ ... }}' --from <admin> -y
```

> **Migration note.** If you are migrating player fragment balances from a four-bucket schema to
> the unified one, do it in batches (`SanguoMigrateFragments`, `limit ≤ 50` per call) to avoid
> out-of-gas. The migration entry point exists specifically for this.

---

## Lessons Learned

| # | Lesson | Fix |
| --- | --- | --- |
| 1 | Contract rejected with `unknown variant` | `#[cw_serde]` enums are `snake_case`; convert variant names centrally in the frontend |
| 2 | Empty signature, wallet popup still shown | `@noble/secp256k1` v2 `sign` is async — `await` it, and reject non-128-hex sigs before popup |
| 3 | "Rare fragments never increase" complaints | You have one unified balance; stop rendering the legacy zero buckets |
| 4 | `missing field card_ids` on cancel/claim | A field required by serde must be sent even when the logic ignores it — send `[]` |
| 5 | VM rejects the wasm | Use Rust 1.69.0; newer `std` emits `sign-ext` the VM 0.55 can't run |
| 6 | Leaked session key drains the account | Bind `session_addr` to `addr(pubkey)` + enforce a daily cap + bump `nonce` per action |

---

## Resources

- **Full running example (this tutorial's source):** the Three Kingdoms contract and web frontend —
  <https://github.com/moxiaomo856/paxi-card-game>
- **PRC-20 standard** (used for the deposit/withdraw channel): [prc_20_standard.md](./prc_20_standard.md)
- **PRC-721 standard** (used for card NFTs): [prc_721_standard.md](./prc_721_standard.md)
- **CosmWasm docs:** <https://docs.cosmwasm.com>
- **Paxi whitepaper:** [whitepaper.md](./whitepaper.md)

---

*Contributions welcome. This doc follows the [contribution guide](../CONTRIBUTING.md): fork →
branch (`docs/...`) → PR. Keep it focused and the examples runnable.*
