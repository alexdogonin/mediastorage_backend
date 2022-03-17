package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func NewMediaList(addr string, s Servicer) http.HandlerFunc {
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
				ThumbURL:    fmt.Sprintf("%s/%s", addr, m.UUID.String()),
				DetailURL:   fmt.Sprintf("%s/%s", addr, m.UUID.String()),
				OriginalURL: fmt.Sprintf("%s/%s", addr, m.UUID.String()),
			})
		}

		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	}
}

func NewMediaListV2(addr string, s Servicer) http.HandlerFunc {
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
				UUID: m.UUID.String(),
				Thumb: &MediaItemInfo{
					URL:    fmt.Sprintf("%s/%s", addr, m.UUID.String()),
					Width:  m.Thumb.Width,
					Height: m.Thumb.Height,
				},
				Detail: &MediaItemInfo{
					URL:    fmt.Sprintf("%s/%s", addr, m.UUID.String()),
					Width:  m.Detail.Width,
					Height: m.Detail.Height,
				},
				Original: &MediaItemInfo{
					URL:    fmt.Sprintf("%s/%s", addr, m.UUID.String()),
					Width:  m.Original.Width,
					Height: m.Original.Height,
				},
			})
		}

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

type MediaItem struct {
	UUID        string         `json:"uuid"`
	ThumbURL    string         `json:"thumb_url,omitempty"`
	DetailURL   string         `json:"detail_url,omitempty"`
	OriginalURL string         `json:"original_url,omitempty"`
	Thumb       *MediaItemInfo `json:"thumb,omitempty"`
	Detail      *MediaItemInfo `json:"detail,omitempty"`
	Original    *MediaItemInfo `json:"original,omitempty"`
}

type MediaItemInfo struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}
