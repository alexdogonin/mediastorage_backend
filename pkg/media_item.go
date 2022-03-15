package root

import "github.com/google/uuid"

type MediaItem struct {
	UUID         uuid.UUID
	ThumbPath    string
	DetailPath   string
	OriginalPath string
}
