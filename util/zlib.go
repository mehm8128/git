package util

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Compress(r io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	defer zw.Close()
	if _, err := io.Copy(zw, r); err != nil {
		return nil, err
	}
	return buf, nil
}
