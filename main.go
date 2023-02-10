package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/mehm8128/git/log/object"
	"github.com/mehm8128/git/log/store"
)

func main() {
	hashString := os.Args[1]
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		log.Fatal(err)
	}

	client, err := store.NewClient(".")
	if err != nil {
		log.Fatal(err)
	}
	if err := client.WalkHistory(hash, func(commit *object.Commit) error {
		fmt.Printf("\x1b[33mcommit %s\x1b[0m\n", commit.Hash.String())
		fmt.Printf("Author: %s <%s>\n", commit.Author.Name, commit.Author.Email)
		fmt.Printf("Date:   %s\n", commit.Author.Timestamp)
		fmt.Println()
		fmt.Printf("    %s\n", commit.Message)
		fmt.Println()
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
