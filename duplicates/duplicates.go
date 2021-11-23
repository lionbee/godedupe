package duplicates

import (
	"fmt"
	"io"
)

// Duplicate Stores the details of two paths that
// are logicall equal
type Duplicate struct {
	Value1 string
	Value2 string
}

// DuplicateHandler function that handles the processing of duplicates
type DuplicateHandler func(Duplicate)

// GetWriter returns a func that will Print the duplicate value to the writer
func GetWriter(writer io.Writer) DuplicateHandler {
	return func(d Duplicate) {
		fmt.Fprintln(writer, d.Value2)
	}
}

// GetCSVWriter returns a func that prints the duplicate value and the value it is a duplicate of as a csv
func GetCSVWriter(writer io.Writer) DuplicateHandler {
	return func(d Duplicate) {
		fmt.Fprintf(writer, "\"%s\",\"%s\"\n", d.Value1, d.Value2)
	}
}

// ApplyFuncToChan iterates over a channel of Duplicate and applies the same function to each value
func ApplyFuncToChan(duplicates <-chan Duplicate, handler DuplicateHandler) {
	for d := range duplicates {
		handler(d)
	}
}
