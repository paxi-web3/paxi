# Contributing to Paxi

Thank you for your interest in contributing to **Paxi**! 🎉  
This document outlines the standard contribution flow: fork → branch → commit → PR → review → merge.

> **Note**  
> - By contributing, you agree that your contributions will be released under the repository’s license (MIT).  
> - Please follow our [Code of Conduct](./CODE_OF_CONDUCT.md).  
> - For quick questions, join the community on Discord: https://discord.com/invite/rA9Xzs69tx  

---

## 📦 Quick Start (Fork → Branch → PR)

**Fork** the repository and clone your fork:
   ```bash
   git clone https://github.com/<your-username>/paxi.git
   cd paxi
```
Add the upstream remote (optional but recommended):

```bash
git remote add upstream https://github.com/paxi-web3/paxi.git
git fetch upstream
git checkout main
git pull --ff-only upstream main
```
Create a feature branch:

```bash
git checkout -b feat/<short-description>
# e.g. feat/add-validator-quickstart
```
Make your changes and commit:

```bash
git add .
git commit -m "docs: add validator quickstart guide"
Push to your fork:
```
```bash
git push origin feat/<short-description>
```
Open a Pull Request from your branch against paxi-web3/paxi:main.

Clearly explain the changes and why they are needed.

Link related issues if applicable.

🧪 Build, Test, Lint
Please make sure CI passes before requesting a review.

Prerequisites
Go >= 1.22

Make (optional)

Docker 26+ (optional, for containerized testing)

Build
```bash
go build ./...
# or
make build
```
Test
```bash
go test ./...
```
Lint (if available)
```bash
golangci-lint run
# or let GitHub Actions run automatically on PRs
```
Shell scripts
If you modify .sh scripts, please check them with:

```bash
shellcheck scripts/*.sh
```
✍️ Commit & Branch Style
Branch naming:

feat/<topic>, fix/<topic>, docs/<topic>, ci/<topic>

e.g. fix/rpc-timeout, docs/add-delegator-guide

Conventional Commits (recommended):

feat(staking): add new delegation option

fix(rpc): handle nil response

docs: translate whitepaper to Indonesian

ci: add golangci-lint workflow

✅ Pull Request Checklist
Before submitting a PR, ensure:

 Changes are focused and atomic (one topic per PR)

 Build and tests pass locally

 Linting passes

 Documentation updated if needed

 No temporary or OS files included (e.g. .DS_Store)

 You have read the Code of Conduct

🐛 Reporting Bugs & Requesting Features
Use GitHub Issues. Please include:

Paxi version (paxid version) / environment

Steps to reproduce

Relevant logs/output

For security issues, please report privately via the official Paxi Discord: https://discord.com/invite/rA9Xzs69tx
(Do not post exploits publicly.)
