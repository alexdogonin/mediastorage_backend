package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (s *Storage) Item(UUID uuid.UUID) (root.MediaItem, error) {
	var mediaItem root.MediaItem

	err := s.s.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.itemKey(UUID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &mediaItem)
		})
	})

	return mediaItem, err
}