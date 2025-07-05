package main

import (
	"blocklite/blockchain"
	"fmt"
)

func main() {
	testBlockChain()
}

func testBlockChain() {
	// Create a new blockchain
	bc := blockchain.NewBlockChain()

	// Perform Proof of Work if there are at least two blocks
	newProof := blockchain.ProofOfWork(1)
	fmt.Println("New Proof: ", newProof)

	// Print Previous/Last Block
	fmt.Print("Previous Block: ")
	latestBlock := bc.GetLatestBlock()

	// Add a new block
	bc.CreateBlock(newProof, latestBlock.CalculateHash())

	// Print the blockchain
	bc.Print()

	// Verify the Proof of Work
	if len(bc.Chain) > 1 {
		if bc.IsChainValid() {
			fmt.Println("Blockchain Valid!")
		} else {
			fmt.Println("Blockchain Invalid!")
		}
	}

	fmt.Println("weeeeeeeee!")
}
