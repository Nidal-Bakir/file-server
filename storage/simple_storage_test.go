package storage

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStramStore(t *testing.T) {
	key := "some_key_dummy"
	store, closer := utilGetStorage(t)
	defer closer()
	err := store.StramStore(key, strings.NewReader("hi"))
	assert.Nil(t, err)
}

func TestStramRead(t *testing.T) {
	key := "some_key_dummy"
	expectedFileBytes := []byte("hi")
	store, closer := utilGetStorage(t)
	defer closer()
	err := store.StramStore(key, bytes.NewReader(expectedFileBytes))
	require.NoError(t, err)

	readerCloser, err := store.StreamRead(key)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, readerCloser.Close())
	}()
	actualFileBytes, err := io.ReadAll(readerCloser)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(expectedFileBytes, actualFileBytes))
}

func TestDelete(t *testing.T) {
	key := "some_key_dummy"
	expectedFileBytes := []byte("hi")
	store, closer := utilGetStorage(t)
	defer closer()
	err := store.StramStore(key, bytes.NewReader(expectedFileBytes))
	require.NoError(t, err)

	err = store.Delete(key)
	assert.NoError(t, err)
	path := store.PathTransformFunc(key)
	fileInfo, err := store.RootDir.Stat(path.FullPath())
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, fileInfo)
}

func TestClear(t *testing.T) {
	store, closer := utilGetStorage(t)
	defer closer()
	assert.NoError(t, store.Clear())
}

func TestHas(t *testing.T) {
	store, closer := utilGetStorage(t)
	defer closer()

	key := "some_key_dummy"
	assert.False(t, store.Has(key))

	err := store.StramStore(key, strings.NewReader("hi"))
	require.NoError(t, err)
	assert.True(t, store.Has(key))

	err = store.Delete(key)
	require.NoError(t, err)
	assert.False(t, store.Has(key))
}

func utilGetStorage(t *testing.T) (*SimpleStorage, func()) {
	root, rootCloser := utilCreateRoot(t)
	params := SimpleStoreParam{
		PathTransformFunc: CASPathTransform,
		RootDir:           root,
	}
	store := NewSimpleStorage(params).(*SimpleStorage)
	closer := func() {
		rootCloser()
	}
	return store, closer
}

func utilCreateRoot(t *testing.T) (*os.Root, func()) {
	root, err := createRoot()
	assert.NoError(t, err)
	closer := func() {
		assert.NoError(t, root.Close())
		assert.NoError(t, os.RemoveAll("./storage_dump/"))
	}
	return root, closer
}
