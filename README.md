# lit2md: a literate program converter from code to Markdown

[![Test status](https://github.com/filmil/lit2md/workflows/Test/badge.svg)](https://github.com/filmil/lit2md/actions/workflows/test.yml)
[![Publish on Bazel Central Registry status](https://github.com/filmil/lit2md/workflows/Publish%20on%20Bazel%20Central%20Registry/badge.svg)](https://github.com/filmil/lit2md/actions/workflows/publish-bcr.yml)
[![Publish to my Bazel registry status](https://github.com/filmil/lit2md/workflows/Publish%20to%20my%20Bazel%20registry/badge.svg)](https://github.com/filmil/lit2md/actions/workflows/publish.yml)
[![Release Binaries status](https://github.com/filmil/lit2md/workflows/Release%20Binaries/badge.svg)](https://github.com/filmil/lit2md/actions/workflows/release.yml)

`lit2md` is a simplistic [literate programming (LP)][litp] converter. It takes
a source code file annotated with special "literate" programming comments, and
outputs a markdown file in which the literate programming comments are
converted into markdown text, and code is converted into code blocks.

[litp]: https://en.wikipedia.org/wiki/Literate_programming

## Example

For example, consider the simplistic `hello.cc`

```
//] # An example Hello World program.
//]
//] Don't forget headers!
#include <iostream>

//] And now, the rest of the code here:
int main() {
		std::cout << "Hello world" << std::endl;
		exit(0);
}
```

Running:

```
lit2md --input=hello.cc --output=hello.cc.md  # enclosed
```

produces [this output][this].

[this]: ./hello.cc.md

## Purpose

LP can be done to various levels of sophistication. `lit2md` does the very
basic massaging of source text.

It is not intended to be used as a stand-alone program. Rather it is intended
to be used as a front-end to a more sophisticated documentation processor, such
as [pandoc][pdc]. An example of use of a similar program is the instructive
website for programmable digital logic designers, https://fpgacpu.ca.

[pdc]: https://pandoc.org

## Prerequisites

* `bazel` to build.

Everything else is downloaded automatically.

## Build

```
bazel build //...
```

## Test

```
bazel test //...
```

## Run

```
bazel run //cmd/lit2md -- --help

```

## Q&A

### Why not reuse the `v2h.py` script from https://fpgacpu.ca/fpga/v2h.py?

* It works with python only,
* It works with Verilog only,
* It outputs HTML only.

I wanted a tad bit more freedom of choice (since I rarely code straight-up
Verilog).

### Why not use [doxygen][dxg]?

Doxygen is great, but its purpose is to provide user documentation. Literate
programs are intended to explain the code instead.

[dxg]: https://www.doxygen.nl

## References

* https://fpgacpu.ca/fpga/v2h.py: direct inspiration for this program. I
  rewrote it in go, because it is then easy to distribute as a small
  self-contained binary.
