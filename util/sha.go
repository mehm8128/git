package util

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1 []byte

func (s SHA1) String() string {
	return hex.EncodeToString(s)
}

func Hash(bytes []byte) SHA1 {
	sha1 := sha1.New()
	sha1.Write(bytes)
	return sha1.Sum(nil)
}
