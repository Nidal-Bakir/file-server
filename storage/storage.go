package storage

import "io"

type Storage interface {
	StramStore(string, io.Reader) (string ,error)
}

