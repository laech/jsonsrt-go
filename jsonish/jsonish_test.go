package jsonish_test

import (
	"jsonsrt/jsonish"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		output jsonish.Node
	}{
		{"\"\"", jsonish.Value{"\"\""}},
		{" \"hello\"", jsonish.Value{"\"hello\""}},
		{"123", jsonish.Value{"123"}},
		{"{}", jsonish.Object{[]jsonish.Member{}, false}},
		{"{}", jsonish.Object{[]jsonish.Member{}, false}},
		{"[]", jsonish.Array{[]jsonish.Node{}, false}},

		{"{\"a\": 1}", jsonish.Object{
			[]jsonish.Member{
				{"\"a\"", jsonish.Value{"1"}},
			},
			false,
		}},

		{"{\"a b\": null,}", jsonish.Object{
			[]jsonish.Member{
				{"\"a b\"", jsonish.Value{"null"}},
			},
			true,
		}},

		{"[true, null]", jsonish.Array{
			[]jsonish.Node{
				jsonish.Value{"true"},
				jsonish.Value{"null"},
			},
			false,
		}},

		{"[0,]", jsonish.Array{
			[]jsonish.Node{jsonish.Value{"0"}},
			true,
		}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			node, err := jsonish.Parse(test.input)
			if err != nil {
				t.Fatalf("\nfailed: %s\n input: %s", err, test.input)
			}
			if !reflect.DeepEqual(node, test.output) {
				t.Fatalf("\nexpected: %s\n     got: %s\n   input: %s",
					test.output, node, test.input)
			}
		})
	}
}

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
			node, err := jsonish.Parse(test.input)
			if err != nil {
				t.Fatalf("\nfailed: %s\n input: %s", err, test.input)
			}
			actual := jsonish.Format(node)
			if actual != test.output {
				t.Fatalf("\nexpected: `%s`\n     got: `%s`\n   input: `%s`\n",
					test.output, actual, test.input)
			}
		})
	}
}
