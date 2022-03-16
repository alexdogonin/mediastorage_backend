package cache

import (
	"errors"

	root "github.com/mediastorage_backend/pkg"
)

func (c *Cache) List(cursorStr string, limit uint) ([]root.MediaItem, string, error) {
	c.itemsMx.RLock()
	defer c.itemsMx.RUnlock()

	cursor := cursor{
		Limit: limit,
	}

	var itemIdx uint
	if len(cursorStr) != 0 {

		if err := cursor.Parse(cursorStr); err != nil {
			return nil, "", err
		}

		var ok bool
		itemIdx, ok = c.itemsIdx[cursor.UUID]
		if !ok {
			return nil, "", errors.New("not found")
		}
	}

	resp := make([]root.MediaItem, 0, cursor.Limit)
	for _, m := range c.items[itemIdx:] {
		resp = append(resp, m)

		if len(resp) == int(cursor.Limit) {
			break
		}
	}

	if len(resp) != 0 {
		cursor.UUID = resp[len(resp)-1].UUID.String()
	}

	return resp, cursor.String(), nil
}
