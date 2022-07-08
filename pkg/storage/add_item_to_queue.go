package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) AddItemToQueue(UUID uuid.UUID) error {
	return s.s.Update(func(txn *badger.Txn) error {
		return txn.Set(s.queuedItemKey(UUID), UUID[:])
	})
}
