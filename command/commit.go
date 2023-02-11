package command

import (
	"os"

	"github.com/mehm8128/git/store"
)

func generateTreeObject(filename string) {

}

func Commit(client *store.Client, filenames []string) {
	for _, filename := range filenames {
		generateTreeObject(filename)
	}
	updateIndex(os.Args[1:])
}
