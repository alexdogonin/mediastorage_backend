package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/dgraph-io/badger"
	apihttp "github.com/mediastorage_backend/pkg/api/http"
	"github.com/mediastorage_backend/pkg/service"
	"github.com/mediastorage_backend/pkg/storage"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

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

	b, err := badger.Open(badger.DefaultOptions("/tmp/mediaserver"))
	if err != nil {
		log.Fatal("opening badger error, ", err)
	}
	defer b.Close()

	strg := storage.NewStorage(b)

	svc := service.New(&strg)

	log.Println("cache filling")
	err = svc.Sync(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()

	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.String())

			h.ServeHTTP(w, r)
		})
	})

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	itemAddr := scheme + "://" + addr + ":" + port + "/media"
	albumAddr := itemAddr + "/albums"
	mux.Get("/media", apihttp.NewMediaList(itemAddr, &svc))
	mux.Get("/media/{id}", apihttp.NewMediaItem(&svc))
	mux.Get("/v2/media", apihttp.NewMediaListV2(itemAddr, &svc))
	mux.Get("/media/albums/{id}", apihttp.NewAlbumHandler(&svc, albumAddr, itemAddr))
	mux.Get("/media/albums", apihttp.NewAlbumHandler(&svc, albumAddr, itemAddr))

	log.Println("start server")
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
