package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/mehm8128/git/log/object"
	"github.com/mehm8128/git/log/store"
)

func main() {
	//hashString := os.Args[1]
	hashString := "8afd98be900ad45406963b3204ff1427c1e3128a"
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		log.Fatal(err)
	}

	client, err := store.NewClient("testrepo")
	if err != nil {
		log.Fatal(err)
	}
	if err := client.WalkHistory(hash, func(commit *object.Commit) error {
		fmt.Println(commit)
		fmt.Println("")
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
