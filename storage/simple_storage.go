package storage

import (
	"fmt"
	"io"
	"os"
)

type SimpleStoreParam struct {
	RootDir           *os.Root
	PathTransformFunc PathTransformFunc
}

type SimpleStorage struct {
	SimpleStoreParam
}

func NewSimpleStorage(params SimpleStoreParam) Storage {
	return &SimpleStorage{SimpleStoreParam: params}
}

func (s SimpleStorage) StramStore(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	if err := s.RootDir.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	fileName := fmt.Sprint(path, string(os.PathSeparator), "some_name.txt")
	f, err := s.RootDir.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("written new file to disk wiht path %s, and using %d bytes \n", fileName, n)

	return nil
}
