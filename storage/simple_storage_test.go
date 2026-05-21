package storage

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStorage(t *testing.T) {
	root := utilCreateRoot()
	defer func() {
		assert.Nil(t, root.Close())
		assert.Nil(t, os.RemoveAll("./storage_dump/"))
	}()

	params := SimpleStoreParam{
		PathTransformFunc: CASPathTransform,
		RootDir:           root,
	}
	store := NewSimpleStorage(params)

	path, err := store.StramStore("hi", strings.NewReader("hi"))
	if(err!=nil){
		panic(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, len(strings.Split(path, "/")), 5)
}

func utilCreateRoot() *os.Root {
	err := os.MkdirAll("./storage_dump/", os.ModePerm)
	if err != nil {
		panic(err)
	}
	root, err := os.OpenRoot("./storage_dump/")
	if err != nil {
		panic(err)
	}
	return root
}
