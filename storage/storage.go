package storage

import "io"

type Storage interface {
	StreamStore(string, io.Reader) error
	StreamRead(string) (io.ReadCloser, error)
	Delete(string) error
	Has(string) bool
	Clear() error
	Close() error
}
