package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) UpsertItemThumb(UUID uuid.UUID, data []byte) error {
	return s.s.Update(func(txn *badger.Txn) error {
		return txn.Set(s.itemThumbKey(UUID), data)
	})
}
