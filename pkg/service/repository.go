package service

import (
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Repository interface {
	Item(uuid.UUID) (root.MediaItem, error)
	List(cursor string, limit uint) ([]root.MediaItem, string, error)
	Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error)
	UpsertItem(root.MediaItem) error
	UpsertAlbum(root.MediaAlbum) error
	AddItemToAlbum(albumUUID, itemUUID uuid.UUID) error
}
