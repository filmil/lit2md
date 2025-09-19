```
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
	"path"
	"strings"
)

// Cfg represents per-language configuration
type Cfg struct {
	commentStr string
	mdLang     string
}

var langMap map[string]Cfg = map[string]Cfg{
	"vhdl": Cfg{
		commentStr: "--",
		mdLang:     "vhdl",
	},
	"vhd": Cfg{
		commentStr: "--",
		mdLang:     "vhdl",
	},
	"c": Cfg{
		commentStr: "//",
		mdLang:     "c",
	},
	".cc": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	"cpp": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	"c++": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	"h": Cfg{
		commentStr: "//",
		mdLang:     "c",
	},
	"h++": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	"hpp": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	"go": Cfg{
		commentStr: "//",
		mdLang:     "go",
	},
	"sh": Cfg{
		commentStr: "#",
		mdLang:     "sh",
	},
	"bash": Cfg{
		commentStr: "#",
		mdLang:     "bash",
	},
	"py": Cfg{
		commentStr: "#",
		mdLang:     "py",
	},
	"txt": Cfg{
		commentStr: "",
		mdLang:     "",
	},
}

// prefixStr is the prefix string of a literate comment. If the language comment
// prefix is `--`, and prefixStr is `]`, then the literate comment begins with
// `--]`.
var prefixStr string = "]"

// State represents the text scanner state.
type State int

const (
	// StateNone denotes neither Code or Text.
	StateNone State = iota
	// StateCode denotes we are scanning code.
	StateCode
	// StateText denotes we are scanning code.
	StateText
)

// DocComment handles documentation comments recognition.
type DocComment struct {
	docPrefix, docPrefix2 string
}

func NewDocComment(commentStr string) DocComment {
	// Used to delineate doc prefixes
	docPrefix := fmt.Sprintf("%s%s", commentStr, prefixStr)
	docPrefix2 := fmt.Sprintf("%s%s ", commentStr, prefixStr)
	return DocComment{
		docPrefix:  docPrefix,
		docPrefix2: docPrefix2,
	}
}

func (self *DocComment) IsPrefixOf(str string) bool {
	strTrim := strings.TrimLeft(str, "\t ")
	return strings.HasPrefix(strTrim, self.docPrefix) ||
		strings.HasPrefix(strTrim, self.docPrefix2)
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

func run(inputFilename, outputFilename string, langMap map[string]Cfg) error {
	if inputFilename == "" {
		return fmt.Errorf("flag --input= is mandatory")
	}

	in, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("error while opening: %q: %v", inputFilename, err)
	}
	defer in.Close()

	ext := path.Ext(inputFilename)
	cfg, ok := langMap[ext]
	if !ok {
		cfg = Cfg{}
	}

	if outputFilename == "" {
		return fmt.Errorf("flag --output=... is mandatory")
	}

	out, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("error while opening: %q: %v", outputFilename, err)
	}
	defer out.Close()

	return convert(in, out, cfg.commentStr, cfg.mdLang)
}

func main() {

	var (
		inputFilename, outputFilename string
	)

	flag.StringVar(&inputFilename, "input", "", "input filename (code)")
	flag.StringVar(&outputFilename, "output", "", "output filename (markdown)")
	flag.StringVar(&prefixStr, "prefix", "]", "The doc comment prefix string")
	flag.Parse()

	if err := run(inputFilename, outputFilename, langMap); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
```
