package storage

import "sync"

type Storage struct {
	sync.RWMutex
	subscriptions map[string]bool
	transactions  map[string][]string
}

func NewStorage() *Storage {
	return &Storage{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]string),
	}
}

func (s *Storage) AddSubscription(address string) {
	s.Lock()
	defer s.Unlock()
	s.subscriptions[address] = true
}

func (s *Storage) IsSubscribed(address string) bool {
	s.RLock()
	defer s.RUnlock()
	return s.subscriptions[address]
}

func (s *Storage) AddTransaction(address, tx string) {
	s.Lock()
	defer s.Unlock()
	s.transactions[address] = append(s.transactions[address], tx)
}

func (s *Storage) GetTransactions(address string) []string {
	s.RLock()
	defer s.RUnlock()
	return s.transactions[address]
}
