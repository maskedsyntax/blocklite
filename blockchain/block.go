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
	Transactions []Transaction
	Proof        int // TODO: for now, we keep it int
	PreviousHash string
}

// CalculateHash Calculate Hash of the Block
func (b *Block) CalculateHash() string {
	txData := ""
	for _, tx := range b.Transactions {
		txData += tx.Sender + tx.Receiver + strconv.FormatFloat(tx.Amount, 'f', -1, 64)
	}
	blockData := strconv.Itoa(b.Index) + b.Timestamp + txData + strconv.Itoa(b.Proof) + b.PreviousHash
	hashedData := utils.SHA256(blockData)
	hashHex := hex.EncodeToString(hashedData[:])
	return hashHex
}

// Print the details of the block
func (b *Block) Print() {
	fmt.Printf("{Index: %d, Timestamp: %s, Transactions: %d, Proof: %d, PreviousHash: %s}\n",
		b.Index, b.Timestamp, len(b.Transactions), b.Proof, b.PreviousHash)
}
