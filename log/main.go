package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	//hashString := os.Args[1]
	hashString := "436b835a37fec9c43adfa03ba199484893ef6afd"
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
