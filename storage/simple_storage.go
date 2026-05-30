package storage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

func (s SimpleStorage) StreamStore(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	if err := s.RootDir.MkdirAll(path.Dir, os.ModePerm); err != nil {
		return err
	}

	file, err := s.RootDir.OpenFile(path.FullPath(), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	byteCount, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	fmt.Printf("written new file to disk wiht path %s, and using %d bytes \n", path.FullPath(), byteCount)

	return nil
}

func (s SimpleStorage) StreamRead(key string) (io.ReadCloser, error) {
	path := s.PathTransformFunc(key)
	return s.RootDir.Open(path.FullPath())
}

func (s SimpleStorage) Delete(key string) error {
	path := s.PathTransformFunc(key)
	return s.RootDir.RemoveAll(path.FullPath())
}

func (s SimpleStorage) Has(key string) bool {
	path := s.PathTransformFunc(key)
	_, err := s.RootDir.Stat(path.FullPath())
	return !errors.Is(err, fs.ErrNotExist)
}

func (s *SimpleStorage) Clear() error {
	entries, err := os.ReadDir(s.RootDir.Name())
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(s.RootDir.Name(), entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func (s *SimpleStorage) Close() error {
	return s.RootDir.Close()
}
