package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {
	rootPath := os.Getenv("ROOT_PATH")

	if len(rootPath) == 0 {
		log.Fatal(errors.New("ROOT_PATH is required"))
	}

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

	mux := http.NewServeMux()

	// mux.HandleFunc("/media", func(rw http.ResponseWriter, req *http.Request) {
	request := struct {
		Cursor string `json:"cursor"`
	}{}

	// err := json.NewDecoder(req.Body).Decode(&request)
	// if err != nil {
	// 	log.Println(err)
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// open dir
	// iterate from last file
	// if last file is null, iterate from the begin

	dir := rootPath
	filename := ""
	if len(request.Cursor) != 0 {
		lastFile, err := base64.RawStdEncoding.DecodeString(request.Cursor)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("last file: ", string(lastFile))

		dir, filename = path.Split(string(lastFile))
		dir = path.Dir(dir)
	}
	_ = filename
	_ = dir

mainLoop:
	for {
		dirEntry, err := fs.ReadDir(os.DirFS(dir), ".")
		if err != nil {
			log.Println(err)
			// rw.WriteHeader(http.StatusInternalServerError)
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

			fmt.Println(path.Join(dir, e.Name()))
			// iterate files and add to resp
			// last file is cursor
		}

		if dir == rootPath {
			break
		}

		dir, filename = path.Split(dir)
		dir = path.Dir(dir)
	}
	// resp := []struct {
	// 	Media []struct {
	// 		ID string
	// 		ThumbURL,
	// 		DetailURL,
	// 		OriginalURL string
	// 	}
	// 	Cursor string
	// }{}

	// err = json.NewEncoder(rw).Encode(resp)
	// if err != nil {
	// 	log.Println(err)
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// }
	// })
	mux.HandleFunc("/media/{id}", func(rw http.ResponseWriter, req *http.Request) {

	})

	// err := http.ListenAndServe(":"+port, mux)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
