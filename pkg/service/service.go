package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

var rootAlbumUUID = uuid.Nil

type Service struct {
	repo Repository
}

type logger struct{}

func (logger) Error(args ...interface{}) {
	log.Println(args...)
}
func (logger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (logger) Info(args ...interface{}) {
	log.Println(args...)
}
func (logger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func New(repo Repository) Service {
	s := Service{repo}

	var log Logger = logger{}
	go func() {
		for range time.Tick(time.Second) {
			log.Info("start processing")
			ts := time.Now()
			err := s.repo.WalkAndPruneQueue(func(UUID uuid.UUID) error {
				return s.processItem(UUID)
			})
			log.Info("processing's been finished, ", time.Since(ts))

			if err != nil {
				log.Error(err)
			}
		}
	}()

	return s
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
