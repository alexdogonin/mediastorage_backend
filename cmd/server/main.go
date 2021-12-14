package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
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

	mux.HandleFunc("/media", func(rw http.ResponseWriter, req *http.Request) {
		request := struct {
			Cursor string
		}{}

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// lastFile := base64.RawStdEncoding.DecodeString(request.Cursor)
		// open dir
		// iterate from last file
		// if last file is null, iterate from the begin

		dirEntry, err := fs.ReadDir(os.DirFS(rootPath), "")
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, e := range dirEntry {
			// iterate files and add to resp
			// last file is cursor
		}

		resp := []struct {
			Media []struct {
				ID string
				ThumbURL,
				DetailURL,
				OriginalURL string
			}
			Cursor string
		}{}

		err = json.NewEncoder(rw).Encode(resp)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/media/{id}", func(rw http.ResponseWriter, req *http.Request) {

	})

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
