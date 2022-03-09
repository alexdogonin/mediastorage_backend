package http

import (
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
)

func NewMediaList(rootPath string, gettingContentURL string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		dir := rootPath
		filename := ""
		if cursor := query.Get("cursor"); len(cursor) != 0 {
			lastFile, err := base64.RawURLEncoding.DecodeString(cursor)
			if err != nil {
				log.Fatal(err)
			}

			dir, filename = path.Split(string(lastFile))
			dir = path.Dir(dir)
		}

		resp := MediaListResponse{}
		limit := 50

	mainLoop:
		for {
			entries, err := fs.ReadDir(os.DirFS(dir), ".")
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			found := len(filename) == 0
			for _, e := range entries {
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
				fUrl := gettingContentURL + "/" + fileID

				resp.Media = append(resp.Media, MediaItem{
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
	}
}

type MediaListResponse struct {
	Media  []MediaItem `json:"media"`
	Cursor string      `json:"cursor"`
}

type MediaItem struct {
	// ID string
	ThumbURL    string `json:"thumb_url"`
	DetailURL   string `json:"detail_url"`
	OriginalURL string `json:"original_url"`
}
