package service

import (
	"time"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Repository interface {
	Item(uuid.UUID) (root.MediaItem, error)
	List(cursor string, limit uint) ([]root.MediaItem, string, error)
	File(path string) (File, error)
	ListFiles() ([]File, error)
	Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error)
	UpsertItem(root.MediaItem) error
	RemoveItem(uuid.UUID) error
	UpsertAlbum(root.MediaAlbum) error
	AddItemToAlbum(albumUUID, itemUUID uuid.UUID) error
}

type File struct {
	Path      string
	UpdatedAt time.Time
	UUID      uuid.UUID
}
