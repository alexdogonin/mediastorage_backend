package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

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

	mux := chi.NewMux()

	mux.HandleFunc("/media", func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		cursor := query.Get("cursor")

		dir := rootPath
		filename := ""
		if len(cursor) != 0 {
			lastFile, err := base64.RawURLEncoding.DecodeString(cursor)
			if err != nil {
				log.Fatal(err)
			}

			dir, filename = path.Split(string(lastFile))
			dir = path.Dir(dir)
		}
		_ = filename
		_ = dir

		resp := struct {
			Media []struct {
				// ID string
				ThumbURL    string `json:"thumb_url"`
				DetailURL   string `json:"detail_url"`
				OriginalURL string `json:"original_url"`
			} `json:"media"`
			Cursor string `json:"cursor"`
		}{}

		limit := 50

	mainLoop:
		for {
			dirEntry, err := fs.ReadDir(os.DirFS(dir), ".")
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			var e fs.DirEntry
			found := len(filename) == 0
			for _, e = range dirEntry {
				_ = e
				if !found {
					if e.Name() == filename {
						found = true
					}
					continue
				}

				if e.IsDir() {
					dir = path.Join(dir, e.Name())
					filename = ""
					continue mainLoop
				}

				fileFullName := path.Join(dir, e.Name())

				fileID := base64.RawURLEncoding.EncodeToString([]byte(fileFullName))
				fUrl := "http://localhost:" + port + "/media/" + fileID

				resp.Media = append(resp.Media, struct {
					ThumbURL    string "json:\"thumb_url\""
					DetailURL   string "json:\"detail_url\""
					OriginalURL string "json:\"original_url\""
				}{
					DetailURL:   fUrl,
					OriginalURL: fUrl,
					ThumbURL:    fUrl,
				})

				if len(resp.Media) == limit {
					resp.Cursor = fileID
					break mainLoop
				}
			}

			if dir == rootPath {
				break
			}

			dir, filename = path.Split(dir)
			dir = path.Dir(dir)
		}

		err := json.NewEncoder(rw).Encode(resp)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/media/{id}", func(rw http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		fname, err := base64.RawURLEncoding.DecodeString(id)
		if err != nil {
			log.Println(err)
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		http.ServeFile(rw, req, string(fname))
	})

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
