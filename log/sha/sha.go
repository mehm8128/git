package sha

import "encoding/hex"

type SHA1 []byte

func (s SHA1) String() string {
	return hex.EncodeToString(s)
}
