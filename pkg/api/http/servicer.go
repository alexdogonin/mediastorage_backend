package http

import (
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Servicer interface {
	Item(uuid.UUID) (root.MediaItem, error)
	List(cursor string, limit uint) ([]root.MediaItem, string, error)
	// Albums(cursor string, limit uint)
	// Album(uuid.UUID)
}
