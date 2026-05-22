package storage

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ContentPath struct {
	Dir      string
	FileName string
}

func (c ContentPath) FullPath() string {
	return filepath.Join(c.Dir, c.FileName)
}

type PathTransformFunc func(string) ContentPath

func PlanePathTransform(key string) ContentPath {
	return ContentPath{Dir: key, FileName: key}
}

func CASPathTransform(key string) ContentPath {
	sha := sha3.Sum256([]byte(key))
	hexStr := hex.EncodeToString(sha[:])

	levels := 4
	l := len(hexStr)
	bucketSize := l / levels
	remanningItems := l % levels
	pathSlice := make([]string, levels)

	for i := range levels {
		from, to := i*bucketSize, (i*bucketSize)+bucketSize
		pathSlice[i] = hexStr[from:to]
	}
	if remanningItems != 0 {
		pathSlice[levels-1] = fmt.Sprint(pathSlice[levels-1], hexStr[l-remanningItems:l])
	}

	return ContentPath{
		Dir:      strings.Join(pathSlice, string(os.PathSeparator)),
		FileName: hexStr,
	}
}
