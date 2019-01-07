package main

import (
	"os"

	"github.com/jfontan/brahma/pkg/brahma"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: brahma <url> <file>")
	}

	repoURL := os.Args[1]
	sivaFile := os.Args[2]

	err := brahma.Download(repoURL, sivaFile)
	if err != nil {
		panic(err)
	}
}
