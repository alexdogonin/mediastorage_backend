package storage

import (
	"encoding/json"
	"errors"

	"github.com/dgraph-io/badger"
	root "github.com/mediastorage_backend/pkg"
)

func (s *Storage) UpsertItem(mediaItem root.MediaItem) error {
	return s.s.Update(func(txn *badger.Txn) error {
		var oldMediaItem *root.MediaItem
		{
			i, err := s.itemByUUIDTx(txn, mediaItem.UUID)
			if err != nil {
				if !errors.Is(err, ErrNotFound) {
					return err
				}
			} else {
				oldMediaItem = &i
			}
		}

		if oldMediaItem != nil && !oldMediaItem.UpdatedAt.Equal(mediaItem.UpdatedAt) {
			err := txn.Delete(s.indexItemsByDateKey(oldMediaItem.UpdatedAt))
			if err != nil {
				return err
			}
		}

		err := txn.Set(s.indexItemsByDateKey(mediaItem.UpdatedAt), mediaItem.UUID[:])
		if err != nil {
			return err
		}

		if oldMediaItem != nil && oldMediaItem.Path != mediaItem.Path {
			err = txn.Delete(s.indexItemsByPathKey(oldMediaItem.Path))
			if err != nil {
				return err
			}
		}

		err = txn.Set(s.indexItemsByPathKey(mediaItem.Path), mediaItem.UUID[:])
		if err != nil {
			return err
		}

		data, err := json.Marshal(mediaItem)
		if err != nil {
			return err
		}

		return txn.Set([]byte("items:"+mediaItem.UUID.String()), data)
	})
}
