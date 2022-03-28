package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

type MediaAlbumResponse struct {
	Cursor string
	Name   string
	Items  []MediaAlbumItem
}

func NewAlbumHandler(service Servicer, albumAddr, itemAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := chi.URLParam(r, "id")
		cursor := query.Get("cursor")

		var UUID uuid.UUID
		var err error
		if len(id) != 0 {
			UUID, err = uuid.Parse(id)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		album, err := service.Album(UUID, cursor)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		resp := MediaAlbumResponse{
			Name: album.Name,
		}

		for _, a := range album.Items {
			item := MediaAlbumItem{
				Type: a.Type.String(),
			}

			switch a.Type {
			case root.AlbumItem_File:
				i, err := service.Item(a.UUID)
				if err != nil {
					log.Println(err)
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}

				item.Thumb = &MediaItemInfo{
					Width:  i.Thumb.Width,
					Height: i.Thumb.Height,
					URL:    fmt.Sprintf("%s/%s", itemAddr, a.UUID.String()),
				}
				item.Detail = &MediaItemInfo{
					Width:  i.Detail.Width,
					Height: i.Detail.Width,
					URL:    fmt.Sprintf("%s/%s", itemAddr, a.UUID.String()),
				}
				item.Original = &MediaItemInfo{
					Width:  i.Original.Width,
					Height: i.Original.Height,
					URL:    fmt.Sprintf("%s/%s", itemAddr, a.UUID.String()),
				}

			case root.AlbumItem_Album:
				item.Album = &MediaAlbumInfo{
					Name: a.Name,
					URL:  fmt.Sprintf("%s/%s", albumAddr, a.UUID.String()),
				}
			}

			resp.Items = append(resp.Items, item)
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}
