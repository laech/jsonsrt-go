package node

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"null", "null"},
		{" true", "true"},
		{"false ", "false"},
		{" 1 ", "1"},
		{"\t-2", "-2"},
		{"-3e10\n", "-3e10"},
		{"{}", "{}"},
		{"[]", "[]"},

		{`{"a":"hello"}`,
			`{
  "a": "hello"
}`,
		},

		{`{"a":"hello", "b":  [1, 2 , false]}`,
			`{
  "a": "hello",
  "b": [
    1,
    2,
    false
  ]
}`,
		},

		{`["a", "hello", null, { "i": "x"}, -1.000 ]`,
			`[
  "a",
  "hello",
  null,
  {
    "i": "x"
  },
  -1.000
]`,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			node, err := Parse(test.input)
			if err != nil {
				t.Fatalf("\nfailed: %s\n input: %s", err, test.input)
			}
			actual := node.Format()
			if actual != test.output {
				t.Fatalf("\nexpected: `%s`\n     got: `%s`\n   input: `%s`\n",
					test.output, actual, test.input)
			}
		})
	}
}
