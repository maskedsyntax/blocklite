<p align="center">
  <img src="assets/blocklite_logo.svg" alt="BlockLite Logo" width="400" />
</p>

<h1 align="center">BlockLite</h1>

<p align="center">
  <b> A Minimalistic Decentralized Cryptocurrency Implementation in Go </b>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Status-Complete-green" />
  <img src="https://img.shields.io/badge/Built_with-Go-blueviolet" />
  <img src="https://img.shields.io/badge/API-Gin-8bc34a" />
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" />
</p>

## Overview

**BlockLite** is an educational, from-scratch implementation of a decentralized cryptocurrency. It demonstrates the core principles of blockchain technology, including Proof of Work (PoW), cryptographic security, distributed consensus, and transaction management.

The project implements a toy cryptocurrency called **MaskedCoins**, allowing users to create wallets, mine blocks for rewards, and transfer value securely between peers.

## Core Features

- **Distributed Ledger**: A tamper-evident blockchain that stores the complete transaction history.
- **Proof of Work (PoW)**: A mining mechanism that secures the network by requiring computational effort (4 leading zeros in SHA-256 hashes).
- **Wallet System**: ECDSA-based cryptographic wallets for secure identity and transaction signing.
- **Mempool & Transactions**: A transaction pool where pending transfers wait to be included in the next mined block.
- **Consensus Algorithm**: Implements the "Longest Chain Rule" to resolve conflicts and synchronize state across multiple nodes.
- **Economic Model**:
    - **Mining Rewards**: Miners are awarded 50 MaskedCoins for every block they successfully mine.
    - **Balance Verification**: Transactions are only accepted if the sender has a sufficient balance, calculated by traversing the blockchain.
- **Persistence**: Automatic state saving and loading via a local JSON file (`blockchain.json`).
- **Thread Safety**: Fully synchronized internal state to handle concurrent API requests safely.

## Use Cases

1.  **Educational Tool**: Understand how blocks are linked, how mining works, and how digital signatures verify ownership.
2.  **Blockchain Prototype**: A base for experimenting with new consensus rules, transaction types, or networking protocols.
3.  **Local Crypto Simulation**: Run multiple instances locally to simulate a small-scale decentralized network.

## Getting Started

### Prerequisites

- Go (1.18+)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/maskedsyntax/blocklite.git
   cd blocklite
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the project:
   ```bash
   go run main.go
   ```
   *Or use `air` for live reloading if installed.*

## How to Use

### 1. Generate a Wallet
Create your cryptographic identity to start receiving coins.
```bash
curl -X POST http://localhost:8080/api/wallet
```
*Save the returned `address` and `private_key`.*

### 2. Mine Blocks for Rewards
Start mining to earn MaskedCoins. You can provide your address to receive the reward.
```bash
curl -X POST http://localhost:8080/api/mine -d '{"miner_address": "YOUR_ADDRESS"}'
```

### 3. Check Your Balance
Verify how many coins you have earned or received.
```bash
curl http://localhost:8080/api/balance/YOUR_ADDRESS
```

### 4. Send Coins
Transfer coins to another address. This requires signing the transaction (currently, the API expects the signature to be provided in the request).
```bash
curl -X POST http://localhost:8080/api/transactions/new -d '{
  "sender": "YOUR_PUBLIC_KEY",
  "receiver": "RECIPIENT_ADDRESS",
  "amount": 10.5,
  "signature": "YOUR_DIGITAL_SIGNATURE"
}'
```

### 5. Network Synchronization
If running multiple nodes, register them and resolve conflicts:
```bash
# Register a neighbor
curl -X POST http://localhost:8080/api/nodes/register -d '{"nodes": ["localhost:8081"]}'

# Sync with the longest chain in the network
curl http://localhost:8080/api/nodes/resolve
```

## API Reference

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/api/full-chain` | `GET` | Retrieve the entire blockchain |
| `/api/mine` | `POST` | Mine a new block and earn rewards |
| `/api/wallet` | `POST` | Generate a new ECDSA wallet |
| `/api/balance/:address` | `GET` | Get the balance of a specific address |
| `/api/transactions/new` | `POST` | Add a new transaction to the mempool |
| `/api/transactions/pending` | `GET` | View pending transactions |
| `/api/nodes/register` | `POST` | Register new neighbor nodes |
| `/api/nodes/resolve` | `GET` | Run the consensus algorithm |

## Testing
The project includes a comprehensive test suite for all core components:
```bash
go test ./...
```

## License
This project is licensed under the [MIT License](LICENSE).
