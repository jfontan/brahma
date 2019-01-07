package main

import (
	"fmt"
	"os"

	"github.com/jfontan/brahma/pkg/brahma"
	"github.com/sanity-io/litter"
)

func client() error {
	c, err := brahma.NewClient("http://localhost:8765")
	if err != nil {
		return err
	}

	for {
		repo, err := c.Repository()
		if err != nil {
			return err
		}

		litter.Dump(repo)
	}
}

func server() error {
	return brahma.StartServer()
}

func main() {
	if len(os.Args) < 2 {
		panic("usage: brahma <server|client> [<arguments>]")
	}

	var err error
	switch os.Args[1] {
	case "client":
		err = client()
	case "server":
		err = server()
	default:
		err = fmt.Errorf("unknown mode %s", os.Args[1])
	}

	if err != nil {
		panic(err)
	}
}
