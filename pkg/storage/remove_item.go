package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) RemoveItem(UUID uuid.UUID) error {
	return s.s.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("index:albums:by_item:" + UUID.String()))
		if err != nil {
			return err
		}

		var albumUUID uuid.UUID
		err = item.Value(func(val []byte) error {
			copy(albumUUID[:], val)

			return nil
		})
		if err != nil {
			return err
		}

		err = txn.Delete([]byte("index:albums:by_item:" + UUID.String()))
		if err != nil {
			return err
		}

		err = txn.Delete(s.indexItemsByPathKey(UUID.String()))
		if err != nil {
			return err
		}

		mediaItem, err := s.itemByUUIDTx(txn, UUID)
		if err != nil {
			return err
		}

		err = txn.Delete(s.indexItemsByDateKey(mediaItem.UpdatedAt))
		if err != nil {
			return err
		}

		err = txn.Delete([]byte("albums:" + albumUUID.String() + ":items:" + UUID.String()))
		if err != nil {
			return err
		}

		return txn.Delete([]byte("items:" + UUID.String()))
	})
}
