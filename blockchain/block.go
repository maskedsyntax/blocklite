package blockchain

import (
	"blocklite/utils"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Block struct {
	Index        int
	Timestamp    string
	Proof        int // TODO: for now, we keep it int
	PreviousHash string
}

// Calculate Hash of the Block
func (b *Block) CalculateHash() string {
	blockData := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.Proof) + b.PreviousHash
	hashedData := utils.SHA256(blockData)
	hashHex := hex.EncodeToString(hashedData[:])
	return hashHex
}

// Print the details of the block
func (b *Block) Print() {
	fmt.Printf("{Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s}\n",
		b.Index, b.Timestamp, b.Proof, b.PreviousHash)
}
