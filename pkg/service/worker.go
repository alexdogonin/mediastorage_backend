package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
	"gopkg.in/gographics/imagick.v2/imagick"
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

	scaleFactor := float64(detailHeight) / float64(imgHeight)

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

/*
Install image magick

	wget https://github.com/ImageMagick/ImageMagick6/archive/6.9.10-11.tar.gz && \
	tar xvzf 6.9.10-11.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib
*/
func (s *Service) processItem1(UUID uuid.UUID) error {
	item, err := s.repo.Item(UUID)
	if err != nil {
		return err
	}

	f, err := os.Open(item.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = mw.ReadImageBlob(fileBytes)
	if err != nil {
		return err
	}

	imageWidth := mw.GetImageWidth()
	imageHeight := mw.GetImageHeight()

	scaleFactor := float64(detailHeight) / float64(imageHeight)

	err = mw.ResizeImage(uint(float64(imageWidth)*scaleFactor), detailHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		return err
	}

	err = s.repo.UpsertItemDetail(UUID, mw.GetImageBlob())
	if err != nil {
		return err
	}

	item.Detail = &root.MediaItemInfo{
		Width:  uint(float64(imageWidth) * scaleFactor),
		Height: uint(detailHeight),
		Format: "image/jpeg",
	}

	err = mw.ResizeImage(uint(float64(imageWidth)*scaleFactor*detailToThumbFactor), thumbHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		return err
	}

	err = s.repo.UpsertItemThumb(UUID, mw.GetImageBlob())
	if err != nil {
		return err
	}

	item.Thumb = &root.MediaItemInfo{
		Width:  uint(float64(imageWidth) * scaleFactor * detailToThumbFactor),
		Height: uint(thumbHeight),
		Format: "image/jpeg",
	}

	return s.repo.UpsertItem(item)
}
