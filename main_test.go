package main

import (
	"bytes"
	"log"
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
