package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

var ErrorLog = log.New(os.Stdout, "Error", log.LstdFlags)
var InfoLog = log.New(os.Stdout, "Info", log.LstdFlags)
var DebugLog = log.New(os.Stdout, "Debug", log.LstdFlags)
var WarnLog = log.New(os.Stdout, "Warn", log.LstdFlags)

const (
	Base     = 10
	KiloByte = 1000
	MegaByte = 1000 * KiloByte
	GigaByte = 1000 * MegaByte
	TeraByte = 1000 * GigaByte
)

type NilEntry struct{}

func (n NilEntry) Handle(information *DuplicatesInformation) error {
	return nil
}

type EntryHandler interface {
	Handle(information *DuplicatesInformation) error
}

type DirEntry struct {
	Fullpath string
}

func (d DirEntry) Handle(information *DuplicatesInformation) error {
	dirFiles, err := ioutil.ReadDir(d.Fullpath)
	if err != nil {
		ErrorLog.Println(err)
		return err
	}
	return traverseDir(information, dirFiles, d.Fullpath)
}

type FileEntry struct {
	FullPath string
	Size     int64
}

func (f FileEntry) Handle(information *DuplicatesInformation) error {
	hashString, hashErr := generateSHA1(f.FullPath)
	if hashErr != nil {
		ErrorLog.Println(hashErr)
		return hashErr
	}
	information.AddDuplicates(hashString, f.FullPath, f.Size)
	return nil
}

func NewEntryHandler(entry os.FileInfo, directory string) EntryHandler {
	fullPath := path.Join(directory, entry.Name())
	if entry.Mode().IsDir() {
		return &DirEntry{
			Fullpath: fullPath,
		}
	}
	if entry.Mode().IsRegular() {
		return &FileEntry{
			fullPath,
			entry.Size(),
		}
	}
	return NilEntry{}
}

type DuplicatesInformation struct {
	Hashes     map[string]string
	Duplicates map[string]string
	DupeSize   *int64
}

func (d *DuplicatesInformation) AddDuplicates(hashString, fullPath string, size int64) {
	if hashEntry, ok := d.Hashes[hashString]; ok {
		d.Duplicates[hashEntry] = fullPath
		atomic.AddInt64(d.DupeSize, size)
	} else {
		d.Hashes[hashString] = fullPath
	}

}

func traverseDir(d *DuplicatesInformation, entries []os.FileInfo, directory string) error {
	for _, entry := range entries {
		entryHandler := NewEntryHandler(entry, directory)
		err := entryHandler.Handle(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateSHA1(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	if _, err := hash.Write(file); err != nil {
		return "", err
	}
	hashSum := hash.Sum(nil)
	return fmt.Sprintf("%x", hashSum), nil
}

func toReadableSize(nbytes int64) string {
	if nbytes > TeraByte {
		return strconv.FormatInt(nbytes/(TeraByte), Base) + " TB"
	}
	if nbytes > GigaByte {
		return strconv.FormatInt(nbytes/(GigaByte), Base) + " GB"
	}
	if nbytes > MegaByte {
		return strconv.FormatInt(nbytes/(MegaByte), Base) + " MB"
	}
	if nbytes > KiloByte {
		return strconv.FormatInt(nbytes/KiloByte, Base) + " KB"
	}
	return strconv.FormatInt(nbytes, Base) + " B"
}

func main() {

	var err error
	dir := flag.String("path", "./fixtures", "the path to traverse searching for duplicates")
	flag.Parse()

	if *dir == "" {
		*dir, err = os.Getwd()
		if err != nil {
			ErrorLog.Println(err)
			panic(err)
		}
	}

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		ErrorLog.Println(err)
		panic(err)
	}
	d := DuplicatesInformation{
		map[string]string{},
		map[string]string{},
		new(int64),
	}
	err = traverseDir(&d, entries, *dir)
	if err != nil {
		ErrorLog.Println(err)
		panic(err)
	}
	fmt.Println("DUPLICATES")
	fmt.Println("TOTAL UNIQUE FILES:", len(d.Hashes))
	fmt.Println("DUPLICATES:", len(d.Duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(*d.DupeSize))
}

// running into problems of not being able to open directories inside .app folders
