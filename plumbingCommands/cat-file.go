package plumbing

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CatFile() {
	filename := os.Args[1]
	filepath := filename[:2] + "/" + filename[2:]
	zr, err := os.Open(fmt.Sprintf("./.git/objects/%s", filepath))
	defer zr.Close()
	if err != nil {
		panic(err)
	}

	r, err := zlib.NewReader(zr)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r)
}
