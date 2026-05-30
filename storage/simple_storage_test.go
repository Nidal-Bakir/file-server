package storage

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Nidal-Bakir/file-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStramStore(t *testing.T) {
	key := "some_key_dummy"
	store := utilGetStorage(t)
	defer teardown(t, store)
	err := store.StreamStore(key, strings.NewReader("hi"))
	assert.Nil(t, err)
}

func TestStramRead(t *testing.T) {
	key := "some_key_dummy"
	expectedFileBytes := []byte("hi")
	store := utilGetStorage(t)
	defer teardown(t, store)

	err := store.StreamStore(key, bytes.NewReader(expectedFileBytes))
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
	store := utilGetStorage(t)
	defer teardown(t, store)

	err := store.StreamStore(key, bytes.NewReader(expectedFileBytes))
	require.NoError(t, err)

	err = store.Delete(key)
	assert.NoError(t, err)
	path := store.PathTransformFunc(key)
	fileInfo, err := store.RootDir.Stat(path.FullPath())
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, fileInfo)
}

func TestClear(t *testing.T) {
	store := utilGetStorage(t)
	assert.NoError(t, store.Clear())
	assert.NoError(t, store.Close())
}

func TestHas(t *testing.T) {
	store := utilGetStorage(t)
	defer teardown(t, store)

	key := "some_key_dummy"
	assert.False(t, store.Has(key))

	err := store.StreamStore(key, strings.NewReader("hi"))
	require.NoError(t, err)
	assert.True(t, store.Has(key))

	err = store.Delete(key)
	require.NoError(t, err)
	assert.False(t, store.Has(key))
}

func utilGetStorage(t *testing.T) *SimpleStorage {
	root := utilCreateRoot(t)
	params := SimpleStoreParam{
		PathTransformFunc: CASPathTransform,
		RootDir:           root,
	}
	store := NewSimpleStorage(params).(*SimpleStorage)
	return store
}

func utilCreateRoot(t *testing.T) *os.Root {
	root, err := utils.OpenRoot("./storage_dump/")
	require.NoError(t, err)
	return root
}

func teardown(t *testing.T, store *SimpleStorage) {
	assert.NoError(t, store.Clear())
	assert.NoError(t, store.Close())
	assert.NoError(t, os.RemoveAll(store.RootDir.Name()))
}
