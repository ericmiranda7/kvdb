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
	readDbIndexes(bufio.NewReader(db), keyIndices)

	var isWrite bool
	setOpts(&isWrite)

	key := flag.Arg(0)

	if isWrite {
		val := flag.Arg(1)
		offset := getAppendOffset(db)
		writeKey(key, val, offset, db)
	} else {
		str := readKey(key, db, keyIndices)
		fmt.Print(str)
	}
}

func getAppendOffset(db *os.File) int64 {
	stat, _ := db.Stat()
	offset := stat.Size()
	return offset
}

func readKey(key string, db io.Reader, keyIndices map[string]int64) string {
	offset := keyIndices[key]
	reader := bufio.NewReader(db)
	_, err := reader.Discard(int(offset))
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := reader.ReadString('\n')
	return kv
}

func writeKey(key string, val string, offset int64, db io.Writer) int64 {
	_, err := db.Write([]byte(fmt.Sprintf("%v: %v\n", key, val)))
	if err != nil {
		log.Fatal(err)
	}

	return offset + 1
}

func setOpts(isWrite *bool) {
	flag.BoolVar(isWrite, "w", false, "Write k, v")
	flag.Parse()
}

func readDbIndexes(db *bufio.Reader, keyIndices map[string]int64) {
	var offset int

	for {
		b, err := db.ReadBytes('\n')
		log.Println(offset)
		if err != nil {
			log.Println(offset, err)
			break
		}

		key := strings.Split(string(b), ": ")[0]
		keyIndices[key] = int64(offset)

		offset += len(b)

	}
}
