package bytecompare

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
)

// MD5Hash Create a MD5 hash for a reader
func MD5Hash(src io.Reader, hashSize int64) (string, error) {
	hash := md5.New()
	if _, err := io.CopyN(hash, src, hashSize); err != nil && err != io.EOF {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// BytesAreEqual compares two bytes slices and returns true if they are equal
func BytesAreEqual(b1 []byte, b2 []byte) bool {
	return bytes.Compare(b1, b2) == 0
}
