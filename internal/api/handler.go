package api

import (
	"net/http"

	"github.com/Thinhhoagn0211/go-parser/internal/blockchain"
	"github.com/gin-gonic/gin"
)

var parser *blockchain.BlockchainParser

func InitParser(p *blockchain.BlockchainParser) {
	parser = p
}

func GetCurrentBlock(c *gin.Context) {
	blockNumber, err := parser.GetCurrentBlockNumber()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block_number": blockNumber})
}

func SubscribeAddress(c *gin.Context) {
	var request struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parser.AddSubscription(request.Address)
	c.JSON(http.StatusOK, gin.H{"status": "subscribed", "address": request.Address})
}

func GetTransactions(c *gin.Context) {
	address := c.Query("address")

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	transactions := parser.GetTransactions(address)
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
