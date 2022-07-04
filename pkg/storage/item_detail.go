package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) ItemDetail(UUID uuid.UUID) ([]byte, error) {
	var data []byte

	err := s.s.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.itemDetailKey(UUID))
		if err != nil {
			return err
		}

		data, err = item.ValueCopy(nil)

		return err
	})

	return data, err
}
