# Paxi Whitepaper

## Introduction

Paxi is a next-generation blockchain protocol designed with simplicity, speed, and radical decentralization at its core. Built on the Cosmos SDK and CometBFT consensus engine — with deep optimizations tailored for enhanced performance and security — and implemented in the Go programming language, Paxi aims to provide a highly efficient, low-fee, and user-friendly blockchain infrastructure that is open to all.

Paxi also supports CosmWasm, a smart contract platform based on Wasm, enabling developers to build secure and efficient contracts in Rust — perfect for DeFi, DAOs, and cross-chain dApps.

Paxi believes in the philosophy of "less is more" — every feature exists to serve purpose and efficiency. By focusing on essential functionality and avoiding unnecessary complexity, Paxi ensures a seamless experience for developers, validators, and end users alike.

---

## Vision

Paxi envisions building a **fast, secure, and truly decentralized** blockchain ecosystem—one that doesn't just belong to the whales or large institutions, but to **everyone**: users, developers, and participants alike. Here, anyone can take part in consensus through staking, contribute to innovation through development, and engage in decision-making through DAO governance proposals.

We believe that only through the collective power of the community can we achieve a truly free and sustainable future. Paxi is committed to promoting a transparent, fair, and open governance system where every participant has a voice, a role, and the ability to help shape the future.

---

## Core Principles

- **Simplicity**: Minimal design, clean APIs, and non-bloated protocol logic
- **Speed**: Optimized performance and high throughput, capable of handling thousands of transactions per second
- **Decentralization**: Open validator set, low participation threshold, governance by the many, not the few
- **Security**: Built on robust Byzantine Fault Tolerance and a lean, auditable codebase
- **Accessibility**: Designed for developers of all levels, including no-code and low-code contract authors

---

## Technology Stack

- **Cosmos SDK**: Modular framework for blockchain application development
- **CometBFT**: BFT consensus engine providing fast finality and secure networking
- **CosmWasm**: Smart contract platform based on WebAssembly (Wasm), allowing developers to build secure and efficient contracts in Rust 
- **Go Language**: Known for its speed, simplicity, and concurrency support — ideal for high-performance blockchain infrastructure

---

## Validator Participation

Paxi lowers both technical and economic barriers, enabling anyone to become a validator and promoting a more open and democratic consensus mechanism.

- No large stake is required — **only 1,000 PAXI is needed to become a validator**  
- Geographic and demographic diversity is encouraged  
- Built-in incentives for rewards and stability  
- Users may also choose to **delegate** their tokens to other validator nodes and earn staking rewards as delegators  

To further enhance **decentralization, fairness, and network security**, Paxi implements a custom validator selection mechanism:

- The validator set is refreshed every fixed number of blocks  
- **50% of validators are selected based on the highest voting power (top N)**  
- **The remaining 50% are chosen via weighted random sampling from the rest of the candidates, proportional to their staked amount**

This hybrid model ensures high validator quality while providing smaller nodes with fair participation opportunities, thereby increasing the network’s censorship resistance and resilience.

---

## Inflation Mechanism

Paxi's native token has an initial supply of **100 million tokens** and adopts a gradually decreasing annual inflation model to incentivize staking participants, enhance network security, and maintain long-term supply stability.

### Annual Inflation Schedule
- **Year 1**: Maximum inflation rate of **8%** of the total supply
- **Year 2**: Fixed inflation rate of **4%**
- **Year 3 and beyond**: Inflation rate capped at **2%**

### Block Reward Distribution
- **95% of newly minted tokens** (block rewards) are distributed to staking participants: Including **validators** and **delegators**
- **5% of newly minted tokens** (block rewards) 5% are allocated to the **DAO community**
- The system will **automatically burn excess tokens** if the actual supply exceeds the expected supply, helping to stabilize inflation

### Important Note: Concentrated Rewards Under Low Staking Ratio
For example, in Year 1:
- Annual inflation of 8% → **8 million tokens** minted
- If only **8% of the total supply** (10 million tokens) is staked,
- Then these **8 million reward tokens** will be **fully distributed among that 8%**

### This design ensures a powerful economic incentive:
- Encourages more users to stake  
- Increases network security and consensus participation  
- Prevents long-term overinflation and excessive value dilution

