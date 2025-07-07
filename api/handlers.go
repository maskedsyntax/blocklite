package api

import (
	"blocklite/blockchain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve all blocks from the blockchain
func GetBlocks(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.Chain)
}

// Add a new block to the blockchain
func CreateBlock(c *gin.Context, bc *blockchain.Blockchain) {
}

// Mine a new block with proof of work
func MineBlock(c *gin.Context, bc *blockchain.Blockchain) {
}

// Calculate the proof of work for the latest block
func GetProofWork(c *gin.Context, bc *blockchain.Blockchain) {
}

// Return the hash of the latest block
func GetPreviousHash(c *gin.Context, bc *blockchain.Blockchain) {
}

// Retrieve a block by its index
func GetBlockByIndex(c *gin.Context, bc *blockchain.Blockchain) {
}

// Return the timestamp of the latest block
func GetTimestamp(c *gin.Context, bc *blockchain.Blockchain) {
}

// Return the number of blocks in the blockchain
func GetLength(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.GetLength())
}

// Return the entire blockchain
func GetFullChain(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.Chain)
}
