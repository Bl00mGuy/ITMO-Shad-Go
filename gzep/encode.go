//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func Encode(data []byte, w io.Writer) error {
	ww := gzipWriterPool.Get().(*gzip.Writer)
	defer gzipWriterPool.Put(ww)

	ww.Reset(w)
	defer func() { _ = ww.Close() }()

	if _, err := ww.Write(data); err != nil {
		return err
	}
	return ww.Flush()
}
