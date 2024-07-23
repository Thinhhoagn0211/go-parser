package blockchain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Thinhhoagn0211/go-parser/models"
	"github.com/Thinhhoagn0211/go-parser/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/mgo.v2/bson"
)

type BlockchainParser struct {
	client  *ethclient.Client
	storage *storage.Storage
}

func NewBlockchainParser(rpcURL string, storage *storage.Storage) (*BlockchainParser, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	return &BlockchainParser{client: client, storage: storage}, nil
}

func (p *BlockchainParser) GetCurrentBlockNumber() (uint64, error) {
	header, err := p.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	return header.Number.Uint64(), nil
}

func (p *BlockchainParser) GetNewestAddress() {
	var mu sync.Mutex

	addrs, err := p.storage.DbStorage.GetAllAddresses()
	if err != nil {
		fmt.Println(err)
	}
	for _, addr := range addrs {
		go p.MonitorTransactions(addr)
	}

	for {
		currentAddrs, err := p.storage.DbStorage.GetAllAddresses()
		if err != nil {
			log.Fatalf("Failed to get current addresses: %v", err)
		}
		mu.Lock()
		newAddrs := slicesDifference(currentAddrs, addrs)

		for _, addr := range newAddrs {
			fmt.Println("add new addr", addr)
			go p.MonitorTransactions(addr)
			addrs = append(addrs, addr)
		}
		mu.Unlock()
	}
}

func slicesDifference(a, b []string) []string {
	// Create a map of elements in b
	bMap := make(map[string]bool)
	for _, addr := range b {
		bMap[addr] = true
	}

	// Find elements in a that are not in b
	var diff []string
	for _, addr := range a {
		if !bMap[addr] {
			diff = append(diff, addr)
		}
	}
	return diff
}

func (p *BlockchainParser) MonitorTransactions(addr string) {
	address := common.HexToAddress(addr)
	// Subscribe to the address' logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	logs := make(chan types.Log)
	sub, err := p.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	// Handle incoming logs
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			// Extract transaction details
			txHash := vLog.TxHash
			tx, isPending, err := p.client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				log.Printf("Failed to get transaction details: %v", err)
				continue
			}
			if isPending {
				fmt.Printf("Transaction %s is still pending\n", txHash.Hex())
				continue
			}

			txData := models.Transaction{
				Hash:      txHash.Hex(),
				To:        tx.To().Hex(),
				Value:     tx.Value().String(),
				Gas:       tx.Gas(),
				GasPrice:  tx.GasPrice().String(),
				Nonce:     tx.Nonce(),
				Data:      tx.Data(),
				Timestamp: time.Now(),
			}
			// Update the document in MongoDB
			filter := bson.M{"address": addr}
			update := bson.M{
				"$push": bson.M{"transactions": txData},
			}
			err = p.storage.DbStorage.Update(filter, update)
			if err != nil {
				log.Printf("Failed to update document in MongoDB: %v", err)
			} else {
				fmt.Printf("Transaction %s added to address %s in MongoDB\n", txHash.Hex(), addr)
			}
		}
	}
}
func (p *BlockchainParser) AddSubscription(address string) error {
	err := p.storage.AddSubscription(address)
	if err != nil {
		return err
	}
	return nil
}

func (p *BlockchainParser) IsSubscribed(address string) bool {
	return p.storage.IsSubscribed(address)
}

func (p *BlockchainParser) AddTransaction(address, tx string) {
	p.storage.AddTransaction(address, tx)
}

func (p *BlockchainParser) GetTransactions(address string) (interface{}, error) {
	return p.storage.GetTransactions(address)
}
