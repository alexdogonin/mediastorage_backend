package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func NewMediaItemThumb(s Servicer) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if len(id) == 0 {
			log.Println("parameter id is required")
			http.Error(rw, "parameter id is required", http.StatusBadRequest)
			return
		}

		UUID, err := uuid.Parse(id)
		if err != nil {
			log.Println(err)
			http.Error(rw, "id is not correct uuid", http.StatusBadRequest)
			return
		}

		data, err := s.ItemDetail(UUID)
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		_, err = rw.Write(data)
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}
	}
}
