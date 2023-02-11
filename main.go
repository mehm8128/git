package main

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/mehm8128/git/command"
	"github.com/mehm8128/git/store"
)

func main() {
	client, err := store.NewClient(".")
	if err != nil {
		log.Fatal(err)
	}

	commandArg := os.Args[1]

	switch commandArg {
	case "add":
		command.Add(client, os.Args[2:])
	case "commit":
		command.Commit(client, os.Args[2:])
	case "log":
		hashString := os.Args[2]
		hash, err := hex.DecodeString(hashString)
		if err != nil {
			log.Fatal(err)
		}
		command.Log(client, hash)
	}
}
