package api

import (
	"blocklite/blockchain"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the API routes with the blockchain instance.
func SetupRoutes(router *gin.Engine, bc *blockchain.Blockchain) {
	router.GET("/api/blocks", func(c *gin.Context) { GetBlocks(c, bc) })
	router.POST("/api/blocks", func(c *gin.Context) { CreateBlock(c, bc) })
	router.POST("/api/mine", func(c *gin.Context) { MineBlock(c, bc) })
	router.GET("/api/proof", func(c *gin.Context) { GetProofOfWork(c, bc) })
	router.GET("/api/previous-hash", func(c *gin.Context) { GetPreviousHash(c, bc) })
	router.GET("/api/blocks/:index", func(c *gin.Context) { GetBlockByIndex(c, bc) })
	router.GET("/api/timestamp", func(c *gin.Context) { GetTimestamp(c, bc) })
	router.GET("/api/length", func(c *gin.Context) { GetLength(c, bc) })
	router.GET("/api/full-chain", func(c *gin.Context) { GetFullChain(c, bc) })
	router.POST("/api/transactions/new", func(c *gin.Context) { NewTransaction(c, bc) })
	router.GET("/api/transactions/pending", func(c *gin.Context) { GetPendingTransactions(c, bc) })
	router.POST("/api/wallet", func(c *gin.Context) { CreateWallet(c) })
	router.POST("/api/nodes/register", func(c *gin.Context) { RegisterNodes(c, bc) })
	router.GET("/api/nodes/resolve", func(c *gin.Context) { Consensus(c, bc) })
	router.GET("/api/balance/:address", func(c *gin.Context) { GetBalance(c, bc) })
}
