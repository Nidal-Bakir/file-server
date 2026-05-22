package storage

import (
	"os"
)

func createRoot() (*os.Root, error) {
	err := os.MkdirAll("./storage_dump/", os.ModePerm)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot("./storage_dump/")
	return root, err
}
