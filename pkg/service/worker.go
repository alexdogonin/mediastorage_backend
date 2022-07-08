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

	// imgToDetailFactor := float64(detailHeight) / float64(img.Bounds().Dy())
	// width := imgToDetailFactor * float64(img.Bounds().Dx())
	// img = resize.Resize(uint(width), detailHeight, img, resize.Lanczos3)

	// c := gg.NewContext(detailHeight, int(width))
	// c.Scale(imgToDetailFactor, imgToDetailFactor)
	// c.DrawImage(img, 0, 0)
	// c.RotateAbout(-math.Pi/2, float64(detailHeight)/2, float64(width)/2)
	// img = c.Image()

	// switch item.Orientation {
	// case root.Orientation90, root.Orientation90Mirrored:
	// 	c := gg.NewContextForImage(img)
	// 	c.Rotate(-math.Pi)
	// 	img = c.Image()
	// case root.Orientation270, root.Orientation270Mirrored:
	// 	c := gg.NewContextForImage(img)
	// 	c.Rotate(math.Pi)
	// 	img = c.Image()
	// }

	buf := bytes.NewBuffer(make([]byte, 0, 200))
	// err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 30})
	// if err != nil {
	// 	return err
	// }

	// err = s.repo.UpsertItemDetail(UUID, buf.Bytes())
	// if err != nil {
	// 	return err
	// }

	// item.Detail = &root.MediaItemInfo{
	// 	Width:  uint(img.Bounds().Dx()),
	// 	Height: uint(img.Bounds().Dy()),
	// 	Format: "image/jpeg",
	// }

	ff := float64(thumbHeight) / float64(img.Bounds().Dx()) // в данном случае (orientation)
	// fff := 1 / ff

	c := gg.NewContext(int(float64(img.Bounds().Dy())*ff), thumbHeight)
	c.RotateAbout(math.Pi/2, float64(c.Width())/2, float64(c.Width())/2)
	c.Scale(float64(ff), float64(ff))
	// c.DrawImageAnchored(img, int(45*fff), 100*6, .5, .5)
	c.DrawImage(img, 0, 0)

	img = c.Image()

	// img = resize.Resize(uint(width), thumbHeight, img, resize.Lanczos3)
	// buf.Reset()

	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	err = s.repo.UpsertItemThumb(UUID, buf.Bytes())
	if err != nil {
		return err
	}

	item.Thumb = &root.MediaItemInfo{
		Width:  uint(img.Bounds().Dx()),
		Height: uint(img.Bounds().Dy()),
		Format: "image/jpeg",
	}

	return s.repo.UpsertItem(item)
}
