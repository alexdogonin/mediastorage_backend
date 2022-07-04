package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (s *Storage) AddItemToAlbum(albumUUID uuid.UUID, albumItem root.MediaAlbumItem) error {
	return s.s.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(albumItem)
		if err != nil {
			return err
		}

		err = txn.Set([]byte("index:albums:by_item:"+albumItem.UUID.String()), albumUUID[:])
		if err != nil {
			return err
		}

		return txn.Set([]byte("albums:"+albumUUID.String()+":items:"+albumItem.UUID.String()), data)
	})
}
