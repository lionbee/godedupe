package duplicates

import (
	"bytes"
	"testing"
)

type MockDuplicateHandler struct {
	count int
}

func (m *MockDuplicateHandler) call(d Duplicate) {
	m.count++
}

var testdupe = Duplicate{"TEST 1", "TEST 2"}

func TestWriter(t *testing.T) {
	buffer := bytes.Buffer{}
	write := GetWriter(&buffer)

	write(testdupe)

	expected := "TEST 2\n"
	received := buffer.String()

	if received != expected {
		t.Errorf("expected %q, received %q", expected, received)
	}
}

func TestCSVWriter(t *testing.T) {
	buffer := bytes.Buffer{}
	write := GetCSVWriter(&buffer)

	write(testdupe)

	expected := "\"TEST 1\",\"TEST 2\"\n"
	received := buffer.String()

	if received != expected {
		t.Errorf("expected %q, received %q", expected, received)
	}
}

func TestApplyFuncToChan(t *testing.T) {
	t.Run("Empty channel works", func(t *testing.T) {
		mockHandler := MockDuplicateHandler{0}
		emptyChannel := make(chan Duplicate)
		close(emptyChannel)
		ApplyFuncToChan(emptyChannel, mockHandler.call)

		if mockHandler.count != 0 {
			t.Error("Handler was called")
		}
	})

	t.Run("Handler is called for each items in channel", func(t *testing.T) {
		expect := 3
		mockHandler := MockDuplicateHandler{0}
		channel := make(chan Duplicate, expect)
		for i := 0; i < expect; i++ {
			channel <- Duplicate{"test", "test"}
		}
		close(channel)

		ApplyFuncToChan(channel, mockHandler.call)

		if mockHandler.count != expect {
			t.Error("Handler was called")
		}
	})

	t.Run("Handler received items in channel", func(t *testing.T) {
		expect := Duplicate{"test", "test"}
		mockHandler := func(d Duplicate) {
			if d != expect {
				t.Error("Wrong item received")
			}
		}
		channel := make(chan Duplicate, 1)
		channel <- expect
		close(channel)

		ApplyFuncToChan(channel, mockHandler)
	})
}
