package p2p

import "io"

type DefaultEncoding struct {
}

func NewDefaultEncoding() Decoder {
	return &DefaultEncoding{}
}

func (d *DefaultEncoding) Decode(r io.Reader, m *Message) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	m.Payload = buf[:n]
	return nil
}
