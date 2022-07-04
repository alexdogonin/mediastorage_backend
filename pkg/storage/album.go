package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (s *Storage) Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error) {
	curs := albumsCursor{
		Limit: limit,
		UUID:  UUID.String(),
	}
	if limit == 0 {
		curs.Limit = defaultLimit
	}

	if len(cursor) != 0 {
		err := curs.Parse(cursor)
		if err != nil {
			return root.MediaAlbum{}, "", err
		}
	}

	var album root.MediaAlbum //TODO create an inner type
	err := s.s.View(func(txn *badger.Txn) error {
		albItem, err := txn.Get([]byte("albums:" + curs.UUID))
		if err != nil {
			return err
		}

		err = albItem.Value(func(val []byte) error {
			return json.Unmarshal(val, &album)
		})
		if err != nil {
			return err
		}

		opt := badger.IteratorOptions{
			Prefix: []byte("albums:" + curs.UUID + ":items:"),
		}
		it := txn.NewIterator(opt)
		defer it.Close()

		it.Seek([]byte("albums:" + curs.UUID + ":items:" + curs.ItemUUID))
		if len(curs.ItemUUID) != 0 {
			it.Next()
		}

		for ; it.Valid(); it.Next() {
			var item root.MediaAlbumItem
			err = it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &item)
			})

			if err != nil {
				return err
			}

			album.Items = append(album.Items, item) //TODO preinit slice
			curs.ItemUUID = item.UUID.String()

			if len(album.Items) == int(curs.Limit) {
				break
			}
		}

		return nil
	})

	if err != nil {
		return root.MediaAlbum{}, "", err
	}

	return album, curs.String(), nil
}
