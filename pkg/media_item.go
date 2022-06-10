package root

import (
	"time"

	"github.com/google/uuid"
)

type MediaItem struct {
	UUID      uuid.UUID
	Thumb     MediaItemInfo
	Detail    MediaItemInfo
	Original  MediaItemInfo
	UpdatedAt time.Time
}

type MediaItemInfo struct {
	Path   string
	Width  uint
	Height uint
	Format string
}
