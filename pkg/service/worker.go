package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"
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

func (s *Service) processItem(UUID uuid.UUID) error {
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

	imgHeight := img.Bounds().Dy()
	imgWidth := img.Bounds().Dx()
	var angle float64

	switch item.Orientation {
	case root.Orientation270, root.Orientation270Mirrored:
		angle = -math.Pi
		imgHeight, imgWidth = imgWidth, imgHeight
	case root.Orientation90, root.Orientation90Mirrored:
		angle = math.Pi / 2
		imgHeight, imgWidth = imgWidth, imgHeight
	}

	scaleFactor := float64(detailHeight) / float64(imgHeight) // в данном случае (orientation)

	c := gg.NewContext(int(float64(imgWidth)*scaleFactor), detailHeight)

	c.RotateAbout(angle, float64(c.Width())/2, float64(c.Width())/2)

	c.Scale(float64(scaleFactor), float64(scaleFactor))
	c.DrawImage(img, 0, 0)

	img = c.Image()

	buf := bytes.NewBuffer(make([]byte, 0, 200))
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 70})
	if err != nil {
		return err
	}

	err = s.repo.UpsertItemDetail(UUID, buf.Bytes())
	if err != nil {
		return err
	}

	item.Detail = &root.MediaItemInfo{
		Width:  uint(c.Width()),
		Height: uint(c.Height()),
		Format: "image/png",
	}

	c = gg.NewContext(int(float64(c.Width())*detailToThumbFactor), thumbHeight)

	c.Scale(float64(detailToThumbFactor), float64(detailToThumbFactor))
	c.DrawImage(img, 0, 0)

	buf.Reset()
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 30})
	if err != nil {
		return err
	}

	err = s.repo.UpsertItemThumb(UUID, buf.Bytes())
	if err != nil {
		return err
	}

	item.Thumb = &root.MediaItemInfo{
		Width:  uint(c.Width()),
		Height: uint(c.Height()),
		Format: "image/png",
	}

	return s.repo.UpsertItem(item)
}
