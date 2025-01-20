package database

import (
	"bufio"
	"github.com/ericmiranda7/kvdb/v2/src/database"
	"os"
	"strings"
	"testing"
)

func TestReadDbIndexes(t *testing.T) {
	tempFile, _ := os.CreateTemp("./", "dehbeh")
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	data := "a: 42\nb: 54\na: 92\n"
	aInd := int64(strings.LastIndex(data, "a"))
	_, _ = tempFile.WriteString(data)

	kmap := make(map[string]int64)
	db := database.NewDb(tempFile.Name(), kmap)
	db.ReadDbIndexes()

	if kmap["a"] != aInd {
		t.Fatalf("bruh got %v want %v", kmap["a"], aInd)
	}
}

func TestWriteKey(t *testing.T) {
	tempFile, _ := os.CreateTemp("./", "dehbeh")
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	kmap := make(map[string]int64)
	db := database.NewDb(tempFile.Name(), kmap)
	db.WriteKey("a", "42")

	br := bufio.NewReader(tempFile)
	got, _ := br.ReadString('\n')

	if got != "a: 42\n" {
		t.Fatalf("got %v want %v", got, "a: 42\n")
	}
}
