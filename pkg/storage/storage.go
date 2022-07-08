package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

const (
	defaultLimit = 50
)

var (
	ErrNotFound = errors.New("not found")
)

var (
	indexItemsByDatePrefix = []byte("index:items:by_date:")
	indexItemsByPathPrefix = []byte("index:items:by_path:")
	albumsPrefix           = []byte("albums:")
	itemsSection           = []byte(":items:")
	itemsPrefix            = []byte("items:")
	thumbSuffix            = []byte(":thumb")
	detailSuffix           = []byte(":detail")
	queuePrefix            = []byte("queue:")
)

type Storage struct {
	s *badger.DB
}

func NewStorage(s *badger.DB) Storage {
	return Storage{s}
}

func (s *Storage) setItemThumb(UUID uuid.UUID, val []byte) error {
	return s.s.Update(func(txn *badger.Txn) error {
		return txn.Set(s.itemThumbKey(UUID), val)
	})
}

func (*Storage) indexItemsByDateKey(tm time.Time) []byte {
	ts := fmt.Sprintf("%010d", tm.Unix())

	return append(indexItemsByDatePrefix, ts...)
}

func (*Storage) indexItemsByPathKey(path string) []byte {
	return append(indexItemsByPathPrefix, path...)
}

func (*Storage) itemKey(UUID uuid.UUID) []byte {
	return append(itemsPrefix, []byte(UUID.String())...)
}

func (*Storage) itemThumbKey(UUID uuid.UUID) []byte {
	key := append(itemsPrefix, []byte(UUID.String())...)
	return append(key, thumbSuffix...)
}

func (*Storage) itemDetailKey(UUID uuid.UUID) []byte {
	key := append(itemsPrefix, []byte(UUID.String())...)
	return append(key, detailSuffix...)
}

func (s *Storage) itemByUUIDTx(txn *badger.Txn, UUID uuid.UUID) (root.MediaItem, error) {
	var mediaItem root.MediaItem

	item, err := txn.Get(s.itemKey(UUID))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return mediaItem, ErrNotFound
		}
		return mediaItem, err
	}

	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &mediaItem)
	})

	return mediaItem, err
}

func (*Storage) queuedItemKey(UUID uuid.UUID) []byte {
	return append(queuePrefix, []byte(UUID.String())...)
}
