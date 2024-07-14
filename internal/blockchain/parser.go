package blockchain

import (
	"context"
	"log"

	"github.com/Thinhhoagn0211/go-parser/internal/storage"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (p *BlockchainParser) MonitorTransactions() {
	headers := make(chan *types.Header)
	sub, err := p.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new blocks: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			block, err := p.client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Printf("Failed to get block: %v", err)
				continue
			}

			for idx, tx := range block.Transactions() {
				msg, err := p.client.TransactionSender(context.Background(), tx, block.Hash(), uint(idx))
				if err != nil {
					log.Printf("Failed to decode transaction sender: %v", err)
					continue
				}

				if p.IsSubscribed(msg.Hex()) {
					p.AddTransaction(msg.Hex(), tx.Hash().Hex())
				}

				if tx.To() != nil && p.IsSubscribed(tx.To().Hex()) {
					p.AddTransaction(tx.To().Hex(), tx.Hash().Hex())
				}
			}
		}
	}
}

func (p *BlockchainParser) AddSubscription(address string) {
	p.storage.AddSubscription(address)
}

func (p *BlockchainParser) IsSubscribed(address string) bool {
	return p.storage.IsSubscribed(address)
}

func (p *BlockchainParser) AddTransaction(address, tx string) {
	p.storage.AddTransaction(address, tx)
}

func (p *BlockchainParser) GetTransactions(address string) []string {
	return p.storage.GetTransactions(address)
}
