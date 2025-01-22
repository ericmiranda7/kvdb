package database

import (
	"bufio"
	"fmt"
	"github.com/ericmiranda7/kvdb/v2/src/database"
	"math/rand"
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

func TestReadKey(t *testing.T) {
	data := "abc: 42\neric: 69\n"
	f, _ := os.CreateTemp("./", "dbfile")
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	_, _ = f.WriteString(data)

	indMap := map[string]int64{}
	indMap["eric"] = 8
	db := database.NewDb(f.Name(), indMap)

	got := db.ReadKey("eric")

	if got != "69" {
		t.Fatalf("got %v want %v", got, 69)
	}
}

func FuzzReadDbIndexes(f *testing.F) {
	f.Add(10, 0, 30)
	f.Fuzz(func(t *testing.T, numLines int, key int, value int) {
		if numLines > 5000 {
			t.Skip("too many lines")
		}
		tempFile, _ := os.CreateTemp("./", "dehbeh")
		defer func() {
			_ = tempFile.Close()
			_ = os.Remove(tempFile.Name())
		}()

		var lines string
		for range numLines {
			k := rune('a' + (rand.Int()+key)%26)
			lines += fmt.Sprintf("%c: %v\n", k, value+rand.Int())
		}

		aInd := int64(strings.LastIndex(lines, "a"))
		_, _ = tempFile.WriteString(lines)

		kmap := make(map[string]int64)
		db := database.NewDb(tempFile.Name(), kmap)
		db.ReadDbIndexes()

		if aInd != -1 && kmap["a"] != aInd {
			t.Log(lines)
			t.Fatalf("bruh got %v want %v", kmap["a"], aInd)
		}
	})
}
