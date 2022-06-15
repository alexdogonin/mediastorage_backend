package root

import (
	"strconv"

	"github.com/google/uuid"
)

type MediaAlbumItemType uint

const (
	AlbumItem_Album MediaAlbumItemType = 0
	AlbumItem_File  MediaAlbumItemType = 1
)

var albumItemTypes = map[MediaAlbumItemType]string{
	AlbumItem_Album: "album",
	AlbumItem_File:  "file",
}

func (t MediaAlbumItemType) String() string {
	if v, ok := albumItemTypes[t]; ok {
		return v
	}

	return strconv.Itoa(int(t))
}

type MediaAlbum struct {
	UUID  uuid.UUID
	Name  string
	Path  string
	Items []MediaAlbumItem
}

type MediaAlbumItem struct {
	Type MediaAlbumItemType
	UUID uuid.UUID
	Name string
}
