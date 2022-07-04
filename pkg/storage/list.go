package storage

import (
	"bytes"
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

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
		opts := badger.DefaultIteratorOptions
		opts.Prefix = indexItemsByDatePrefix
		opts.Reverse = true

		it := txn.NewIterator(opts)
		defer it.Close()

		keySuffix := curs.TS
		if len(keySuffix) == 0 {
			keySuffix = "z"
		}
		it.Seek(append(indexItemsByDatePrefix, keySuffix...))

		if len(curs.TS) != 0 {
			it.Next()
		}

		for ; it.Valid(); it.Next() {
			var item *badger.Item

			err := it.Item().Value(func(val []byte) error {
				var UUID uuid.UUID
				copy(UUID[:], val)

				var err error
				item, err = txn.Get([]byte("items:" + UUID.String()))
				return err
			})
			if err != nil {
				return err
			}

			var mediaItem root.MediaItem
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &mediaItem)
			})
			if err != nil {
				return err
			}

			mediaItems = append(mediaItems, mediaItem)

			curs.TS = string(bytes.TrimPrefix(it.Item().Key(), indexItemsByDatePrefix))

			if len(mediaItems) >= int(curs.Limit) {
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
