package storage

import (
	"sync"

	"github.com/Thinhhoagn0211/go-parser/database"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	sync.RWMutex
	Subscriptions map[string]bool
	transactions  map[string][]string
	DbStorage     database.DbConfig
	Address       string
}

func NewStorage(db database.DbConfig) *Storage {
	return &Storage{
		Subscriptions: make(map[string]bool),
		transactions:  make(map[string][]string),
		DbStorage:     db,
	}
}

func (s *Storage) AddSubscription(address string) error {
	s.Lock()
	defer s.Unlock()
	doc := bson.M{
		"address":      address,
		"transactions": []interface{}{},
	}
	err := s.DbStorage.Insert(doc)
	if err != nil {
		return err
	}
	s.Address = address
	return nil
}

func (s *Storage) IsSubscribed(address string) bool {
	s.RLock()
	defer s.RUnlock()
	return s.Subscriptions[address]
}

func (s *Storage) AddTransaction(address, tx string) {
	s.Lock()
	defer s.Unlock()
	s.transactions[address] = append(s.transactions[address], tx)

}

func (s *Storage) GetTransactions(address string) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	transaction, err := s.DbStorage.GetTransactionsByAddress(address)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
