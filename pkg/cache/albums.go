package cache

import (
	"errors"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (c *Cache) Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error) {
	var album root.MediaAlbum
	var offset uint
	var curs albumsCursor

	if len(cursor) != 0 {
		err := curs.Parse(cursor)
		if err != nil {
			return album, "", err
		}

		UUID, err = uuid.Parse(curs.AlbumUUID)
		if err != nil {
			return album, "", err
		}

		limit = curs.Limit
		offset = curs.Offset
	}

	if UUID == (uuid.UUID{}) {
		UUID = c.rootAlbumUUID
	}

	c.mx.RLock()
	defer c.mx.RUnlock()

	aInd, ok := c.albumsIdx[UUID.String()]
	if !ok {
		return album, "", errors.New("album is not found")
	}

	album = c.albums[aInd]

	rBound := int(offset + limit)
	if rBound > len(album.Items) {
		rBound = len(album.Items)
	}

	album.Items = album.Items[offset:rBound]

	curs.Limit = limit
	curs.AlbumUUID = UUID.String()
	curs.Offset = uint(rBound)

	return album, curs.String(), nil
}
