package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	apihttp "github.com/mediastorage_backend/pkg/api/http"
	"github.com/mediastorage_backend/pkg/cache"

	"github.com/go-chi/chi/v5"
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

	fmt.Println("cur dir: " + os.Getenv("PWD"))
	fmt.Println("root path: " + rootPath)

	cache := cache.NewCache()

	log.Println("cache filling")
	err := cache.Fill(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()

	mux.HandleFunc("/media", apihttp.NewMediaList("http://localhost:"+port+"/media", cache))
	mux.HandleFunc("/media/{id}", apihttp.NewMediaItem(cache))

	log.Println("start server")
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
