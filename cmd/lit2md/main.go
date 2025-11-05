//] # `lit2md`: a self-referential [literate programming][lp] presentation.
//]
//] Let's first do away with the blurbs.

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

//] There is some rudimentary per-language configuration. The default "literate"
//] comment string is `//]` for C-like languages, or `--]` for example in
//] [VHDL][vhdl].  the `]` at the end was chosen to be easy to type.
//]
//] [vhdl]: https://en.wikipedia.org/wiki/VHDL

// Cfg represents per-language configuration
type Cfg struct {
	// commentStr is the comment str for this language.
	commentStr string
	// mdLang is the label used in code blocks in Markdown to activate the
	// correct syntax highlighting.
	mdLang string
}

//] Preload some language configurations.

var langMap map[string]Cfg = map[string]Cfg{
	"vhdl": Cfg{
		commentStr: "--",
		mdLang:     "vhdl",
	},
	".lua": Cfg{
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
	".cpp": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	".c++": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	".h": Cfg{
		commentStr: "//",
		mdLang:     "c",
	},
	".h++": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	".hpp": Cfg{
		commentStr: "//",
		mdLang:     "c++",
	},
	".go": Cfg{
		commentStr: "//",
		mdLang:     "go",
	},
	".sh": Cfg{
		commentStr: "#",
		mdLang:     "sh",
	},
	".bash": Cfg{
		commentStr: "#",
		mdLang:     "bash",
	},
	".py": Cfg{
		commentStr: "#",
		mdLang:     "python",
	},
	".txt": Cfg{
		commentStr: "",
		mdLang:     "",
	},
	".bazel": Cfg{
		commentStr: "#",
		mdLang:     "python",
	},
}

//] We allow the user to choose the prefix, if somehow it conflicts with other
//] uses. One potential us is to set this value to `!`, so that comments become
//] the same as Doxygen comments. I am not sure how useful that is.  This
//] is settable through flags, which are further down in the source.

// prefixStr is the prefix string of a literate comment. If the language comment
// prefix is `--`, and prefixStr is `]`, then the literate comment begins with
// `--]`.
var prefixStr string = "]"

//] The source text parsing is *extremely* simplified compared to its
//] [literate programming][lp] paragon. There are no code reorderings, there
//] are no "tangle" and "weave", because they are extremely hard to use
//] effectively in every day work, and they destroy editor cooperation. Oh well.
//]
//] [lp]: https://en.wikipedia.org/Literate_Programming
//]
//] We parse the text by simply alternating between "code" blocks and "text"
//] blocks, and filling in the appropriate code block fences as we go. This is
//] fully streaming, so you can do this to your heart's content.

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

//] We introduce a struct `DocComment` which takes on recognizing the literate
//] comment prefix.  The rules are very simple: a line is a part of a literate
//] comment if the first nonempty chars in the line is the literate comment
//] prefix.

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

// IsPrefixOf checks whether this DocComment is a prefix of `str`.
func (self *DocComment) IsPrefixOf(str string) bool {
	strTrim := strings.TrimLeft(str, "\t ")
	return strings.HasPrefix(strTrim, self.docPrefix) ||
		strings.HasPrefix(strTrim, self.docPrefix2)
}

// UnapplyPrefix removes the DocComment prefix from `str`. This results
// in a line that can be used in the text block.
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

//] `convert` makes a straight-line input-to-output conversion of the input
//] text file, on a line-buffered basis.

func convert(in io.Reader, out io.Writer, commentStr, lang string) error {
	//] This is some regular initialization. We start from `StateNone`, which
	//] allows us to eat initial empty lines.

	d := NewDocComment(commentStr)
	state := StateNone
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	_ = out

	//] This is done for each input line of text.  We unapply the prefix for
	//] each line, it's cheap to do, even if unapplied lines

	for scanner.Scan() {
		t := scanner.Text()
		np := d.UnapplyPrefix(t)

		//] I heard that goto is considered harmful.

	l:

		switch state {

		//] This is the beginning of the file. I added a `goto` just because
		//] I can.

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

		//] This is how code state is processed.  If the code ends, we also emit
		//] the Markdown end of code block.
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

		//] This is how text is processed. The idea is very similar to above.

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

	//] Don't forget to close your open code blocks at the very end.

	// If we end in the "code" state, close the code block.
	if state == StateCode {
		fmt.Fprintf(out, "```\n")
	}
	return nil
}

//] run converts `inputFilename` to `outputFilename`. We use per-language map
//] `langMap`, where the language is determined by file extension. This is
//] simplistic, but enough for now.
//]
//] `run` only exists to convert the filenames to a reader and a writer, and
//] to emit any initialization errors. This allows us to `convert` and `run`
//] in unit tests easily.

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

//] `main` declares the command line flags, and almost immediately offloads
//] to `convert`.
//]
//] Use:
//] ```
//] lit2md --help
//] ```
//] to print usage, as you'd expect.

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
