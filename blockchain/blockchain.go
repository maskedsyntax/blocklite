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

// NewBlockChain creates a new blockchain and adds the genesis block
func NewBlockChain() *Blockchain {
	bc := &Blockchain{}
	bc.CreateBlock(1, "0")

	return bc
}

// Return the last/previous Block in the Blockchain
func (bc *Blockchain) GetLatestBlock() Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	lastBlock.Print()
	return lastBlock
}

// Verify the proof
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

// Function to check if the blockchain is valid
func (bc *Blockchain) IsChainValid() bool {

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

// PrintBlocks prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Chain {
		// fmt.Printf("Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s\n", block.Index, block.Timestamp, block.Proof, block.PreviousHash)
		block.Print()
	}
}
