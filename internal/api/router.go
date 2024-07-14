package api

import (
	"github.com/Thinhhoagn0211/go-parser/internal/blockchain"
	"github.com/Thinhhoagn0211/go-parser/internal/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize storage and parser
	store := storage.NewStorage()
	parser, err := blockchain.NewBlockchainParser("https://mainnet.infura.io/v3/415ca3500623445ba1aabaded14974f2", store)
	if err != nil {
		panic(err)
	}
	InitParser(parser)

	// Start monitoring in the background
	go parser.MonitorTransactions()

	// Define your API endpoints
	router.GET("/current-block", GetCurrentBlock)
	router.POST("/subscribe", SubscribeAddress)
	router.GET("/transactions", GetTransactions)

	return router
}
