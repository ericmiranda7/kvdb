package main

import (
	"bufio"
	"github.com/ericmiranda7/kvdb/v2/src/database"
	"log"
	"net/http"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := database.NewDb("./dbfile", make(map[string]int64))

	db.ReadDbIndexes()

	smux := http.DefaultServeMux
	s := http.Server{
		Addr:    "localhost:8081",
		Handler: smux,
	}
	smux.HandleFunc("POST /", writeHandler(db))
	smux.HandleFunc("GET /{key}", readHandler(db))

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func readHandler(db database.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		println(key)
		value := db.ReadKey(key)
		w.Write([]byte(value))
	}
}

func writeHandler(db database.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sc := bufio.NewScanner(r.Body)
		sc.Scan()
		kv := sc.Text()
		key := strings.Split(kv, ": ")[0]
		val := strings.Split(kv, ": ")[1]
		db.WriteKey(key, val)
		w.WriteHeader(http.StatusOK)
	}
}
