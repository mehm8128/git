package store

import (
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehm8128/git/object"
	"github.com/mehm8128/git/util"
)

type Client struct {
	objectDir string
}

// pathをもらってルートディレクトリを探す
func NewClient(path string) (*Client, error) {
	rootDir, err := util.FindGitRoot(path)
	if err != nil {
		return nil, err
	}
	return &Client{
		objectDir: filepath.Join(rootDir, ".git", "objects"),
	}, nil
}

func (c *Client) GetObject(hash util.SHA1) (*object.Object, error) {
	hashString := hash.String()
	objectPath := filepath.Join(c.objectDir, hashString[:2], hashString[2:])

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

func (c *Client) GetHeadCommit() (util.SHA1, error) {
	fp, err := os.Open(filepath.Join(c.objectDir, "HEAD"))
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	//refのパスを取得してopen
	bytes, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	headRef := strings.Split(string(bytes), " ")[1]
	fp2, err := os.Open(filepath.Join(c.objectDir, headRef))
	if err != nil {
		return nil, err
	}
	defer fp2.Close()
	bytes, err = io.ReadAll(fp2)
	if err != nil {
		return nil, err
	}
	return util.SHA1(bytes), nil
}
