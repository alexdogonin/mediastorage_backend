package cache

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func walkDir(rootPath string, dir string, f fs.WalkDirFunc) error {
	filename := ""
	if len(dir) == 0 {
		dir = rootPath
	}

mainLoop:
	for {
		dirEntry, err := fs.ReadDir(os.DirFS(dir), ".")
		if err != nil {
			return err
		}

		var e fs.DirEntry
		found := len(filename) == 0
		for _, e = range dirEntry {
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

			err = f(filepath.Join(dir, e.Name()), e, err)
			if err != nil {
				return err
			}
		}

		if dir == rootPath {
			break
		}

		dir, filename = path.Split(dir)
		dir = path.Dir(dir)
	}

	return nil
}
