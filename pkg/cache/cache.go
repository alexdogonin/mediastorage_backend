package cache

import (
	"errors"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Cache struct {
	items    []root.MediaItem
	itemsIdx map[string]uint
	itemsMx  sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items:    make([]root.MediaItem, 0, 100),
		itemsIdx: make(map[string]uint),
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

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		cfg, format, err := image.DecodeConfig(f)
		if err != nil {
			return err
		}

		uuid := uuid.New()
		if _, ok := c.itemsIdx[uuid.String()]; ok {
			return errors.New(uuid.String() + " already exists")
		}

		info := root.MediaItemInfo{
			Path:   path,
			Width:  uint(cfg.Width),
			Height: uint(cfg.Height),
			Format: format,
		}

		c.items = append(c.items, root.MediaItem{
			UUID:     uuid,
			Original: info,
			Detail:   info,
			Thumb:    info,
		})

		c.itemsIdx[uuid.String()] = uint(len(c.items))

		return nil
	})
}
