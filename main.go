package main

import (
	"blocklite/api"
	"blocklite/blockchain"
	"blocklite/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	testBlockChain()
}

func testBlockChain() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize blockchain (singleton)
	bc := blockchain.NewBlockChain()

	// Set up Gin router
	router := gin.Default()
	api.SetupRoutes(router, bc)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
