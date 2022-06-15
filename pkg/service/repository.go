package service

import (
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Repository interface {
	Item(uuid.UUID) (root.MediaItem, error)
	ItemByPath(p string) (uuid.UUID, bool, error)
	List(cursor string, limit uint) ([]root.MediaItem, string, error)

	// TODO create separated methods Album and AlbumItems
	// method Album returns description of an album, AlbumItems returns album items
	Album(UUID uuid.UUID, itemsLimit uint, cursor string) (root.MediaAlbum, string, error)
	AlbumByPath(p string) (uuid.UUID, bool, error)

	// ItemAlbum(itemUUID uuid.UUID, itemsLimit uint, cursor string) (root.MediaAlbum, string, error)

	UpsertItem(root.MediaItem) error
	// RemoveItem(uuid.UUID) error

	UpsertAlbum(root.MediaAlbum) error
	AddItemToAlbum(albumUUID uuid.UUID, albumItem root.MediaAlbumItem) error
}
