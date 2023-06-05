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
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
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

	println(node.String())
}
