package cache

import (
	"errors"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func (c *Cache) Item(UUID uuid.UUID) (root.MediaItem, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	itemInd, ok := c.itemsIdx[UUID.String()]
	if !ok {
		return root.MediaItem{}, errors.New("not found")
	}

	return c.items[itemInd], nil
}
