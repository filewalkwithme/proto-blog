package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var wd string

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wd = workingDir
}

func main() {
	var b blog
	b.loadConfig("blog.cfg")
	b.start()
}
