package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) ItemByPath(p string) (uuid.UUID, bool, error) {
	var UUID uuid.UUID
	var ok bool

	err := s.s.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("index:items:by_path:" + p))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}

		err = item.Value(func(val []byte) error {
			copy(UUID[:], val)

			return nil
		})

		if err != nil {
			return err
		}

		ok = true
		return nil
	})

	return UUID, ok, err
}
