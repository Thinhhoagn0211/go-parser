package api

import (
	"net/http"

	"github.com/Thinhhoagn0211/go-parser/blockchain"
	"github.com/gin-gonic/gin"
)

var parser *blockchain.BlockchainParser

func InitParser(p *blockchain.BlockchainParser) {
	parser = p
}

// GetCurrentBlock godoc
// @Summary      Get Current Block
// @Description  Get current block which address is pointing
// @Accept       json
// @Tags 	     blockchain
// @Produce      json
// @Success      200  {array}   models.Block
// @Failure      400  {object}  models.ErrorResponse
// @Router       /current-block [get]
func GetCurrentBlock(c *gin.Context) {
	blockNumber, err := parser.GetCurrentBlockNumber()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block_number": blockNumber})
}

// SubscribeAddress godoc
// @Summary      Subcribe address
// @Description  Subcribe address into database
// @Accept       json
// @Tags blockchain
// @Produce      json
// @Param        address body  object  true  "Address to subscribe" example({"address": "string"})
// @Success      200
// @Failure      400  {object}  models.ErrorResponse
// @Router       /subscribe [post]
func SubscribeAddress(c *gin.Context) {
	var request struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := parser.AddSubscription(request.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "subscribed", "address": request.Address})
}

// GetTransactions godoc
// @Summary      Show an transactions
// @Description  get transactions in subscribed address
// @Accept       json
// @Tags         blockchain
// @Produce      json
// @Param        address query  string true  "Address to get transactions"
// @Success      200  {array}   []models.Transaction
// @Failure      400  {object}  models.ErrorResponse
// @Router       /transactions [get]
func GetTransactions(c *gin.Context) {
	address := c.Query("address")

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	transactions, err := parser.GetTransactions(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
