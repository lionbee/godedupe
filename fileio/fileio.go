package fileio

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lionbee/godedupe/bytecompare"
)

// FileInfo Information to indicate if the file is a directory
type FileInfo interface {
	IsDir() bool
	Size() int64
}

// WalkFn function that is called for each element found
// during Walk
type WalkFn = func(path string, info FileInfo, err error) error

// FileIO file io operations interface
type FileIO interface {
	Delete(string) error
	FilesBytesAreEqual(string, string) bool
	MD5HashFile(string) (string, error)
	Walk(string, WalkFn) error
}

// FS file io operations
type FS struct{}

func readfile(path string) []byte {
	if data, err := ioutil.ReadFile(path); err == nil {
		return data
	}
	return nil
}

// Delete deleted the path
// If there is an error, it will be of type *PathError.
func (FS) Delete(path string) error {
	return os.Remove(path)
}

// FilesBytesAreEqual returns true if two files are exactly equal
// Will return false if a failure occurs
func (FS) FilesBytesAreEqual(path1 string, path2 string) bool {
	b1 := readfile(path1)
	b2 := readfile(path2)
	return bytecompare.BytesAreEqual(b1, b2)
}

// MD5HashFile creates an MD5 hash for a file
func (FS) MD5HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return bytecompare.MD5Hash(file)
}

// Walk walks the entire file path from to root directory calling
// walkFn for each item found
func (FS) Walk(root string, walkFn WalkFn) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return walkFn(path, info, err)
	})
}