---

## Token Distribution

| Category                      | Allocation | Unlock Schedule                                         | Purpose & Description |
|-------------------------------|------------|---------------------------------------------------------|------------------------|
| **Founding Team & Advisors**  | 15%        | 3% released at launch, remaining over 4 phases (every 4 months) | Incentivizes long-term commitment and prevents short-term profit-taking. Aligns team with project success. |
| **Paxi Foundation**           | 10%        | 4% released at launch, remaining over 4 phases (every 4 months) | Funds core maintenance, legal compliance, branding, and global strategic initiatives. |
| **Paxi DAO**                  | 5%         | Fully unlocked at launch                                | Managed by community through DAO. Funds ecosystem development, tooling, and decentralized governance. |
| **Private & Strategic Investors** | 15%    | 3% released at launch, remaining over 4 phases (every 4 months) | Supports partnerships in infrastructure, tech, marketing, and exchange listings. Encourages strategic collaboration. |
| **Public Sale**               | 45%        | Fully released or unlocked in strategic phases          | Provides liquidity, enables wide user participation, and funds exchange listings, marketing, and ecosystem support. |
| **User Incentives & Promotions** | 10%     | Dynamically released via tasks, referrals, and campaigns | Fuels user acquisition, community growth, and viral engagement through incentivized activities. |

---

## Paxi DAO (Decentralized Autonomous Organization)

A DAO (Decentralized Autonomous Organization) is an organizational structure managed in a decentralized manner through smart contracts and blockchain technology. Unlike traditional organizations, DAOs have no single leader or central authority. Members use governance tokens to vote collectively on the organization's operations, fund allocations, rule-making, and future directions.

### How Paxi DAO Works
- **Smart Contracts**: All rules and processes of the DAO are encoded into smart contracts on the blockchain, ensuring transparency and immutability
- **Governance Tokens**: Members holding governance tokens obtain voting rights, typically proportional to the number of tokens held
- **Proposals and Voting**: Members can propose initiatives (such as fund allocation, parameter changes, or project support), and upon community approval through voting, the smart contracts automatically execute the agreed-upon actions

### Key Functions of Paxi DAO
- **Governance Decisions**: Community voting determines key blockchain parameters such as inflation rates, transaction fees, and staking rules
- **Fund Management**: Management of the community treasury to allocate resources effectively, such as funding project development and community promotion activities
- **Collaboration Platform**: Provides a transparent, fair, and traceable collaboration environment to foster trust and cooperation among members

### Use Cases
- **On-chain Parameter Changes**: Voting to adjust inflation rates or governance thresholds
- **Software Upgrade Proposals**: DAO-driven decisions on blockchain node software upgrades to ensure stability and feature enhancements
- **Community Fund Utilization**: Voting to allocate funds to support new decentralized applications (dApps) or ecosystem-building projects
- **Permission Management**: Granting or revoking governance permissions to specific smart contracts, ensuring decentralized and transparent authority management

Through the DAO model, Paxi achieves true decentralized autonomy, allowing the community to collectively shape the future direction of the ecosystem, thereby enhancing sustainability and community engagement.

---

## Use Cases

Paxi is designed to support a wide range of real-world applications:

- **DeFi**: Low-fee, high-speed infrastructure for decentralized finance
- **GameFi**: Scalable engine for blockchain gaming and NFT integration
- **Social & Identity**: Trustless systems for social apps, identity, and credentialing
- **Enterprise & IoT**: Minimalist design makes it suitable for enterprise-grade and embedded devices

---

## Developer Experience

Paxi is committed to lowering the barrier to entry for blockchain application development. To achieve this, it will include:

- An integrated development environment (IDE) tailored for Paxi smart contracts
- Intuitive SDKs and APIs
- No-code/low-code smart contract options
- Rich documentation and tutorials

Even non-programmers will be able to deploy decentralized applications (dApps) on Paxi.

---

## Conclusion

Paxi is not just another blockchain. It is a principled approach to solving the complexities and inefficiencies found in modern Layer 1 networks. By staying true to simplicity, performance, and decentralization, Paxi aims to redefine what a truly user-first blockchain can be.

We invite developers, validators, creators, and dreamers to join the Paxi movement.

**Build less. Achieve more. Build on Paxi.**
