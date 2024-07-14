package main

import (
	"github.com/Thinhhoagn0211/go-parser/internal/api"
)

func main() {
	router := api.SetupRouter()
	router.Run(":8080")
}
