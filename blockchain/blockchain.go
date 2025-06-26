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
func VerifyProof(proof, lastproof int) bool {
	code := strconv.Itoa(proof) + strconv.Itoa(lastproof)
	fmt.Printf("Generated code: %s\n", code)

	hashedCode := utils.SHA256(code)
	hashedCodeHex := hex.EncodeToString(hashedCode[:])

	fmt.Printf("Encoded Hex Code: %s\n", hashedCodeHex)
	return hashedCodeHex[:4] == "0000"
}

// PrintBlocks prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Chain {
		// fmt.Printf("Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s\n", block.Index, block.Timestamp, block.Proof, block.PreviousHash)
		block.Print()
	}
}
