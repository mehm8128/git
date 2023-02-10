package command

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/mehm8128/git/log/sha"
)

func compress(r io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	defer zw.Close()
	if _, err := io.Copy(zw, r); err != nil {
		return nil, err
	}
	return buf, nil
}

func generateObject(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)

	header := fmt.Sprintf("blob %d\x00", len(fileBytes))
	hash := sha.Hash(append([]byte(header), fileBytes...))
	hashStr := fmt.Sprintf("%x", hash)
	fileDirectory := filepath.Join(".git", "objects", hashStr[:2])
	filePath := filepath.Join(".git", "objects", hashStr[:2], hashStr[2:])
	zr, err := compress(bytes.NewBufferString(header + string(fileBytes)))
	if err != nil {
		panic(err)
	}
	if f, err := os.Stat(fileDirectory); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(fileDirectory, 0777); err != nil {
			panic(err)
		}
	}
	if f, err := os.Stat(filePath); !os.IsNotExist(err) && !f.IsDir() {
		panic(err)
	}
	fp, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if _, err := io.Copy(fp, zr); err != nil {
		panic(err)
	}
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

func updateIndex(filenames []string) {
	indexPath := filepath.Join(".git", "index")
	fp, err := os.OpenFile(indexPath, os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
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

	content := make([]byte, 0)
	content = append(content, headerByte...)

	entries := make([]IndexEntry, len(filenames))
	for i, filename := range filenames {
		info, err := os.Stat(filename)
		if err != nil {
			panic(err)
		}
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

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

		fileByte, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		entries[i].Sha1 = sha.Hash(append([]byte(fmt.Sprintf("blob %d\x00", len(fileByte))), fileByte...))
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
}

func main() {
	for _, filename := range os.Args[1:] {
		generateObject(filename)
	}
	updateIndex(os.Args[1:])
}
