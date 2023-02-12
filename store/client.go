package store

import (
	"compress/zlib"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehm8128/git/object"
	"github.com/mehm8128/git/util"
)

type Client struct {
	RootDir string
}

// pathをもらってルートディレクトリを探す
func NewClient(path string) (*Client, error) {
	rootDir, err := util.FindGitRoot(path)
	if err != nil {
		return nil, err
	}
	return &Client{
		RootDir: rootDir,
	}, nil
}

func (c *Client) GetObject(hash util.SHA1) (*object.Object, error) {
	hashString := hash.String()
	objectPath := filepath.Join(c.RootDir, "objects", hashString[:2], hashString[2:])

	objectFile, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	defer objectFile.Close()

	zr, err := zlib.NewReader(objectFile)
	if err != nil {
		return nil, err
	}
	obj, err := object.ReadObject(zr)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

type WalkFunc func(commit *object.Commit) error

func (c *Client) WalkHistory(hash util.SHA1, walkFunc WalkFunc) error {
	ancestors := []util.SHA1{hash}
	cycleCheck := map[string]struct{}{}

	for len(ancestors) > 0 {
		currentHash := ancestors[0]
		if _, ok := cycleCheck[string(currentHash)]; ok {
			ancestors = ancestors[1:]
			continue
		}
		cycleCheck[string(currentHash)] = struct{}{}

		obj, err := c.GetObject(currentHash)
		if err != nil {
			return err
		}

		current, err := object.NewCommit(obj)
		if err != nil {
			return err
		}

		if err := walkFunc(current); err != nil {
			return err
		}

		ancestors = append(ancestors[1:], current.Parents...)
	}

	return nil
}

func (c *Client) GetHeadRef() (string, error) {
	bytes, err := os.ReadFile(filepath.Join(c.RootDir, "HEAD"))
	if err != nil {
		return "", err
	}
	return strings.Split(string(bytes), " ")[1], nil
}

func (c *Client) GetHeadCommit() (util.SHA1, error) {
	headRef, err := c.GetHeadRef()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(filepath.Join(c.RootDir, headRef))
	if err != nil {
		return nil, err
	}
	return util.SHA1(bytes), nil
}
