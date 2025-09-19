//] # The unit tests for `lit2md`

//] Blurbs away.

// LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

package main

import (
	"strings"
	"testing"
)

//] I prefer tabular tests

func TestOne(t *testing.T) {

	//] Marking the test as `Parallel` allows each individual test run to run
	//] separately. Not crucial but useful.
	t.Parallel()

	//] I like to define tabular tests in terms of "inputs" and expected outputs.
	//] This is very common in go.
	tests := []struct {
		summary  string
		in       string
		expected string
	}{
		//] Always start with the most trivial tests, then build up to more
		//] complex ones.
		{"empty", "", ""},
		{"basic", "Code", "```vhdl\nCode\n```\n"},
		//] This shows the first significant tests that combines a code block
		//] and a text block.
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
		//] From here on, we deal with specific interesting corner cases. Some
		//] are more annoying to deal with than others.
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

	//] Once the test table is constructed, the test is basically formulaic,
	//] and mostly reduces to presenting the errors in a way that readily
	//] shows what the problem might be.
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
