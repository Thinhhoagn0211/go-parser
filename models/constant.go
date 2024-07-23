package models

import "time"

type Transaction struct {
	Hash      string    `bson:"hash"`
	To        string    `bson:"to"`
	Value     string    `bson:"value"`
	Gas       uint64    `bson:"gas"`
	GasPrice  string    `bson:"gas_price"`
	Nonce     uint64    `bson:"nonce"`
	Data      []byte    `bson:"data"`
	Timestamp time.Time `bson:"timestamp"`
}

type Document struct {
	Transactions []Transaction `bson:"transactions"`
}

type Block struct {
	BlockNumber string `json:"string"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
