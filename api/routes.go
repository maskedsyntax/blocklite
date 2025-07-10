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
}
