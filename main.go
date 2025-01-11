package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var ind map[string]int64

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ind = make(map[string]int64)
	db, _ := os.OpenFile("./dbfile", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	readDbIndexes("./dbfile")

	var isWrite bool
	setOpts(&isWrite)

	key := flag.Arg(0)

	if isWrite {
		writeKey(key, db)
	} else {
		readKey(key, db)
	}
}

func readKey(key string, db *os.File) {
	offset := ind[key]
	reader := bufio.NewReader(db)
	_, err := reader.Discard(int(offset))
	if err != nil {
		log.Fatal(err)
	}
	kv, _ := reader.ReadString('\n')
	fmt.Print(kv)
}

func writeKey(key string, db *os.File) {
	val := flag.Arg(1)
	stat, _ := db.Stat()
	offset := stat.Size()
	ind[key] = offset + 1
	_, err := db.WriteString(fmt.Sprintf("%v: %v\n", key, val))
	if err != nil {
		log.Fatal(err)
	}
}

func setOpts(isWrite *bool) {
	flag.BoolVar(isWrite, "w", false, "Write k, v")
	flag.Parse()
}

func readDbIndexes(path string) {
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
		ind[key] = offset

		offset += int64(len(bytes))
	}
}
