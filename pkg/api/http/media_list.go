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
		}

		for _, m := range media {
			resp.Media = append(resp.Media, MediaItem{
				UUID: m.UUID.String(),
				// ThumbURL:    fmt.Sprintf("%s/%s/%s", addr, m.UUID.String(), "thumb"),
				// DetailURL:   fmt.Sprintf("%s/%s/%s", addr, m.UUID.String(), "detail"),
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

type MediaListResponse struct {
	Media  []MediaItem `json:"media"`
	Cursor string      `json:"cursor"`
}

type MediaItem struct {
	UUID        string `json:"uuid"`
	ThumbURL    string `json:"thumb_url"`
	DetailURL   string `json:"detail_url"`
	OriginalURL string `json:"original_url"`
}
