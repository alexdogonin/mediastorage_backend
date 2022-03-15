package cache

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Cache struct {
	items    []root.MediaItem
	itemsIdx map[string]uint
	itemsMx  sync.RWMutex
}

func NewCache() Cache {
	return Cache{
		items: make([]root.MediaItem, 100),
	}
}

func (c *Cache) Fill(rootDir string) error {
	if c.items == nil {
		return errors.New("cache isn't initialized")
	}

	c.itemsMx.Lock()
	defer c.itemsMx.Unlock()

	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		uuid := uuid.New()
		if _, ok := c.itemsIdx[uuid.String()]; ok {
			return errors.New(uuid.String() + " already exists")
		}

		c.items = append(c.items, root.MediaItem{
			UUID:         uuid,
			OriginalPath: path,
			DetailPath:   path,
			ThumbPath:    path,
		})

		c.itemsIdx[uuid.String()] = uint(len(c.items))

		return nil
	})
}
