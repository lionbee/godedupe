package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lionbee/godedupe/duplicates"
	"github.com/lionbee/godedupe/fileio"
)

var fs fileio.FileIO

// Filehash An MD5 hash string with the path to the hashed file
type Filehash struct {
	hash string
	path string
	size int64
}

// HashFilesInPath HashFilesInPath recursively walks the directory tree
// creating a MD5 hash for each for.
func HashFilesInPath(rootDir string) <-chan Filehash {
	hashChannel := make(chan Filehash)

	go func() {
		defer close(hashChannel)
		fs.Walk(rootDir, func(path string, info fileio.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			hash, err := fs.MD5HashFile(path)
			if err != nil {
				log.Fatal(err)
			} else {
				hashChannel <- Filehash{hash, path, info.Size()}
			}
			return nil
		})
	}()

	return hashChannel
}

func filesAreEqual(f1 *Filehash, f2 *Filehash) bool {
	return f1.size == f2.size && fs.FilesBytesAreEqual(f1.path, f2.path)
}

func emitDuplicates(hashChannel <-chan Filehash, dupesChannel chan<- duplicates.Duplicate) {
	hashMap := make(map[string]Filehash)
	for hf := range hashChannel {
		if val, ok := hashMap[hf.hash]; ok && filesAreEqual(&val, &hf) {
			dupesChannel <- duplicates.Duplicate{Value1: val.path, Value2: hf.path}
		}
		hashMap[hf.hash] = hf
	}
}

// FindDuplicates returns a new channel containing all the duplicates
// found in the hashChannel
func FindDuplicates(hashChannel <-chan Filehash) <-chan duplicates.Duplicate {
	dupesChannel := make(chan duplicates.Duplicate)

	go func() {
		defer close(dupesChannel)
		emitDuplicates(hashChannel, dupesChannel)
	}()

	return dupesChannel
}

// GetDuplicateFileDeleter returns a func that deletes the duplicate file from disc
// and writes a status message to the supplied writer
func GetDuplicateFileDeleter(writer io.Writer) duplicates.DuplicateHandler {
	return func(d duplicates.Duplicate) {
		path := d.Value2
		fmt.Fprintf(writer, "DELETING: %s\n", path)
		fs.Delete(d.Value2)
	}
}

// ProcessDuplicateFiles process all files found recursively in dir, and checks if files with matching MD5
// are equal using the provided equal function. All duplicates are sent to the supplied dupeHandler
// function
func ProcessDuplicateFiles(dir string, dupeHandler duplicates.DuplicateHandler) {
	duplicates.ApplyFuncToChan(
		FindDuplicates(
			HashFilesInPath(dir)), dupeHandler)
}

// SetFS sets an alternative FS handler
func SetFS(newFS fileio.FileIO) {
	fs = newFS
}

func main() {
	csv := flag.Bool("c", false, "Print duplicate values as a CSV to the console")
	del := flag.Bool("d", false, "Delete all duplicate values")

	flag.Parse()

	dupeHandler := duplicates.GetWriter(os.Stdout)
	if *csv {
		dupeHandler = duplicates.GetCSVWriter(os.Stdout)
	} else if *del {
		dupeHandler = GetDuplicateFileDeleter(os.Stdout)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [options] directory\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Calling without any options does a dry run and lists the files to be deleted")
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}

	dir := flag.Arg(0)

	SetFS(fileio.FS{})
	ProcessDuplicateFiles(dir, dupeHandler)
}
