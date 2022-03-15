package cache

import (
	"errors"

	root "github.com/mediastorage_backend/pkg"
)

func (c *Cache) List(cursor string, limit uint) ([]root.MediaItem, string, error) {
	c.itemsMx.RLock()
	defer c.itemsMx.RUnlock()

	var itemIdx uint
	if len(cursor) != 0 {
		var ok bool
		itemIdx, ok = c.itemsIdx[cursor]
		if !ok {
			return nil, "", errors.New("not found")
		}

	}

	resp := make([]root.MediaItem, 0, limit)
	for _, m := range c.items[itemIdx:] {
		resp = append(resp, m)
		cursor = m.UUID.String()

		if len(resp) == int(limit) {
			break
		}
	}

	return resp, cursor, nil
}
