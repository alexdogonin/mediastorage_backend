package http

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMediaItem() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		fname, err := base64.RawURLEncoding.DecodeString(id)
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		http.ServeFile(rw, req, string(fname))
	}
}
