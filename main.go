package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

const (
	Base     = 10
	KiloByte = 1000
	MegaByte = 1000 * KiloByte
	GigaByte = 1000 * MegaByte
	TeraByte = 1000 * GigaByte
)

func traverseDir(hashes, duplicates map[string]string, dupeSize *int64, entries []os.FileInfo, directory string) {
	for _, entry := range entries {
		fullpath := path.Join(directory, entry.Name())
		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}
		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				panic(err)
			}
			traverseDir(hashes, duplicates, dupeSize, dirFiles, fullpath)
			continue
		}
		hashString, hashErr := generateSHA1(fullpath)
		if hashErr != nil {
			panic(hashErr)
		}
		if hashEntry, ok := hashes[hashString]; ok {
			duplicates[hashEntry] = fullpath
			atomic.AddInt64(dupeSize, entry.Size())
		} else {
			hashes[hashString] = fullpath
		}
	}
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
	dir := flag.String("path", "", "the path to traverse searching for duplicates")
	flag.Parse()

	if *dir == "" {
		*dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	hashes := map[string]string{}
	duplicates := map[string]string{}
	var dupeSize int64

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		panic(err)
	}

	traverseDir(hashes, duplicates, &dupeSize, entries, *dir)

	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(hashes))
	fmt.Println("DUPLICATES:", len(duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(dupeSize))
}

// running into problems of not being able to open directories inside .app folders
