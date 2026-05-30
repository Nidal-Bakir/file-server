package utils

import "os"

func OpenRoot(dirName string) (*os.Root, error) {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(dirName)
	return root, err
}
