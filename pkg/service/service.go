package service

import (
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
	"github.com/rwcarlsen/goexif/exif"
)

var rootAlbumUUID = uuid.Nil

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo}
}

func (s *Service) Item(uuid uuid.UUID) (root.MediaItem, error) {
	return s.repo.Item(uuid)
}

func (s *Service) List(cursor string, limit uint) ([]root.MediaItem, string, error) {
	return s.repo.List(cursor, limit)
}

func (s *Service) Album(UUID uuid.UUID, limit uint, cursor string) (root.MediaAlbum, string, error) {
	return s.repo.Album(UUID, limit, cursor)
}

func (s *Service) ItemThumb(UUID uuid.UUID) ([]byte, error) {
	return s.repo.ItemThumb(UUID)
}

func (s *Service) ItemDetail(UUID uuid.UUID) ([]byte, error) {
	return s.repo.ItemDetail(UUID)
}
