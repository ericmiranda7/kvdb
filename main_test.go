package main

import (
	"bufio"
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
	got := readKey("a", &f, keyIndices)
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

func TestReadDbIndexes(t *testing.T) {
	inp := "a: 42\nb: 36\na: 92\n"
	f := bytes.NewBufferString(inp)

	ki := make(map[string]int64)

	readDbIndexes(bufio.NewReader(f), ki)

	if v, ok := ki["a"]; !ok || v != 12 {
		log.Printf("a got %v want %v\n", v, 12)
		t.Fail()
	}
	if v, ok := ki["b"]; !ok || v != 6 {
		log.Printf("b got %v want %v\n", v, 6)
		t.Fail()
	}

}
