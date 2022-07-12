package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	root "github.com/mediastorage_backend/pkg"
)

func NewMediaList(s Servicer, originalUrl, thumbUrl, detailUrl func(UUID uuid.UUID) string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		cursor := query.Get("cursor")
		limit := 50
		if l := query.Get("limit"); len(l) != 0 {
			var err error
			limit, err = strconv.Atoi(l)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}
		}

		media, cursor, err := s.List(cursor, uint(limit))
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		resp := MediaListResponse{
			Cursor: cursor,
			Media:  make([]MediaItem, 0, len(media)),
		}

		for _, m := range media {
			resp.Media = append(resp.Media, MediaItem{
				UUID:        m.UUID.String(),
				ThumbURL:    thumbUrl(m.UUID),
				DetailURL:   detailUrl(m.UUID),
				OriginalURL: originalUrl(m.UUID),
			})
		}

		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	}
}

func NewMediaListV2(s Servicer, originalUrl, thumbUrl, detailUrl func(UUID uuid.UUID) string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		cursor := query.Get("cursor")
		limit := 50
		if l := query.Get("limit"); len(l) != 0 {
			var err error
			limit, err = strconv.Atoi(l)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}
		}

		media, cursor, err := s.List(cursor, uint(limit))
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		resp := MediaListResponse{
			Cursor: cursor,
			Media:  make([]MediaItem, 0, len(media)),
		}

		for _, m := range media {
			mediaItem := MediaItem{
				UUID: m.UUID.String(),
				Original: &MediaItemInfo{
					URL:    originalUrl(m.UUID),
					Width:  m.Original.Width,
					Height: m.Original.Height,
				},
			}
			switch m.Orientation {
			case root.Orientation270, root.Orientation270Mirrored, root.Orientation90, root.Orientation90Mirrored:
				mediaItem.Original.Height, mediaItem.Original.Width = mediaItem.Original.Width, mediaItem.Original.Height
			}

			mediaItem.Thumb = mediaItem.Original
			if m.Thumb != nil {
				mediaItem.Thumb = &MediaItemInfo{
					URL:    thumbUrl(m.UUID),
					Width:  m.Thumb.Width,
					Height: m.Thumb.Height,
				}
			}

			mediaItem.Detail = mediaItem.Original
			if m.Detail != nil {
				mediaItem.Detail = &MediaItemInfo{
					URL:    detailUrl(m.UUID),
					Width:  m.Detail.Width,
					Height: m.Detail.Height,
				}
			}

			resp.Media = append(resp.Media, mediaItem)
		}

		rw.Header().Set("content-type", "application/json")
		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	}
}

type MediaListResponse struct {
	Media  []MediaItem `json:"media"`
	Cursor string      `json:"cursor"`
}
