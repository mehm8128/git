package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
)

func hashBySha1(bytes []byte) []byte {
	sha1 := sha1.New()
	sha1.Write(bytes)
	return sha1.Sum(nil)
}

func main() {
	filename := os.Args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	bytes = append([]byte(fmt.Sprintf("commit %d\x00", len(bytes))), bytes...)

	str := hashBySha1(bytes)
	fmt.Printf("%x\n", str)
}
