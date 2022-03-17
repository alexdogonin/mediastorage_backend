package root

import "github.com/google/uuid"

type MediaAlbumItemType uint

const (
	AlbumItem_Album MediaAlbumItemType = 0
	AlbumItem_File  MediaAlbumItemType = 1
)

type MediaAlbum struct {
	UUID  uuid.UUID
	Name  string
	Items []MediaAlbumItem
}

type MediaAlbumItem struct {
	Type MediaAlbumItemType
	UUID uuid.UUID
}
