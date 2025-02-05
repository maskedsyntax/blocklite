# Block Lite
<!-- A basic crypto coin in Python3 using hashlib (sha-256 encryption). -->


## Simple Blockchain Implementation

This repository contains a simple implementation of a blockchain in Python. The code demonstrates the basic concepts of a blockchain, including block creation, mining, proof of work, and data validation.

## Features

- **Block Creation**: Each block contains an index, proof number, previous hash, data, and a timestamp.
- **Proof of Work**: A simple proof-of-work algorithm is used to mine new blocks.
- **Data Validation**: Ensures that each block is valid by checking the index, hash, proof, and timestamp.
- **Mining Rewards**: Reward system for miners for creating new blocks.
- **Node Management**: Ability to add nodes to the network.

## Getting Started

### Prerequisites

- Python 3.x

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/simple-blockchain.git
   cd simple-blockchain
   ```

2. Run the blockchain script:
   ```bash
   python blockchain.py
   ```

### Usage

The script demonstrates the creation of a simple blockchain, mining of new blocks, and adding data to the blockchain.

1. **Initialize the Blockchain**:
   ```python
   blockchain = BlockChain()
   ```

2. **Mine a New Block**:
   ```python
   print("***Mining Coinxyz about to start***")
   print(blockchain.chain)

   last_block = blockchain.latest_block
   last_proof_no = last_block.proof_no
   proof_no = blockchain.proof_of_work(last_proof_no)

   blockchain.new_data(
       sender="0",
       recipient="Jane Doe",
       quantity=1,
   )

   last_hash = last_block.calculate_hash
   block = blockchain.construct_block(proof_no, last_hash)

   print("***Mining Coinxyz has been successful***")
   print(blockchain.chain)
   ```

3. **View the Blockchain**:
   ```python
   for block in blockchain.chain:
       print(block)
   ```

## Code Overview

### Block Class

The `Block` class represents a single block in the blockchain.

- **Attributes**:
  - `index`: The position of the block in the blockchain.
  - `proof_no`: The proof number of the block.
  - `prev_hash`: The hash of the previous block.
  - `data`: The data stored in the block.
  - `timestamp`: The time the block was created.

- **Methods**:
  - `calculate_hash`: Calculates the SHA-256 hash of the block.
  - `__repr__`: Returns a string representation of the block.

### BlockChain Class

The `BlockChain` class manages the entire blockchain.

- **Attributes**:
  - `chain`: A list of all the blocks in the blockchain.
  - `current_data`: A list of data to be added to the next block.
  - `nodes`: A set of nodes in the network.

- **Methods**:
  - `construct_genesis`: Creates the genesis block.
  - `construct_block`: Adds a new block to the blockchain.
  - `check_validity`: Validates a block against the previous block.
  - `new_data`: Adds new data to the current data list.
  - `proof_of_work`: Performs the proof-of-work algorithm.
  - `verifying_proof`: Verifies the proof of a block.
  - `latest_block`: Returns the latest block in the blockchain.
  - `block_mining`: Mines a new block and rewards the miner.
  - `create_node`: Adds a new node to the network.
  - `obtain_block_object`: Creates a block object from block data.
