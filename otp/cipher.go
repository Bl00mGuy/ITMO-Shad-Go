//go:build !solution

package otp

import "io"

type reader struct {
	r    io.Reader
	prng io.Reader
}

type writer struct {
	w    io.Writer
	prng io.Reader
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &reader{
		r:    r,
		prng: prng,
	}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &writer{
		w:    w,
		prng: prng,
	}
}

func (reader *reader) Read(data []byte) (int, error) {
	n, err := reader.readData(&data)
	if n > 0 {
		prngBytes, prgnErr := reader.readPRNGBytes(&n)
		if reader.handleError(prgnErr) != nil {
			return 0, prgnErr
		}
		reader.xorBytes(&data, &prngBytes, n)
	}
	if err == io.EOF && n > 0 {
		return n, nil
	}
	return n, err
}

func (reader *reader) readData(data *[]byte) (int, error) {
	return reader.r.Read(*data)
}

func (reader *reader) readPRNGBytes(n *int) ([]byte, error) {
	prngBytes := make([]byte, *n)
	_, err := reader.prng.Read(prngBytes)
	return prngBytes, err
}

func (reader *reader) xorBytes(data, prngBytes *[]byte, n int) {
	for i := 0; i < n; i++ {
		(*data)[i] ^= (*prngBytes)[i]
	}
}

func (reader *reader) handleError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func (writer *writer) Write(data []byte) (int, error) {
	n := len(data)
	prngBytes, err := writer.readPRNGBytes(&n)
	if writer.handleError(err) != nil {
		return 0, err
	}
	encrypted := writer.applyXOR(&data, &prngBytes)
	return writer.writeData(&encrypted)
}

func (writer *writer) readPRNGBytes(n *int) ([]byte, error) {
	prngBytes := make([]byte, *n)
	_, err := writer.prng.Read(prngBytes)
	return prngBytes, err
}

func (writer *writer) applyXOR(data, prngBytes *[]byte) []byte {
	encrypted := make([]byte, len(*data))
	for i := range *data {
		encrypted[i] = (*data)[i] ^ (*prngBytes)[i]
	}
	return encrypted
}

func (writer *writer) writeData(data *[]byte) (int, error) {
	n, err := writer.w.Write(*data)
	if err == io.EOF && n > 0 {
		return n, nil
	}
	return n, err
}

func (writer *writer) handleError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
