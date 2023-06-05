package main

import (
	"flag"
	"fmt"
	"io"
	"jsonsrt/jsonish"
	"os"
)

func main() {
	sortByName := flag.Bool("sort-by-name", false, "sort by name")
	sortByValue := flag.String("sort-by-value", "", "sort by value")
	path := flag.String("file", "-", "file")
	flag.Parse()

	var input []byte
	var err error
	if *path == "-" {
		input, err = io.ReadAll(os.Stdin)
	} else {
		input, err = os.ReadFile(*path)
	}
	if err != nil {
		if _, err := os.Stderr.WriteString(err.Error()); err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	node, err := jsonish.Parse(string(input))
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, "failed to parse input:", err.Error()); err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	if *sortByName {
		node.SortByName()
	}
	if *sortByValue != "" {
		node.SortByValue(*sortByValue)
	}

	if *path == "-" {
		println(node.String())
	} else {
		if err := os.WriteFile(*path, []byte(node.String()+"\n"), 0666); err != nil {
			panic(err)
		}
	}
}
