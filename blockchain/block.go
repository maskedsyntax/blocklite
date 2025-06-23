package blockchain

import "fmt"

type Block struct {
	Index        int
	Timestamp    string
	Proof        int // TODO: for now, we keep it int
	PreviousHash string
}

// Print the details of the block
func (b *Block) Print() {
	fmt.Printf("{Index: %d, Timestamp: %s, Proof: %d, PreviousHash: %s}\n", b.Index, b.Timestamp, b.Proof, b.PreviousHash)
}
