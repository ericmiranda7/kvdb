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
	writeDb, _ := os.OpenFile("./dbfile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	readDb, _ := os.OpenFile("./dbfile", os.O_RDONLY, 0644)

	readDbIndexes(readDb, keyIndices)

	var isWrite bool
	var isStdin bool
	setOpts(&isWrite, &isStdin)

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		s := strings.Split(sc.Text(), " ")
		op := s[0]
		key := s[1]

		switch op {
		case "w":
			{
				val := s[2]
				keyIndices[key] = getAppendOffset(writeDb)
				writeKey(key, val, keyIndices[key], writeDb)
			}
		case "r":
			{
				println("reading")
				str := readKey(key, readDb, keyIndices)
				fmt.Print(str)
			}
		}
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
	db.(*os.File).Seek(0, 0)
	return kv
}

func writeKey(key string, val string, offset int64, db io.Writer) {
	_, err := db.Write([]byte(fmt.Sprintf("%v: %v\n", key, val)))
	if err != nil {
		log.Fatal(err)
	}
}

func setOpts(isWrite *bool, isStdin *bool) {
	flag.BoolVar(isWrite, "w", false, "Write k, v")
	flag.BoolVar(isStdin, "i", false, "Interactive mode")
	flag.Parse()
}

func readDbIndexes(db io.Reader, keyIndices map[string]int64) {
	var offset int
	reader := bufio.NewReader(db)

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		key := strings.Split(string(b), ": ")[0]
		keyIndices[key] = int64(offset)

		offset += len(b)
	}
	fmt.Println("offsets: ", keyIndices)
	_, _ = db.(*os.File).Seek(0, 0)
}
