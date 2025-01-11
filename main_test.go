package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestReadKey(t *testing.T) {
	inp := "a: 22\na: 42\n"
	f := bytes.Buffer{}
	f.WriteString(inp)

	keyIndices := make(map[string]int64)
	keyIndices["a"] = 6
	got := ReadKey("a", &f, keyIndices)
	expected := "a: 42\n"

	if got != expected {
		log.Fatalf("got %v want %v", got, expected)
	}
}

func TestWriteKey(t *testing.T) {
	os.Args = []string{"a", "42"}

	f := bytes.Buffer{}

	exp := "a: 42\n"

	writeKey("a", "42", 0, &f)

	got, _ := f.ReadString('\n')
	if got != exp {
		log.Fatalf("got %v want %v", got, exp)
	}
}
