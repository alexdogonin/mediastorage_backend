package root

import (
	"time"

	"github.com/google/uuid"
)

type Orientation uint

const (
	OrientationNone        Orientation = 0
	Orientation0           Orientation = 1
	Orientation0Mirrored   Orientation = 2
	Orientation90          Orientation = 3
	Orientation90Mirrored  Orientation = 4
	Orientation180         Orientation = 5
	Orientation180Mirrored Orientation = 6
	Orientation270         Orientation = 7
	Orientation270Mirrored Orientation = 8
)

type MediaItem struct {
	UUID        uuid.UUID
	Path        string
	UpdatedAt   time.Time
	Original    MediaItemInfo
	Thumb       *MediaItemInfo
	Detail      *MediaItemInfo
	Orientation Orientation
}

type MediaItemInfo struct {
	Width  uint
	Height uint
	Format string
}
