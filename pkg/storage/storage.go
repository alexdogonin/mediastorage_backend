package storage

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

const (
	defaultLimit = 50
)

type Storage struct {
	s *badger.DB
}

func NewStorage(s *badger.DB) Storage {
	return Storage{s}
}

func (s *Storage) Item(uuid uuid.UUID) (root.MediaItem, error) {
	var mediaItem root.MediaItem

	err := s.s.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("items:" + uuid.String()))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &mediaItem)
		})
	})

	return mediaItem, err
}

func (s *Storage) List(cursor string, limit uint) ([]root.MediaItem, string, error) {
	mediaItems := make([]root.MediaItem, 0, limit)

	curs := itemsCursor{
		Limit: limit,
	}
	if curs.Limit == 0 {
		curs.Limit = defaultLimit
	}

	if len(cursor) != 0 {
		err := curs.Parse(cursor)
		if err != nil {
			return nil, "", err
		}
	}

	err := s.s.View(func(txn *badger.Txn) error {
		opts := badger.IteratorOptions{
			Prefix: []byte("items:"),
		}

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte("items:" + curs.UUID)); it.Valid(); it.Next() {
			err := it.Item().Value(func(val []byte) error {
				var item root.MediaItem

				err := json.Unmarshal(val, &item)
				if err != nil {
					return err
				}

				mediaItems = append(mediaItems, item)
				return nil
			})

			if err != nil {
				return err
			}

			if len(mediaItems) >= int(limit) {
				curs.UUID = mediaItems[len(mediaItems)-1].UUID.String()
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return mediaItems, curs.String(), nil
}

func (s *Storage) Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error) {
	panic("not implemented")
}

func (s *Storage) UpsertItem(item root.MediaItem) error {
	return s.s.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		return txn.Set([]byte("items:"+item.UUID.String()), data)
	})
}

func (s *Storage) UpsertAlbum(root.MediaAlbum) error {
	panic("not implemented")
}

func (s *Storage) AddItemToAlbum(albumUUID, itemUUID uuid.UUID) error {
	panic("not implemented")
}
