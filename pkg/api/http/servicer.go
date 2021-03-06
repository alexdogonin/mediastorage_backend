package http

import (
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Servicer interface {
	Item(uuid.UUID) (root.MediaItem, error)
	ItemThumb(uuid.UUID) ([]byte, error)
	ItemDetail(uuid.UUID) ([]byte, error)
	List(cursor string, limit uint) ([]root.MediaItem, string, error)
	Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error)
}
