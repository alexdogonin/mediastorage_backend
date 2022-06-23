package service

import (
	"errors"
	"image"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

var rootAlbumUUID = uuid.Nil

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

	return s.repo.AddItemToAlbum(baseAlbum, root.MediaAlbumItem{
		Type: root.AlbumItem_Album,
		UUID: UUID,
		Name: name,
	})
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

func newItemFromFile(filePath string) (root.MediaItem, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return root.MediaItem{}, err
	}

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return root.MediaItem{}, err
	}

	stat, err := f.Stat()
	if err != nil {
		return root.MediaItem{}, err
	}

	uuid := uuid.New()

	info := root.MediaItemInfo{
		Path:   filePath,
		Width:  uint(cfg.Width),
		Height: uint(cfg.Height),
		Format: format,
	}

	item := root.MediaItem{
		UUID:      uuid,
		Original:  info,
		Detail:    info,
		Thumb:     info,
		UpdatedAt: stat.ModTime(),
	}

	return item, nil
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
	if err := s.refreshDirectoryData(rootDir); err != nil {
		return err
	}

	return s.refreshCachedData()
}

func (s *Service) refreshDirectoryData(rootDir string) error {
	stat, err := os.Stat(rootDir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("rootDir must be a directory")
	}

	_, ok, err := s.repo.AlbumByPath(rootDir)
	if err != nil {
		return err
	}

	if !ok {
		err = s.repo.UpsertAlbum(root.MediaAlbum{
			UUID: rootAlbumUUID,
			Name: path.Base(rootDir),
			Path: rootDir,
		})
		if err != nil {
			return err
		}
	}

	return filepath.WalkDir(rootDir, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if e.IsDir() {
			_, ok, err := s.repo.AlbumByPath(p)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}

			d := filepath.Dir(p)
			baseAlbum, ok, err := s.repo.AlbumByPath(d)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("album " + d + " is not found")
			}

			album := root.MediaAlbum{
				UUID: uuid.New(),
				Name: e.Name(),
				Path: p,
			}
			err = s.repo.UpsertAlbum(album)
			if err != nil {
				return err
			}

			return s.repo.AddItemToAlbum(baseAlbum, root.MediaAlbumItem{
				Type: root.AlbumItem_Album,
				UUID: album.UUID,
				Name: e.Name(),
			})
		}

		mediaItemUUID, ok, err := s.repo.ItemByPath(p)
		if err != nil {
			return err
		}

		if ok {
			item, err := s.repo.Item(mediaItemUUID)
			if err != nil {
				return err
			}

			info, err := e.Info()
			if err != nil {
				return err
			}

			if info.ModTime().Equal(item.UpdatedAt) {
				return nil
			}
		}

		d := filepath.Dir(p)

		baseAlbumUUID, ok, err := s.repo.AlbumByPath(d)
		if err != nil {
			return err
		}
		if !ok {
			baseAlbumUUID = uuid.New()
			err = s.repo.UpsertAlbum(root.MediaAlbum{
				UUID: baseAlbumUUID,
				Name: e.Name(),
				Path: d,
			})

			if err != nil {
				return err
			}
		}

		item, err := newItemFromFile(p)
		if err != nil {
			return err
		}

		err = s.repo.UpsertItem(item)
		if err != nil {
			return err
		}

		return s.repo.AddItemToAlbum(baseAlbumUUID, root.MediaAlbumItem{
			Type: root.AlbumItem_File,
			UUID: item.UUID,
		})
	})
}

func (s *Service) refreshCachedData() error {
	const itemsPerPage = 2000

	var err error
	var cursor string
	var media []root.MediaItem

	for {
		media, cursor, err = s.repo.List(cursor, itemsPerPage)
		if err != nil {
			return err
		}

		if len(media) == 0 {
			break
		}

		for _, item := range media {
			_, err := os.Stat(item.Original.Path)

			if err == nil {
				continue
			}

			if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			err = s.repo.RemoveItem(item.UUID)
			if err != nil {
				return err
			}
		}

		if len(media) < itemsPerPage {
			break
		}
	}

	return s.refreshAlbum(rootAlbumUUID)
}

func (s *Service) refreshAlbum(UUID uuid.UUID) error {
	const itemsPerPage = 2000

	var album root.MediaAlbum
	var cursor string
	var err error

	for firstIter := true; ; firstIter = false {
		album, cursor, err = s.repo.Album(UUID, itemsPerPage, cursor)
		if err != nil {
			return err
		}

		if firstIter {
			if len(album.Items) == 0 {
				return s.repo.RemoveAlbum(UUID)
			}

			_, err := os.Stat(album.Path)
			if err == os.ErrNotExist {
				return s.removeAlbumCascade(UUID)
			}
			if err != nil {
				return err
			}
		}

		for _, item := range album.Items {
			if item.Type != root.AlbumItem_Album {
				continue
			}

			err = s.refreshAlbum(item.UUID)
			if err != nil {
				return err
			}
		}

		if len(album.Items) < itemsPerPage {
			break
		}
	}

	return nil
}

func (s *Service) removeAlbumCascade(UUID uuid.UUID) error {
	const itemsPerPage = 2000

	var album root.MediaAlbum
	var cursor string
	var err error

	for {
		album, cursor, err = s.repo.Album(UUID, itemsPerPage, cursor)
		if err != nil {
			return err
		}

		for _, item := range album.Items {
			if item.Type != root.AlbumItem_Album {
				continue
			}

			err = s.removeAlbumCascade(item.UUID)
			if err != nil {
				return err
			}
		}

		if len(album.Items) < itemsPerPage {
			break
		}
	}

	return s.repo.RemoveAlbum(UUID)
}
