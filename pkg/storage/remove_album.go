package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

func (s *Storage) RemoveAlbum(UUID uuid.UUID) error {
	return s.s.Update(func(txn *badger.Txn) error {
		opt := badger.IteratorOptions{
			Prefix: []byte("albums:" + UUID.String() + ":items:"),
		}

		it := txn.NewIterator(opt)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			err := txn.Delete(it.Item().Key())
			if err != nil {
				return err
			}
		}

		var a struct {
			Path string `json:"path"`
		}
		item, err := txn.Get([]byte("albums:" + UUID.String()))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &a)
		})
		if err != nil {
			return err
		}

		err = txn.Delete([]byte("index:albums:by_path:" + a.Path))
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte("index:albums:by_item:" + UUID.String()))
		if err != nil {
			return err
		}

		var baseAlbumUUID uuid.UUID
		err = item.Value(func(val []byte) error {
			copy(baseAlbumUUID[:], val)

			return nil
		})
		if err != nil {
			return err
		}

		err = txn.Delete([]byte("index:albums:by_item:" + UUID.String()))
		if err != nil {
			return err
		}

		err = txn.Delete([]byte("albums:" + baseAlbumUUID.String() + ":items:" + UUID.String()))
		if err != nil {
			return err
		}

		return txn.Delete([]byte("albums:" + UUID.String()))
	})
}
