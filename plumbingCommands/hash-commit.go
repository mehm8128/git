package plumbing

import (
	"fmt"
	"os"

	"github.com/mehm8128/git/util"
)

func HashCommit() {
	filename := os.Args[1]
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	bytes = append([]byte(fmt.Sprintf("commit %d\x00", len(bytes))), bytes...)

	str := util.SHA1(bytes).Hash()
	fmt.Printf("%x\n", str)
}
