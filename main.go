package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/pullrequestrfb/omalley/addrbook"
)

var (
	name       = flag.String("name", "nobody", "The name everyone identifies you with")
	masterAddr = flag.String("masterAddr", "127.0.0.1:4088", "address of the master instance, defaults to 127.0.0.1:4088")
	master     = flag.Bool("master", false, "whether or not this instance is the master instance")
	dbPath     = flag.String("dbpath", "/var/lib/omalley", "path to the persistent store for master, defaults to /var/lib/omalley")
)

func handleMainErr(err error) {
	if strings.Compare(os.Getenv("OMALLEY_ENV"), "production") == 0 {
		log.Fatal(err)
		return
	}
	panic(err)
}

func main() {
	abook, err := addrbook.New(*master, *dbPath)
	if err != nil {
		handleMainErr(err)
	}
}
