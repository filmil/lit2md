```go
// LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

package main

import (
	"strings"
	"testing"
)

func TestOne(t *testing.T) {
	t.Parallel()
	tests := []struct {
		summary  string
		in       string
		expected string
	}{
		{"empty", "", ""},
		{"basic", "Code", "```vhdl\nCode\n```\n"},
		{
			summary: "basic mixed",
			in: strings.Join([]string{
				"--] # Hello world",
				"--] This is text",
				"This is code.",
			}, "\n"),
			expected: strings.Join([]string{
				"# Hello world",
				"This is text",
				"",
				"```vhdl",
				"This is code.",
				"```\n",
			}, "\n"),
		},
		{
			summary: "empty line before code",
			in: strings.Join([]string{
				"",
				"--] # Hello world",
				"--] This is text",
				"This is code.",
			}, "\n"),
			expected: strings.Join([]string{
				"# Hello world",
				"This is text",
				"",
				"```vhdl",
				"This is code.",
				"```\n",
			}, "\n"),
		},
		{
			summary: "empty line befor code in the middle",
			in: strings.Join([]string{
				"",
				"--] # Hello world",
				"--] This is text",
				"",
				"This is code.",
			}, "\n"),
			expected: strings.Join([]string{
				"# Hello world",
				"This is text",
				"",
				"```vhdl",
				"This is code.",
				"```\n",
			}, "\n"),
		},
		{
			summary: "empty line befor code in the middle",
			in: `This is code.

--] This is text

This is code again.
`,
			expected: strings.Join([]string{
				"```vhdl",
				"This is code.",
				"",
				"```",
				"This is text",
				"",
				"```vhdl",
				"This is code again.",
				"```\n",
			}, "\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.summary, func(t *testing.T) {
			inRead := strings.NewReader(test.in)
			var outWrite strings.Builder

			if err := convert(inRead, &outWrite, "--", "vhdl"); err != nil {
				t.Fatalf("error while reading: %v", err)
			}

			actual := outWrite.String()

			if test.expected != actual {
				t.Errorf("want:\n%v, got:\n%v", test.expected, actual)
			}

		})
	}
}
```
