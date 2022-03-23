package cache

import (
	"errors"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (c *Cache) Album(UUID uuid.UUID, cursor string) (root.MediaAlbum, error) {
	var album root.MediaAlbum

	if UUID == (uuid.UUID{}) {
		UUID = c.rootAlbumUUID
	}

	c.mx.RLock()
	defer c.mx.RUnlock()

	aInd, ok := c.albumsIdx[UUID.String()]
	if !ok {
		return album, errors.New("album is not found")
	}

	album = c.albums[aInd]

	return album, nil
}
