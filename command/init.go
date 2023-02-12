package command

import (
	"os"
)

func generateFiles() error {
	err := os.Mkdir(".git", 0777)
	if err != nil {
		return err
	}
	err = os.Mkdir(".git/objects", 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(".git/refs/heads", 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(".git/HEAD")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return err
	}
	return nil
}

func Init() {
	err := generateFiles()
	if err != nil {
		panic(err)
	}
}
