package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ericmiranda7/kvdb/v2/src/database"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := database.NewDb("../dbfile")

	db.ReadDbIndexes()

	var isStdin bool
	setOpts(&isStdin)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := strings.Split(sc.Text(), " ")
		op := s[0]
		key := s[1]

		switch op {
		case "w":
			{
				val := s[2]
				db.WriteKey(key, val)
			}
		case "r":
			{
				str := db.ReadKey(key)
				fmt.Print(str)
			}
		}
	}
}

func setOpts(isStdin *bool) {
	flag.BoolVar(isStdin, "i", false, "Interactive mode")
	flag.Parse()
}
