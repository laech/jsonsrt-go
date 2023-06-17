package main

import (
	"flag"
	"fmt"
	"io"
	"jsonsrt/jsonish"
	"os"
)

const usage = `Usage: jsonsrt [OPTION]... [FILE]
Sort JSON contents.

With no FILE, read standard input and write standard output.

  --sort-by-name           sort objects by key names
  --sort-by-value KEY      sort object arrays by comparing the values of KEY
  --help                   display this help text and exit
`

var (
	sortByName  bool
	sortByValue string
	file        string
	help        bool
)

func main() {
	flag.Usage = func() {
		if _, err := fmt.Fprintf(os.Stderr, "%s", usage); err != nil {
			panic(err)
		}
	}
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&sortByName, "sort-by-name", false, "")
	flag.StringVar(&sortByValue, "sort-by-value", "", "")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() > 0 {
		file = flag.Arg(0)
	}

	var input []byte
	var err error
	if file == "" {
		input, err = io.ReadAll(os.Stdin)
	} else {
		input, err = os.ReadFile(file)
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

	if sortByName {
		node.SortByName()
	}
	if sortByValue != "" {
		node.SortByValue(sortByValue)
	}

	output := node.String() + "\n"
	if file == "" {
		fmt.Print(output)
	} else if err := os.WriteFile(file, []byte(output), 0666); err != nil {
		panic(err)
	}
}
