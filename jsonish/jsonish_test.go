package jsonish

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		output Node
	}{
		{`""`, Value(`""`)},
		{` "hello"`, Value(`"hello"`)},
		{"123", Value("123")},
		{"{}", Object{}},
		{"{}", Object{}},
		{"[]", Array{}},
		{"{\"a\": 1}", Object{{`"a"`, Value("1")}}},
		{"{\"a b\": null,}", Object{{`"a b"`, Value("null")}}},
		{"[true, null]", Array{Value("true"), Value("null")}},
		{"[0,]", Array{Value("0")}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			node, err := Parse(test.input)
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
			node, err := Parse(test.input)
			if err != nil {
				t.Fatalf("\nfailed: %s\n input: %s", err, test.input)
			}
			actual := node.String()
			if actual != test.output {
				t.Fatalf("\nexpected: `%s`\n     got: `%s`\n   input: `%s`\n",
					test.output, actual, test.input)
			}
		})
	}
}

func TestSortByName(t *testing.T) {
	tests := []struct {
		input  Node
		output Node
	}{
		{Value("1"), Value("1")},
		{Object{}, Object{}},
		{
			Object{{"1", Value("a")}},
			Object{{"1", Value("a")}},
		},
		{
			Object{{"1", Value("a")}, {"2", Value("b")}},
			Object{{"1", Value("a")}, {"2", Value("b")}},
		},
		{
			Object{{"2", Value("b")}, {"1", Value("a")}},
			Object{{"1", Value("a")}, {"2", Value("b")}},
		},
		{
			Object{
				{"2", Value("b")},
				{"1", Value("a")},
				{"3", Object{
					{"1", Value("one")},
					{"0", Value("zero")},
				}},
			},
			Object{
				{"1", Value("a")},
				{"2", Value("b")},
				{"3", Object{
					{"0", Value("zero")},
					{"1", Value("one")},
				}},
			},
		},
		{
			Object{
				{"2", Value("b")},
				{"1", Value("a")},
				{"3", Array{Object{
					{"1", Value("one")},
					{"0", Value("zero")},
				}}},
			},
			Object{
				{"1", Value("a")},
				{"2", Value("b")},
				{"3", Array{Object{
					{"0", Value("zero")},
					{"1", Value("one")}},
				}},
			},
		},
		{Array{}, Array{}},
		{
			Array{Object{
				{"1", Value("one")},
				{"0", Value("zero")},
			}},
			Array{Object{
				{"0", Value("zero")},
				{"1", Value("one")},
			}},
		},
		{
			Array{Object{
				{"1", Value("one")},
				{"0", Array{Object{
					{"y", Value("yy")},
					{"x", Value("xx")},
				}}},
			}},
			Array{Object{
				{"0", Array{Object{
					{"x", Value("xx")},
					{"y", Value("yy")},
				}}},
				{"1", Value("one")},
			}},
		},
	}

	for _, test := range tests {
		str := test.input.String()
		t.Run(str, func(t *testing.T) {
			test.input.SortByName()
			if !reflect.DeepEqual(test.input, test.output) {
				t.Fatalf("\nexpected: `%s`\n     got: `%s`\n   input: `%s`\n",
					test.output, test.input, str)
			}
		})
	}
}

func TestSortByValue(t *testing.T) {
	tests := []struct {
		name   string
		input  Node
		output Node
	}{
		{"", Value("1"), Value("1")},
		{"", Object{}, Object{}},
		{"", Array{}, Array{}},
		{
			"name",
			Array{
				Object{{`"name"`, Value("1")}},
				Object{{`"name"`, Value("2")}},
			},
			Array{
				Object{{`"name"`, Value("1")}},
				Object{{`"name"`, Value("2")}},
			},
		},
		{
			"name",
			Array{
				Object{{`"name"`, Value("2")}},
				Object{{`"name"`, Value("1")}},
			},
			Array{
				Object{{`"name"`, Value("1")}},
				Object{{`"name"`, Value("2")}},
			},
		},
		{
			"name",
			Object{
				{`"name"`, Array{
					Object{{`"name"`, Value("2")}},
					Object{{`"name"`, Value("1")}},
				}},
			},
			Object{
				{`"name"`, Array{
					Object{{`"name"`, Value("1")}},
					Object{{`"name"`, Value("2")}},
				}},
			},
		},
		{
			"a",
			Array{
				Object{{`"a"`, Value("1")}},
				Object{{`"a"`, Value("2")}},
				Object{{`"a"`, Value("0")}},
			},
			Array{
				Object{{`"a"`, Value("0")}},
				Object{{`"a"`, Value("1")}},
				Object{{`"a"`, Value("2")}},
			},
		},
	}

	for _, test := range tests {
		str := test.input.String()
		t.Run(str, func(t *testing.T) {
			test.input.SortByValue(test.name)
			if !reflect.DeepEqual(test.input, test.output) {
				t.Fatalf("\nexpected: `%s`\n     got: `%s`\n   input: `%s`\n",
					test.output, test.input, str)
			}
		})
	}
}
