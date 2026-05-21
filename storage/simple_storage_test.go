package storage

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStorage(t *testing.T) {
	root := utilCreateRoot()

	params := SimpleStoreParam{
		PathTransformFunc: CASPathTransform,
		RootDir:           root,
	}
	store := NewSimpleStorage(params)

	err := store.StramStore("hi", strings.NewReader("hi"))
	assert.Nil(t, err)
}

func utilCreateRoot() *os.Root {
	err := os.MkdirAll("./storage/", os.ModePerm)
	if err != nil {
		panic(err)
	}
	root, err := os.OpenRoot("./storage/")
	if err != nil {
		panic(err)
	}
	return root
}
