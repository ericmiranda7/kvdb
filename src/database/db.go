package database

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Db struct {
	store      *os.File
	keyIndices map[string]int64
}

func NewDb(path string) Db {
	f, _ := os.OpenFile("./dbfile", os.O_RDWR|os.O_CREATE, 0644)
	return Db{
		store:      f,
		keyIndices: make(map[string]int64),
	}
}

func (d Db) ReadDbIndexes() {
	_, _ = d.store.Seek(0, 0)

	var offset int
	reader := bufio.NewReader(d.store)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		key := strings.Split(string(line), ": ")[0]
		d.keyIndices[key] = int64(offset)

		offset += len(line)
	}
	fmt.Println("offsets: ", d.keyIndices)
}

func (d Db) WriteKey(key string, val string) {
	seek, err := d.store.Seek(2, 0)
	d.keyIndices[key] = seek

	if err != nil {
		log.Fatal(err)
	}
	_, err = d.store.Write([]byte(fmt.Sprintf("%v: %v\n", key, val)))
	if err != nil {
		log.Fatal(err)
	}
}

func (d Db) ReadKey(key string) string {
	_, err := d.store.Seek(d.keyIndices[key], 0)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(d.store)

	kv, _ := reader.ReadString('\n')
	return kv
}
