package api

import (
	"blocklite/blockchain"
	"blocklite/wallet"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBlocks Retrieve all blocks from the blockchain
func GetBlocks(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.Chain)
}

// CreateBlock Add a new block to the blockchain
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

// MineBlock Mine a new block with proof of work
func MineBlock(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	newProof := blockchain.ProofOfWork(latestBlock.Proof)
	newBlock := bc.CreateBlock(newProof, latestBlock.CalculateHash())

	c.JSON(http.StatusOK, gin.H{"message": "Congratulations! You just mined a block", "block": newBlock})
}

// GetProofOfWork Calculate the proof of work for the latest block
func GetProofOfWork(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	proof := blockchain.ProofOfWork(latestBlock.Proof)
	c.JSON(http.StatusOK, gin.H{"proof": proof})
}

// GetPreviousHash Return the hash of the latest block
func GetPreviousHash(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"previousHash": latestBlock.CalculateHash()})
}

// GetBlockByIndex Retrieve a block by its index
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

// GetTimestamp Return the timestamp of the latest block
func GetTimestamp(c *gin.Context, bc *blockchain.Blockchain) {
	latestBlock := bc.GetLatestBlock()
	c.JSON(http.StatusOK, gin.H{"timestamp": latestBlock.Timestamp})
}

// GetLength Return the number of blocks in the blockchain
func GetLength(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, gin.H{"length": bc.GetLength()})
}

// GetFullChain Return the entire blockchain
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

// NewTransaction Create a new transaction
func NewTransaction(c *gin.Context, bc *blockchain.Blockchain) {
	var tx blockchain.Transaction

	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify signature if sender is not "0" (system)
	if tx.Sender != "0" {
		data := tx.Sender + tx.Receiver + strconv.FormatFloat(tx.Amount, 'f', -1, 64)
		if !wallet.Verify(tx.Sender, data, tx.Signature) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	index := bc.AddTransaction(tx.Sender, tx.Receiver, tx.Amount, tx.Signature)
	c.JSON(http.StatusCreated, gin.H{"message": "Transaction will be added to Block " + strconv.Itoa(index)})
}

// GetPendingTransactions Return the list of pending transactions
func GetPendingTransactions(c *gin.Context, bc *blockchain.Blockchain) {
	c.JSON(http.StatusOK, bc.CurrentTransactions)
}

// CreateWallet Generate a new wallet
func CreateWallet(c *gin.Context) {
	w := wallet.NewWallet()
	c.JSON(http.StatusCreated, gin.H{
		"private_key": hex.EncodeToString(w.PrivateKey.D.Bytes()),
		"public_key":  w.GetAddress(),
		"address":     w.GetAddress(),
	})
}

// RegisterNodes Add new nodes to the blockchain network
func RegisterNodes(c *gin.Context, bc *blockchain.Blockchain) {
	var input struct {
		Nodes []string `json:"nodes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for _, node := range input.Nodes {
		bc.RegisterNode(node)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New nodes have been added",
		"nodes":   bc.Nodes,
	})
}

// Consensus Resolve conflicts between nodes
func Consensus(c *gin.Context, bc *blockchain.Blockchain) {
	replaced := bc.ResolveConflicts()

	if replaced {
		c.JSON(http.StatusOK, gin.H{
			"message":   "Our chain was replaced",
			"new_chain": bc.Chain,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Our chain is authoritative",
			"chain":   bc.Chain,
		})
	}
}
