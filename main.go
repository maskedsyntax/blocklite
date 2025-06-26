package main

import (
	"blocklite/blockchain"
	"blocklite/utils"
	"fmt"
)

func main() {
	testBlockChain()
}

func testBlockChain() {
	// Test SHA256 function
	first_hash := utils.SHA256("Hello, World!")
	fmt.Printf("%x\n", first_hash)

	// Create a new blockchain
	bc := blockchain.NewBlockChain()

	// Add a new block
	bc.CreateBlock(2, "02343402ABC12")

	// Print the blockchain
	bc.PrintBlocks()

	// Print Previous/Last Block
	fmt.Print("Previous Block: ")
	bc.GetLatestBlock()

	// Perform Proof of Work if there are at least two blocks
	if len(bc.Chain) > 1 {
		block1 := bc.Chain[len(bc.Chain)-1]
		block2 := bc.Chain[len(bc.Chain)-2]
		bc.VerifyProof(block1.Proof, block2.Proof)
	}

	fmt.Println("weeeeeeeee!")
}
