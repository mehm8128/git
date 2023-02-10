package sha

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1 []byte

func (s SHA1) String() string {
	return hex.EncodeToString(s)
}

func (s SHA1) Hash() []byte {
	sha1 := sha1.New()
	sha1.Write(s)
	return sha1.Sum(nil)
}
