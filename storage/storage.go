package storage

import "io"

type Storage interface {
	StramStore(string, io.Reader) error
	StreamRead(string) (io.ReadCloser, error)
	Delete(string) error
	Has(string) bool
	Clear() error
}
