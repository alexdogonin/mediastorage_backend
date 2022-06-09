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

// type Repository interface {
// 	Item(uuid.UUID) (root.MediaItem, error)
// 	List(cursor string, limit uint) ([]root.MediaItem, string, error)
// 	Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error)

// 	UpsertItem(root.MediaItem) error
// 	RemoveItem(uuid.UUID) error

// 	UpsertAlbum(root.MediaAlbum) error
// 	AddItemToAlbum(albumUUID, itemUUID uuid.UUID) error
// }

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
			var item root.MediaItem

			err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &item)
			})
			if err != nil {
				return err
			}

			mediaItems = append(mediaItems, item)

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

		var item root.MediaAlbumItem
		for it.Seek([]byte("albums:" + curs.UUID + ":items:")); it.Valid(); it.Next() {
			err = it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &item)
			})

			if err != nil {
				return err
			}

			album.Items = append(album.Items, item)

			if len(album.Items) == int(curs.Limit) {
				curs.ItemUUID = item.UUID.String()
			}
		}

		return nil
	})

	if err != nil {
		return root.MediaAlbum{}, "", err
	}

	return album, curs.String(), nil
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
