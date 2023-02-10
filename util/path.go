package util

import (
	"errors"
	"os"
	"path/filepath"
)

func FindGitRoot(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() && file.Name() == ".git" {
			return path, nil
		}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if abs == "/" {
		return "", errors.New("not found .git directory")
	}
	return FindGitRoot(filepath.Join(path, "."))
}
