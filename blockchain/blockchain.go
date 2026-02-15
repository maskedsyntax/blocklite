package blockchain

import (
	"blocklite/utils"
	"blocklite/wallet"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ValidChain check if a given blockchain is valid
func (bc *Blockchain) ValidChain(chain []Block) bool {
	if len(chain) == 0 {
		return false
	}

	previousBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]

		// Check that the hash of the block is correct
		if block.PreviousHash != previousBlock.CalculateHash() {
			return false
		}

		// Check that the Proof of Work is correct
		if valid, _ := VerifyProof(block.Proof, previousBlock.Proof); !valid {
			return false
		}

		// Verify all transaction signatures in the block
		for _, tx := range block.Transactions {
			if tx.Sender != "0" {
				data := tx.Sender + tx.Receiver + strconv.FormatFloat(tx.Amount, 'f', -1, 64)
				if !wallet.Verify(tx.Sender, data, tx.Signature) {
					return false
				}
			}
		}

		previousBlock = block
		currentIndex++
	}

	return true
}

var timeNow = time.Now

const BlockchainFile = "blockchain.json"
const MiningReward = 50.0

// Blockchain The entire blockchain
type Blockchain struct {
	Chain               []Block
	CurrentTransactions []Transaction
	Nodes               map[string]bool
	mux                 sync.Mutex
}

// Transaction represents a transfer of value
type Transaction struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature,omitempty"`
}

// CreateBlock adds a new block to the blockchain
func (bc *Blockchain) CreateBlock(proof int, previousHash string) Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	block := Block{
		Index: len(bc.Chain) + 1,
		// Timestamp: time.Now().String(),
		Timestamp:    timeNow().UTC().Format(time.RFC3339),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	// Reset the current list of transactions
	bc.CurrentTransactions = []Transaction{}

	bc.Chain = append(bc.Chain, block)
	
	// Persistence: save the chain after every new block
	_ = bc.Save(BlockchainFile)
	
	return block
}

// NewBlockChain creates a new blockchain and adds the genesis block
func NewBlockChain() *Blockchain {
	bc := &Blockchain{
		Chain:               []Block{},
		CurrentTransactions: []Transaction{},
		Nodes:               make(map[string]bool),
	}

	// Persistence: try to load existing chain
	if err := bc.LoadFromFile(BlockchainFile); err == nil && len(bc.Chain) > 0 {
		return bc
	}

	bc.CreateBlock(1, "0")

	return bc
}

// AddTransaction creates a new transaction to go into the next mined Block
func (bc *Blockchain) AddTransaction(sender, receiver string, amount float64, signature string) int {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Signature: signature,
	})

	return len(bc.Chain) + 1
}

// GetBalance returns the balance of a given address
func (bc *Blockchain) GetBalance(address string) float64 {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	balance := 0.0
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.Sender == address {
				balance -= tx.Amount
			}
			if tx.Receiver == address {
				balance += tx.Amount
			}
		}
	}
	return balance
}

// GetLatestBlock Return the last/previous Block in the Blockchain
func (bc *Blockchain) GetLatestBlock() Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	lastBlock := bc.Chain[len(bc.Chain)-1]
	lastBlock.Print()
	return lastBlock
}

// VerifyProof Verify the proof
func VerifyProof(proof, lastProof int) (bool, string) {
	blockData := strconv.Itoa(proof) + strconv.Itoa(lastProof)
	// fmt.Printf("Generated code: %s\n", blockData)

	hashedData := utils.SHA256(blockData)
	hashHex := hex.EncodeToString(hashedData[:])

	// fmt.Printf("Encoded Hex Code: %s\n", hashHex)
	return hashHex[:4] == "0000", hashHex
}

// ProofOfWork is a simple algorithm that identifies a new proof number
// such that the hash of the concatenation of the previous proof and the new proof
// contains 4 leading zeroes. The previous proof is provided as input.
func ProofOfWork(lastProof int) int {
	proofNumber := 0

	// Keep incrementing the proof number until a valid proof is found
	// A valid proof is one that, when hashed with the last proof, results in a hash
	// with 4 leading zeroes.
	valid, hashHex := VerifyProof(proofNumber, lastProof)
	for !valid {
		proofNumber += 1
		valid, hashHex = VerifyProof(proofNumber, lastProof)
	}

	fmt.Printf("Encoded Hex Code: %s\n", hashHex)

	return proofNumber
}

// IsChainValid Function to check if the blockchain is valid
func (bc *Blockchain) IsChainValid() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	previousBlock := bc.Chain[0]
	blockIndex := 1

	for blockIndex < len(bc.Chain) {

		block := bc.Chain[blockIndex]

		// compare the indices
		if block.Index != previousBlock.Index+1 {
			fmt.Println("Index Invalid")
			return false
		}
		// Compare the previous hash of the block with the hash of the previous block
		if block.PreviousHash != previousBlock.CalculateHash() {
			fmt.Println("Hash Invalid")
			return false
		}
		// Verify the proof of work
		if valid, _ := VerifyProof(block.Proof, previousBlock.Proof); !valid {
			fmt.Println("Proof Invalid")
			return false
		}
		// Compare the timestamps
		if block.Timestamp <= previousBlock.Timestamp {
			fmt.Println("Timestamp Invalid")
			return false
		}

		previousBlock = bc.Chain[blockIndex]
		blockIndex += 1

	}

	return true
}

// Print PrintBlocks prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Chain {
		// fmt.Printf("Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s\n", block.Index, block.Timestamp, block.Proof, block.PreviousHash)
		block.Print()
	}
}

// GetLength Return the number of blocks in the blockchain.
func (bc *Blockchain) GetLength() int {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return len(bc.Chain)
}

// GetBlockByIndex Return the block at the specified index (1-based)
func (bc *Blockchain) GetBlockByIndex(index int) (Block, bool) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	if index < 1 || index > len(bc.Chain) {
		return Block{}, false
	}
	return bc.Chain[index-1], true
}

// Save serializes the blockchain and saves it to a file
func (bc *Blockchain) Save(filename string) error {
	data, err := json.MarshalIndent(bc.Chain, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// LoadFromFile loads the blockchain from a file
func (bc *Blockchain) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &bc.Chain)
}

// RegisterNode adds a new node to the list of nodes
func (bc *Blockchain) RegisterNode(address string) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.Nodes[address] = true
}

// ResolveConflicts implements our consensus algorithm.
// It replaces our chain with the longest one in the network.
func (bc *Blockchain) ResolveConflicts() bool {
	bc.mux.Lock()
	nodes := []string{}
	for node := range bc.Nodes {
		nodes = append(nodes, node)
	}
	bc.mux.Unlock()

	var newChain []Block
	maxLength := bc.GetLength()

	for _, node := range nodes {
		resp, err := http.Get("http://" + node + "/api/full-chain")
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var result struct {
				Length int     `json:"length"`
				Chain  []Block `json:"chain"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			if result.Length > maxLength && bc.ValidChain(result.Chain) {
				maxLength = result.Length
				newChain = result.Chain
			}
		} else {
			resp.Body.Close()
		}
	}

	if newChain != nil {
		bc.mux.Lock()
		bc.Chain = newChain
		bc.mux.Unlock()
		_ = bc.Save(BlockchainFile)
		return true
	}

	return false
}
