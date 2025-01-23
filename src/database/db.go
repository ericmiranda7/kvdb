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
	store    io.ReadWriteSeeker
	keyIndex map[string]int64
}

func NewDb(path string, keyIndex map[string]int64) Db {
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	return Db{
		store:    f,
		keyIndex: keyIndex,
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
		d.keyIndex[key] = int64(offset)

		offset += len(line) + 1
	}
	fmt.Println("offsets: ", d.keyIndex)
}

func (d Db) WriteKey(key string, val string) {
	seek, err := d.store.Seek(0, 2)
	if err != nil {
		log.Fatal(err)
	}
	d.keyIndex[key] = seek + 1

	_, err = d.store.Write([]byte(fmt.Sprintf("%v: %v\n", key, val)))
	if err != nil {
		log.Fatal(err)
	}
}

func (d Db) ReadKey(key string) string {
	_, err := d.store.Seek(d.keyIndex[key], 0)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(d.store)

	kv, _ := reader.ReadString('\n')
	value := strings.TrimSpace(strings.Split(kv, ":")[1][1:])
	return value
}
