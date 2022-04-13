package service

import (
	"errors"
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo}
}

func (s *Service) Fill(rootDir string) error {
	stat, err := os.Stat(rootDir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("rootDir must be a directory")
	}

	//TODO make stack of directories entry
	albums := map[string]uuid.UUID{}

	return filepath.WalkDir(rootDir, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		d := filepath.Dir(p)
		baseAlbumUUID := albums[d]

		if e.IsDir() {
			UUID := uuid.New()

			albums[p] = UUID

			return s.addAlbum(UUID, e.Name(), baseAlbumUUID)
		}

		return s.addItem(p, baseAlbumUUID)
	})
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

func (s *Service) addAlbum(UUID uuid.UUID, name string, baseAlbum uuid.UUID) error {
	a := root.MediaAlbum{
		Name: name,
		UUID: UUID,
	}

	err := s.repo.UpsertAlbum(a)
	if err != nil {
		return err
	}

	return s.repo.AddItemToAlbum(baseAlbum, a.UUID)
}

func (s *Service) addItem(filePath string, baseAlbum uuid.UUID) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		log.Println(filePath, " parse error: ", err)
		return nil
	}

	uuid := uuid.New()

	info := root.MediaItemInfo{
		Path:   filePath,
		Width:  uint(cfg.Width),
		Height: uint(cfg.Height),
		Format: format,
	}

	item := root.MediaItem{
		UUID:     uuid,
		Original: info,
		Detail:   info,
		Thumb:    info,
	}

	err = s.repo.UpsertItem(item)
	if err != nil {
		return err
	}

	// return s.repo.AddItemToAlbum(baseAlbum, item.UUID)
	return nil
}

// TODO необходимо добавить новые файла и альбомы, удалить удалённые
// идентифицировать файл можно по пути к файлу
// 1.
// пробежаться по директории. если файл есть, то
//   проверить не изменился ли он. если изменился, то
//     обновить сохранив uuid
//   иначе
//     пропустить
// иначе
//   добавить новый файл
// 2.
// пробежаться по кэшу. аналогично пункту 1, но теперь смотрим наличие и совпадение файла на жестком диске
func (s *Service) Sync(rootDir string) error {
	if err := s.checkCachedData(); err != nil {
		return err
	}

	return s.checkDirectoryData(rootDir)
}

func (s *Service) checkDirectoryData(rootDir string) error {
	stat, err := os.Stat(rootDir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("rootDir must be a directory")
	}

	//TODO make stack of directories entry
	albums := map[string]uuid.UUID{}

	return filepath.WalkDir(rootDir, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		f, err := s.repo.File(p)
		if err != nil {
			return err
		}

		info, err := e.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Equal(f.UpdatedAt) {
			return nil
		}

		d := filepath.Dir(p)
		baseAlbumUUID := albums[d]

		if e.IsDir() {
			UUID := uuid.New()

			albums[p] = UUID

			return s.addAlbum(UUID, e.Name(), baseAlbumUUID)
		}

		return s.addItem(p, baseAlbumUUID)
	})
}

func (s *Service) checkCachedData() error {
	files, err := s.repo.ListFiles()
	if err != nil {
		return err
	}

	for _, p := range files {
		info, err := os.Stat(p.Path)
		
		if err != nil {
			if err == os.ErrNotExist {
				err = s.repo.RemoveItem(p.UUID)
				if err != nil {
					return err
				}

			}

			return err
		}

		if !info.ModTime().Equal(p.UpdatedAt) {
			// update detail and thumb
		}
	}

	return nil
}
