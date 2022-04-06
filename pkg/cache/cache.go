package cache

import (
	"errors"
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Cache struct {
	items         []root.MediaItem
	itemsIdx      map[string]uint
	albums        []root.MediaAlbum
	albumsIdx     map[string]uint
	rootAlbumUUID uuid.UUID

	mx sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items:     make([]root.MediaItem, 0, 100),
		itemsIdx:  make(map[string]uint),
		albums:    make([]root.MediaAlbum, 0, 100),
		albumsIdx: make(map[string]uint),
	}
}

func (c *Cache) Fill(rootDir string) error {
	if c.items == nil {
		return errors.New("cache isn't initialized")
	}

	stat, err := os.Stat(rootDir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("rootDir must be a directory")
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	albums := map[string]uuid.UUID{}

	return filepath.WalkDir(rootDir, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		d := filepath.Dir(p)

		if e.IsDir() {
			a := root.MediaAlbum{
				Name: e.Name(),
				UUID: uuid.New(),
			}

			c.albums = append(c.albums, a)
			c.albumsIdx[a.UUID.String()] = uint(len(c.albums) - 1)
			albums[p] = a.UUID

			if p == rootDir {
				c.rootAlbumUUID = a.UUID
			}

			baseAlbumUUID := albums[d]
			baseAlbumIdx, ok := c.albumsIdx[baseAlbumUUID.String()]
			if !ok {
				return nil
				// return errors.New("album " + curAlbUUID.String() + " doesn't exist (" + d + ")")
			}

			baseAlbum := &c.albums[baseAlbumIdx]
			baseAlbum.Items = append(baseAlbum.Items, root.MediaAlbumItem{
				Type: root.AlbumItem_Album,
				UUID: a.UUID,
				Name: a.Name,
			})

			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}

		cfg, format, err := image.DecodeConfig(f)
		if err != nil {
			// return err
			log.Println(p, " parse error: ", err)
			return nil
		}

		uuid := uuid.New()
		if _, ok := c.itemsIdx[uuid.String()]; ok {
			return errors.New(uuid.String() + " already exists")
		}

		info := root.MediaItemInfo{
			Path:   p,
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

		c.itemsIdx[uuid.String()] = uint(len(c.items)) - 1
		curAlbUUID := albums[d]

		curAlbumIdx, ok := c.albumsIdx[curAlbUUID.String()]
		if !ok {
			return errors.New("album " + curAlbUUID.String() + " doesn't exist (" + d + ")")
		}

		curAlbum := &c.albums[curAlbumIdx]
		curAlbum.Items = append(curAlbum.Items, root.MediaAlbumItem{
			Type: root.AlbumItem_File,
			UUID: uuid,
		})

		return nil
	})
}
