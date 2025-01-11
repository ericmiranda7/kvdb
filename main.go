package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	keyIndices := make(map[string]int64)
	db, _ := os.OpenFile("./dbfile", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	readDbIndexes("./dbfile", keyIndices)

	var isWrite bool
	setOpts(&isWrite)

	key := flag.Arg(0)

	if isWrite {
		writeKey(key, db, keyIndices)
	} else {
		str := ReadKey(key, db, keyIndices)
		fmt.Print(str)
	}
}

func ReadKey(key string, db io.Reader, keyIndices map[string]int64) string {
	offset := keyIndices[key]
	reader := bufio.NewReader(db)
	_, err := reader.Discard(int(offset))
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := reader.ReadString('\n')
	return kv
}

func writeKey(key string, db *os.File, keyIndices map[string]int64) {
	val := flag.Arg(1)
	stat, _ := db.Stat()
	offset := stat.Size()
	keyIndices[key] = offset + 1
	_, err := db.WriteString(fmt.Sprintf("%v: %v\n", key, val))
	if err != nil {
		log.Fatal(err)
	}
}

func setOpts(isWrite *bool) {
	flag.BoolVar(isWrite, "w", false, "Write k, v")
	flag.Parse()
}

func readDbIndexes(path string, keyIndices map[string]int64) {
	f, _ := os.Open(path)
	reader := bufio.NewReader(f)
	var offset int64
	stat, _ := f.Stat()
	fileSize := stat.Size()

	for offset < fileSize {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}

		key := strings.Split(string(bytes), ": ")[0]
		keyIndices[key] = offset

		offset += int64(len(bytes))
	}
}
