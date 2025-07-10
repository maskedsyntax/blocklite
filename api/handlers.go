package api

import (
	"blocklite/blockchain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Retrieve all blocks from the blockchain
func GetBlocks(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.Chain)
}

// Add a new block to the blockchain
func CreateBlock(c *gin.Context, bc *blockchain.Blockchain) {
	var newBlock blockchain.Block

	if err := c.ShouldBindJSON(&newBlock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	latestBlock := bc.GetLatestBlock()
	newProof := blockchain.ProofOfWork(latestBlock.Proof)
	newBlock = bc.CreateBlock(newProof, latestBlock.CalculateHash())

	c.JSON(http.StatusCreated, gin.H{"message": "Block mined successfully", "block": newBlock})
}

// Mine a new block with proof of work
func MineBlock(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	newProof := blockchain.ProofOfWork(latestBlock.Proof)
	newBlock := bc.CreateBlock(newProof, latestBlock.CalculateHash())

	c.JSON(http.StatusOK, gin.H{"message": "Congratulations! You just mined a block", "block": newBlock})
}

// Calculate the proof of work for the latest block
func GetProofOfWork(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	proof := blockchain.ProofOfWork(latestBlock.Proof)
	c.JSON(http.StatusOK, gin.H{"proof": proof})
}

// Return the hash of the latest block
func GetPreviousHash(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"previousHash": latestBlock.Proof})
}

// Retrieve a block by its index
func GetBlockByIndex(c *gin.Context, bc *blockchain.Blockchain) {
	index := c.Param("index")

	idx, err := strconv.Atoi(index)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	block, ok := bc.GetBlockByIndex(idx)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block Not Found"})
		return
	}

	c.JSON(http.StatusOK, block)
}

// Return the timestamp of the latest block
func GetTimestamp(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"timestamp": latestBlock.Timestamp})
}

// Return the number of blocks in the blockchain
func GetLength(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"length": latestBlock.Timestamp})
}

// Return the entire blockchain
func GetFullChain(c *gin.Context, bc *blockchain.Blockchain) {
	response := struct {
		Length int                `json:"length"`
		Chain  []blockchain.Block `json:"chain"`
	}{
		Length: bc.GetLength(),
		Chain:  bc.Chain,
	}

	c.JSON(http.StatusOK, response)
}
