package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
	apihttp "github.com/mediastorage_backend/pkg/api/http"
	"github.com/mediastorage_backend/pkg/cache"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	rootPath := os.Getenv("ROOT_PATH")

	if len(rootPath) == 0 {
		log.Fatal(errors.New("ROOT_PATH is required"))
	}

	rootPath = path.Clean(rootPath)

	port := os.Getenv("PORT")
	{
		if len(port) == 0 {
			log.Fatal(errors.New("PORT is required"))
		}

		_, err := strconv.ParseUint(port, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
	}
	addr := os.Getenv("ADDRESS")
	if len(addr) == 0 {
		addr = "0.0.0.0"
	}

	scheme := os.Getenv("SCHEME")
	if len(scheme) == 0 {
		scheme = "http"
	}

	fmt.Println("cur dir: " + os.Getenv("PWD"))
	fmt.Println("root path: " + rootPath)

	cache := cache.NewCache()

	log.Println("cache filling")
	err := cache.Fill(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()

	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)

			h.ServeHTTP(w, r)
		})
	})

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))
	mux.Get("/media", apihttp.NewMediaList(scheme+"://"+addr+":"+port+"/media", cache))
	mux.Get("/media/{id}", apihttp.NewMediaItem(cache))
	mux.Get("/v2/media", apihttp.NewMediaListV2(scheme+"://"+addr+":"+port+"/media", cache))
	mux.Get("/media/albums/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		album, err := cache.Album(UUID, cursor)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		
	})

	log.Println("start server")
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
