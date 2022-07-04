package root

import (
	"time"

	"github.com/google/uuid"
)

type MediaItem struct {
	UUID      uuid.UUID
	Path      string
	UpdatedAt time.Time
	Original  MediaItemInfo
	Thumb     *MediaItemInfo
	Detail    *MediaItemInfo
}

type MediaItemInfo struct {
	Width  uint
	Height uint
	Format string
}
