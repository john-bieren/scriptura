package main

import (
	"flag"
	"fmt"
)

// version is the current version of scriptura.
const version = "v1.0.0"

// main is scriptura's entry point.
func main() {
	// overwrite default usage function to print custom message
	flag.Usage = usage
	flag.Parse()
	processExitFlags()
	args := flag.Args()

	if len(args) == 1 {
		printPassage(args[0], "")
	} else if len(args) == 2 {
		printPassage(args[0], args[1])
	} else {
		fmt.Println("Invalid number of arguments")
		usage()
	}
}
