package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
	dirseq "github.com/sarkarshuvojit/dirseq/pkg"
)

const (
	configDirName = ".config/direseq"
)

var padding int


func setupFlags() {
	flag.IntVar(&padding, "p", 0, "Pad the output with leading zeros to reach the specified length (e.g., -p 4 => 0055)")
	flag.Parse()
}

func main() {
	setupFlags()
	dirseq.Show(configDirName, padding)
}
