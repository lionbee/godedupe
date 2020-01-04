package bytecompare

import (
	"bytes"
	"errors"
	"testing"
)

type readerror struct {
	message string
}

func (r readerror) Read(p []byte) (n int, err error) {
	return 0, errors.New(r.message)
}

func TestMD5Hash(t *testing.T) {
	t.Run("MD5 is correct", func(t *testing.T) {
		expect := "0480aa34aa3db358b37cde2ab6b65326"
		received, _ := MD5Hash(bytes.NewReader([]byte("Thisisatest")), 2000)

		if expect != received {
			t.Errorf("expected %s, received %s", expect, received)
		}
	})

	t.Run("Errors are returned", func(t *testing.T) {
		expected := "FAIL"
		_, err := MD5Hash(readerror{expected}, 2000)
		received := err.Error()
		if received != expected {
			t.Errorf("Expected %s, received %s", expected, received)
		}
	})
}

func TestBytesAreEqual(t *testing.T) {
	testbytes := []byte("These are test bytes")

	t.Run("Matching", func(t *testing.T) {
		t.Parallel()
		if !BytesAreEqual(testbytes, testbytes) {
			t.Error("Bytes did not match")
		}
	})

	t.Run("Not matching", func(t *testing.T) {
		t.Parallel()
		borkedbytes := []byte("These bytes are borked")
		if BytesAreEqual(testbytes, borkedbytes) {
			t.Error("Bytes matched")
		}
	})
}
