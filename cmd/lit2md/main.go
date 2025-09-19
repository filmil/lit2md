// LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

// A conversion from simplified literate programs to markdown.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type State int

const (
	StateNone State = iota
	StateCode
	StateText
)

type DocComment struct {
	docPrefix, docPrefix2 string
}

func NewDocComment(commentStr string) DocComment {
	// Used to delineate doc prefixes
	docPrefix := fmt.Sprintf("%s]", commentStr)
	docPrefix2 := fmt.Sprintf("%s] ", commentStr)
	return DocComment{
		docPrefix:  docPrefix,
		docPrefix2: docPrefix2,
	}
}

func (self *DocComment) IsPrefixOf(str string) bool {
	strTrim := strings.TrimLeft(str, "\t ")
	return strings.HasPrefix(strTrim, self.docPrefix) || strings.HasPrefix(strTrim, self.docPrefix2)
}

func (self *DocComment) UnapplyPrefix(str string) string {
	strTrim := strings.TrimLeft(str, "\t ")

	after, cut := strings.CutPrefix(strTrim, self.docPrefix2)
	if cut {
		return after
	}
	after, cut = strings.CutPrefix(strTrim, self.docPrefix)
	if cut {
		return after
	}
	return strings.TrimRight(strTrim, "\n\r")
}

func convert(in io.Reader, out io.Writer, commentStr, lang string) error {
	d := NewDocComment(commentStr)
	state := StateNone
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	_ = out

	for scanner.Scan() {
		t := scanner.Text()
		np := d.UnapplyPrefix(t)

	l:

		switch state {
		case StateNone:
			trimmed := strings.Trim(t, "\n\t ")
			if trimmed != "" {
				// It's something, check what.
				if d.IsPrefixOf(t) {
					state = StateText
				} else {
					fmt.Fprintf(out, "```%s\n", lang) // Begin code
					state = StateCode
				}
				goto l
			}
			// Otherwise, it's empty string, so read next line.
		case StateCode:
			if d.IsPrefixOf(t) {
				// Incoming text line!
				fmt.Fprintf(out, "```\n") // End code
				fmt.Fprintf(out, "%v\n", np)
				// TODO: identify file type.
				state = StateText
			} else {
				// Incoming another code line!
				// Output verbatim.
				fmt.Fprintf(out, "%v\n", t)
			}
		case StateText:
			if d.IsPrefixOf(t) {
				// Still text.
				fmt.Fprintf(out, "%v\n", np)
			} else {
				trimmed := strings.Trim(t, "\n\t ")
				if trimmed != "" {
					// It's code now!
					fmt.Fprintf(out, "\n```%s\n", lang) // Begin code
					fmt.Fprintf(out, "%v\n", np)
					state = StateCode
				}
			}
		default:
			panic(fmt.Sprintf("unexpected main.State: %#v", state))
		}
	}

	// If we end in the "code" state, close the code block.
	if state == StateCode {
		fmt.Fprintf(out, "```\n")
	}
	return nil
}

func run(inputFilename, outputFilename string) error {
	if inputFilename == "" {
		return fmt.Errorf("flag --input= is mandatory")
	}

	in, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("error while opening: %q: %v", inputFilename, err)
	}
	defer in.Close()

	if outputFilename == "" {
		return fmt.Errorf("flag --output=... is mandatory")
	}

	out, err := os.Open(outputFilename)
	if err != nil {
		return fmt.Errorf("error while opening: %q: %v", outputFilename, err)
	}
	defer out.Close()

	return convert(in, out, "--", "vhdl")
}

func main() {

	var (
		inputFilename, outputFilename string
	)

	flag.StringVar(&inputFilename, "input", "", "input filename (code)")
	flag.StringVar(&outputFilename, "output", "", "output filename (markdown)")
	flag.Parse()

	if err := run(inputFilename, outputFilename); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
