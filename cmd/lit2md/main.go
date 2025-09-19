// LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

// A conversion from simplified literate programs to markdown.
package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	var (
		inputFilename, outputFilename string
	)

	flag.StringVar(&inputFilename, "input", "", "input filename (code)")
	flag.StringVar(&outputFilename, "output", "", "output filename (markdown)")
	flag.Parse()

	if inputFilename == "" {
		log.Printf("flag --input= is mandatory")
		os.Exit(1)
	}

	in, err := os.Open(inputFilename)
	if err != nil {
		log.Printf("error while opening: %q: %v", inputFilename, err)
	}
}
