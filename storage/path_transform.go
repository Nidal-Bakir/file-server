package storage

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type PathTransformFunc func(string) string

func PlanePathTransform(key string) string {
	return key
}

func CASPathTransform(key string) string {
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
	
	casPath := strings.Join(pathSlice, string(os.PathSeparator))
	return casPath
}
