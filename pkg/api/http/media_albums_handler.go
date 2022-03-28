package http

import "net/http"

type MediaAlbumResponse struct {
	Cursor string
	Name   string
	Items  []MediaAlbumItem
}

func NewAlbumHandler(service Servicer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

	}
}
