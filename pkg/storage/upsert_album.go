package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	root "github.com/mediastorage_backend/pkg"
)

func (s *Storage) UpsertAlbum(album root.MediaAlbum) error {
	return s.s.Update(func(txn *badger.Txn) error {
		items := album.Items
		album.Items = nil

		data, err := json.Marshal(album)
		if err != nil {
			return err
		}

		err = txn.Set([]byte("albums:"+album.UUID.String()), data)
		if err != nil {
			return err
		}

		for _, i := range items {
			data, err = json.Marshal(i)
			if err != nil {
				return err
			}

			err = txn.Set([]byte("albums:"+album.UUID.String()+":items:"+i.UUID.String()), data)
			if err != nil {
				return err
			}
		}

		return txn.Set([]byte("index:albums:by_path:"+album.Path), album.UUID[:])
	})
}
