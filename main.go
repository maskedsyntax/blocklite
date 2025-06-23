package main

import (
	"blocklite/blockchain"
	"blocklite/utils"
	"fmt"
)

func main() {
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
	bc.GetPreviousBlock()

}
