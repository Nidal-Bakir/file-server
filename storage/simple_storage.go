package storage

import (
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

func (s SimpleStorage) StramStore(key string, r io.Reader) (string, error) {
	path := s.PathTransformFunc(key)
	if err := s.RootDir.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	oldFileName := fmt.Sprint(path, string(os.PathSeparator), uuid.New().String())
	file, err := s.RootDir.OpenFile(oldFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1000)
	var byteCount = 0
	sha := sha3.New256()
	for {
		n, err := r.Read(buf)
		byteCount += n
		s := buf[:n]
		if err == nil && n != 0 {
			sha.Write(s)
			file.Write(s)
			continue
		}
		if n != 0 {
			sha.Write(s)
			file.Write(s)
			continue
		}
		if errors.Is(err, io.EOF) || n == 0 {
			break
		}
		return "", err
	}

	var out [32]byte
	sha.Sum(out[:0])
	hexFileName := hex.EncodeToString(out[:])
	dir := filepath.Dir(oldFileName)
	newFileName := filepath.Clean(
		fmt.Sprint(
			dir,
			string(os.PathSeparator),
			hexFileName,
		),
	)
	err = s.RootDir.Rename(oldFileName, newFileName)
	if err != nil {
		return "", err
	}

	fmt.Printf("written new file to disk wiht path %s, and using %d bytes \n", newFileName, byteCount)

	return newFileName, nil
}
