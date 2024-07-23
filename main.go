package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/Thinhhoagn0211/go-parser/api"
	"github.com/Thinhhoagn0211/go-parser/blockchain"
	"github.com/Thinhhoagn0211/go-parser/common"
	"github.com/Thinhhoagn0211/go-parser/database"
	"github.com/Thinhhoagn0211/go-parser/docs"
	"github.com/Thinhhoagn0211/go-parser/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Main structure holding the router and HTTP client
type Main struct {
	router *gin.Engine
	client http.Client
}

// defFoo recovers from panics and logs the error
func defFoo() {
	if r := recover(); r != nil {
		log.Println("This program is packing with value", r)
	}
}

// initServer initializes the server and sets up routes
func (m *Main) initServer() error {
	var err error
	common.LoadConfig(".")
	m.router = gin.Default()

	// Dynamically set Access-Control-Allow-Origin based on the request's Origin header
	m.router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Origin", origin) // Set the received Origin as the allowed origin
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, crossdomain")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	return err
}

// @title           Swagger BlockChain Parser API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	defer defFoo()
	m := Main{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	m.client = *client

	// Initialize server
	if m.initServer() != nil {
		return
	}

	// Get host and port from the configuration
	webserverHost := common.Config.WebserverHost
	webserverPort := common.Config.WebserverPort

	// Dynamically set the host for Swagger
	swaggerHost := fmt.Sprintf("%s:%s", webserverHost, webserverPort)
	fmt.Println(swaggerHost)
	// Override Swagger's host dynamically
	docs.SwaggerInfo.Host = swaggerHost
	var database = database.DbConfig{}

	database.InitMongo("mongodb://"+common.Config.MongoHost+":"+common.Config.MongoPort, "metadata", "blockchain")
	database.Connect()
	// Initialize storage and parser
	store := storage.NewStorage(database)
	parser, err := blockchain.NewBlockchainParser("wss://mainnet.infura.io/ws/v3/9b3a250bcd62480aad9d3f814e32317b", store)
	if err != nil {
		panic(err)
	}
	api.InitParser(parser)

	// Start monitoring in the background
	go parser.GetNewestAddress()

	// Define your API endpoints
	v1 := m.router.Group("/api/v1")
	{
		v1.GET("/current-block", api.GetCurrentBlock)
		v1.POST("/subscribe", api.SubscribeAddress)
		v1.GET("/transactions", api.GetTransactions)
	}
	m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	m.router.Run(":" + webserverPort)
}
