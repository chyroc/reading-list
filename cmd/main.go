package main

import (
	"errors"
	"fmt"
	"os"
)

const pageTitle = "chyroc's reading list"

func run() error {
	if len(os.Args) == 1 || !(os.Args[1] == "add" || os.Args[1] == "generateSite") {
		return fmt.Errorf("usage: %s [add|generateSite]", os.Args[0])
	}

	var f func() error

	switch os.Args[1] {
	case "add":
		f = AddURL
	case "generateSite":
		f = GenerateSite
	default:
		return errors.New("unreachable")
	}

	return f()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
