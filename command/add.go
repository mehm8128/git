package command

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/mehm8128/git/store"
	"github.com/mehm8128/git/util"
)

func generateBlobObject(filename string) error {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	//データの準備
	header := fmt.Sprintf("blob %d\x00", len(fileBytes))
	hash := util.Hash(append([]byte(header), fileBytes...))
	hashStr := fmt.Sprintf("%x", hash)

	//ファイルの準備
	fileDirectory := filepath.Join(".git", "objects", hashStr[:2])
	filePath := filepath.Join(".git", "objects", hashStr[:2], hashStr[2:])
	zr, err := util.Compress(bytes.NewBufferString(header + string(fileBytes)))
	if err != nil {
		panic(err)
	}
	//ディレクトリの存在確認と準備
	if _, err := os.Stat(fileDirectory); err != nil {
		if err := os.Mkdir(fileDirectory, 0777); err != nil {
			panic(err)
		}
	}
	//ファイルの存在確認
	if _, err := os.Stat(filePath); err == nil {
		panic("file already exists")
	}
	//書き込み
	fp, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if _, err := io.Copy(fp, zr); err != nil {
		panic(err)
	}
	return nil
}

type IndexHeader struct {
	Signature  []byte
	Version    []byte
	EntryCount []byte
}
type IndexEntry struct {
	Ctime           []byte
	CtimeNanosecond []byte
	Mtime           []byte
	MtimeNanosecond []byte
	Dev             []byte
	Ino             []byte
	Mode            []byte
	Uid             []byte
	Gid             []byte
	FileSize        []byte
	Sha1            []byte
	fileNameSize    []byte
	Name            []byte
}

func updateIndex(filenames []string) error {
	//存在確認とファイルの準備
	indexPath := filepath.Join(".git", "index")
	_, err := os.Stat(indexPath)
	var fp *os.File
	if err != nil {
		fp, err = os.Create(indexPath)
		if err != nil {
			panic(err)
		}
	} else {
		fp, err = os.OpenFile(indexPath, os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
	}
	defer fp.Close()

	//データの準備

	//ヘッダーの準備
	var header IndexHeader
	header.Signature = make([]byte, 4)
	header.Version = make([]byte, 4)
	header.EntryCount = make([]byte, 4)
	copy(header.Signature[:], []byte("DIRC"))
	copy(header.Version[:], []byte{0, 0, 0, 2})
	//todo
	copy(header.EntryCount[:], []byte{0, 0, 0, uint8(len(filenames))})
	headerByte := make([]byte, 0)
	headerByte = append(headerByte, header.Signature[:]...)
	headerByte = append(headerByte, header.Version[:]...)
	headerByte = append(headerByte, header.EntryCount[:]...)

	//エントリーの準備
	content := make([]byte, 0)
	content = append(content, headerByte...)

	entries := make([]IndexEntry, len(filenames))
	for i, filename := range filenames {
		info, err := os.Stat(filename)
		if err != nil {
			panic(err)
		}
		fileByte, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		entries[i].Ctime = make([]byte, 4)
		entries[i].CtimeNanosecond = make([]byte, 4)
		entries[i].Mtime = make([]byte, 4)
		entries[i].MtimeNanosecond = make([]byte, 4)
		entries[i].Dev = make([]byte, 4)
		entries[i].Ino = make([]byte, 4)
		entries[i].Mode = make([]byte, 4)
		entries[i].Uid = make([]byte, 4)
		entries[i].Gid = make([]byte, 4)
		entries[i].FileSize = make([]byte, 4)
		entries[i].Sha1 = make([]byte, 20)
		entries[i].fileNameSize = make([]byte, 2)
		entries[i].Name = make([]byte, len(filename))
		binary.BigEndian.PutUint32(entries[i].Ctime, uint32(info.Sys().(*syscall.Stat_t).Ctim.Sec))
		binary.BigEndian.PutUint32(entries[i].CtimeNanosecond, uint32(info.Sys().(*syscall.Stat_t).Ctim.Nsec))
		binary.BigEndian.PutUint32(entries[i].Mtime, uint32(info.Sys().(*syscall.Stat_t).Mtim.Sec))
		binary.BigEndian.PutUint32(entries[i].MtimeNanosecond, uint32(info.Sys().(*syscall.Stat_t).Mtim.Nsec))
		binary.BigEndian.PutUint32(entries[i].Dev, uint32(info.Sys().(*syscall.Stat_t).Dev))
		binary.BigEndian.PutUint32(entries[i].Ino, uint32(info.Sys().(*syscall.Stat_t).Ino))
		binary.BigEndian.PutUint32(entries[i].Mode, uint32(info.Sys().(*syscall.Stat_t).Mode))
		binary.BigEndian.PutUint32(entries[i].Uid, uint32(info.Sys().(*syscall.Stat_t).Uid))
		binary.BigEndian.PutUint32(entries[i].Gid, uint32(info.Sys().(*syscall.Stat_t).Gid))
		binary.BigEndian.PutUint32(entries[i].FileSize, uint32(info.Size()))

		entries[i].Sha1 = util.Hash(append([]byte(fmt.Sprintf("blob %d\x00", len(fileByte))), fileByte...))
		//todo
		entries[i].fileNameSize = []byte{0, uint8(len(info.Name()))}
		entries[i].Name = []byte(info.Name())

		entryByte := make([]byte, 0)
		entryByte = append(entryByte, entries[i].Ctime[:]...)
		entryByte = append(entryByte, entries[i].CtimeNanosecond[:]...)
		entryByte = append(entryByte, entries[i].Mtime[:]...)
		entryByte = append(entryByte, entries[i].MtimeNanosecond[:]...)
		entryByte = append(entryByte, entries[i].Dev[:]...)
		entryByte = append(entryByte, entries[i].Ino[:]...)
		entryByte = append(entryByte, entries[i].Mode[:]...)
		entryByte = append(entryByte, entries[i].Uid[:]...)
		entryByte = append(entryByte, entries[i].Gid[:]...)
		entryByte = append(entryByte, entries[i].FileSize[:]...)
		entryByte = append(entryByte, entries[i].Sha1[:]...)
		entryByte = append(entryByte, entries[i].fileNameSize[:]...)
		entryByte = append(entryByte, entries[i].Name[:]...)

		//todo: add flags

		content = append(content, entryByte...)
	}
	if _, err := fp.Write(content); err != nil {
		panic(err)
	}
	return nil
}

func Add(client *store.Client, filenames []string) {
	for _, filename := range filenames {
		err := generateBlobObject(filename)
		if err != nil {
			panic(err)
		}
	}
	err := updateIndex(filenames)
	if err != nil {
		panic(err)
	}
}
