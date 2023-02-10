package plumbing

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

func hexToInt64(bytes []byte) int64 {
	str := hex.EncodeToString(bytes)
	for _, v := range str {
		if v == '0' {
			str = str[1:]
		} else {
			break
		}
	}
	result, err := strconv.ParseInt(str, 16, 32)
	if err != nil {
		panic(err)
	}
	return result
}

func readEntry(bytes []byte) int64 {
	mode := bytes[24:28]
	modeInt64 := hexToInt64(mode)
	hash := bytes[40:60]
	filenameSize := bytes[60:62]
	filenameSizeInt64 := hexToInt64(filenameSize)
	filename := bytes[62 : 62+filenameSizeInt64]
	padding := 8 - (int64(62)+hexToInt64(filenameSize))%8
	fmt.Printf("%o %x %d       %s\n", modeInt64, hash, 0, filename)
	return int64(62) + hexToInt64(filenameSize) + padding
}

func IsFiles() {
	bytes, err := os.ReadFile("./.git/index")
	if err != nil {
		panic(err)
	}
	indexCount := hexToInt64(bytes[8:12])
	nextEntry := int64(12)
	for i := 0; i < int(indexCount); i++ {
		nextEntry += readEntry(bytes[nextEntry:])
	}
}
