package blockchain

import (
	"blocklite/utils"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// The entire blockchain
type Blockchain struct {
	Chain []Block
}

// CreateBlock adds a new block to the blockchain
func (bc *Blockchain) CreateBlock(proof int, previousHash string) Block {
	block := Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now().String(),
		Proof:        proof,
		PreviousHash: previousHash,
	}

	bc.Chain = append(bc.Chain, block)
	return block
}

// Return the last/previous Block in the Blockchain
func (bc *Blockchain) GetLatestBlock() Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	lastBlock.Print()
	return lastBlock
}

// NewBlockChain creates a new blockchain and adds the genesis block
func NewBlockChain() *Blockchain {
	bc := &Blockchain{}
	bc.CreateBlock(1, "0")

	return bc
}

// Verify the proof
func VerifyProof(proof, lastProof int) bool {
	code := strconv.Itoa(proof) + strconv.Itoa(lastProof)
	fmt.Printf("Generated code: %s\n", code)

	hashedCode := utils.SHA256(code)
	hashedCodeHex := hex.EncodeToString(hashedCode[:])

	fmt.Printf("Encoded Hex Code: %s\n", hashedCodeHex)
	return hashedCodeHex[:4] == "0000"
}

// ProofOfWork is a simple algorithm that identifies a new proof number
// such that the hash of the concatenation of the previous proof and the new proof
// contains 4 leading zeroes. The previous proof is provided as input.
func ProofOfWork(lastProof int) int {
	proofNumber := 0

	// Keep incrementing the proof number until a valid proof is found
	// A valid proof is one that, when hashed with the last proof, results in a hash
	// with 4 leading zeroes.
	for VerifyProof(proofNumber, lastProof) {
		proofNumber += 1
	}

	return proofNumber
}

// Function to check if the blockchain is valid
func IsChainValid(block Block, previousBlock Block) bool {
	// compare the indices
	if block.Index != previousBlock.Index+1 {
		return false
	}
	// Compare the previous hash of the block with the hash of the previous block
	if block.PreviousHash != previousBlock.CalculateHash() {
		return false
	}
	// Verify the proof of work
	if !VerifyProof(previousBlock.Proof, block.Proof) {
		return false
	}
	// Compare the timestamps
	if block.Timestamp <= previousBlock.Timestamp {
		return false
	}

	return true
}

// PrintBlocks prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Chain {
		// fmt.Printf("Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s\n", block.Index, block.Timestamp, block.Proof, block.PreviousHash)
		block.Print()
	}
}
