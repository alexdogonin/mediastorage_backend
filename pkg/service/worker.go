package service

import (
	"bytes"
	"image"
	"os"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

const (
	thumbHeight  = 200
	detailHeight = 1200

	detailToThumbFactor = float64(thumbHeight) / float64(detailHeight)
)

func (s *Service) processQueue() error {
	return s.repo.WalkAndPruneQueue(func(UUID uuid.UUID) error {
		item, err := s.repo.Item(UUID)
		if err != nil {
			return err
		}

		f, err := os.Open(item.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		canvas := gg.NewContextForImage(img)

		factor := float64(detailHeight) / float64(item.Original.Height)
		canvas.Scale(factor, factor)

		buf := bytes.NewBuffer(make([]byte, 0, 200))
		err = canvas.EncodePNG(buf)
		if err != nil {
			return err
		}

		err = s.repo.UpsertItemDetail(UUID, buf.Bytes())
		if err != nil {
			return err
		}

		item.Detail = &root.MediaItemInfo{
			Width:  uint(canvas.Width()),
			Height: uint(canvas.Height()),
			Format: "image/png",
		}

		canvas.Scale(detailToThumbFactor, detailToThumbFactor)
		buf.Reset()

		err = canvas.EncodePNG(buf)
		if err != nil {
			return err
		}

		err = s.repo.UpsertItemThumb(UUID, buf.Bytes())
		if err != nil {
			return err
		}

		item.Thumb = &root.MediaItemInfo{
			Width:  uint(canvas.Width()),
			Height: uint(canvas.Height()),
			Format: "image/png",
		}

		return s.repo.UpsertItem(item)
	})
}
