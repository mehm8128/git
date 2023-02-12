package command

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mehm8128/git/object"
	"github.com/mehm8128/git/store"
	"github.com/mehm8128/git/util"
)

func generateTreeObject(filename string) {

}

func generateCommitObject(tree util.SHA1, message string, client *store.Client) util.SHA1 {
	parent, err := client.GetHeadCommit()
	if err != nil {
		panic(err)
	}
	commit := &object.Commit{
		Hash: util.SHA1{},
		Tree: tree,
		Parents: []util.SHA1{
			parent,
		},
		Author: object.Sign{
			Name:      "mehm8128",
			Email:     "",
			Timestamp: time.Now(),
		},
		Committer: object.Sign{
			Name:      "mehm8128",
			Email:     "",
			Timestamp: time.Now(),
		},
		Message: message,
	}
	commitBytes := []byte{}
	commitBytes = append(commitBytes, []byte("tree "+tree.String()+"\n")...)
	commitBytes = append(commitBytes, []byte("parent "+parent.String()+"\n")...)
	commitBytes = append(commitBytes, []byte("author "+commit.Author.String()+"\n")...)
	commitBytes = append(commitBytes, []byte("committer "+commit.Committer.String()+"\n")...)
	commitBytes = append(commitBytes, []byte("\n")...)
	commitBytes = append(commitBytes, []byte(commit.Message)...)

	commit.Size = len(commitBytes)
	commitBytes = append([]byte(fmt.Sprintf("commit %d\x00", commit.Size)), commitBytes...)

	commit.Hash = util.Hash(commitBytes)
	hashStr := commit.Hash.String()

	fp, err := os.Create(".git/objects/" + hashStr[:2] + "/" + hashStr[2:])
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	_, err = fp.Write(commitBytes)
	if err != nil {
		panic(err)
	}
	return commit.Hash
}

func updateHead(client *store.Client, commitHash util.SHA1) {
	headRef, err := client.GetHeadRef()
	if err != nil {
		panic(err)
	}
	fp, err := os.Open(filepath.Join(client.RootDir, "objects", headRef))
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	_, err = fp.Write(commitHash)
	if err != nil {
		panic(err)
	}
}

func Commit(client *store.Client, filenames []string, message string) {
	for _, filename := range filenames {
		generateTreeObject(filename)
	}
	tree := util.SHA1{}
	commitHash := generateCommitObject(tree, message, client)
	updateHead(client, commitHash)
}
