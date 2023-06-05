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
		{"\"\"", jsonish.Value("\"\"")},
		{" \"hello\"", jsonish.Value("\"hello\"")},
		{"123", jsonish.Value("123")},
		{"{}", jsonish.Object([]jsonish.Member{})},
		{"{}", jsonish.Object([]jsonish.Member{})},
		{"[]", jsonish.Array([]jsonish.Node{})},

		{"{\"a\": 1}", jsonish.Object(
			[]jsonish.Member{
				{"\"a\"", jsonish.Value("1")},
			},
		)},

		{"{\"a b\": null,}", jsonish.Object(
			[]jsonish.Member{
				{"\"a b\"", jsonish.Value("null")},
			},
		)},

		{"[true, null]", jsonish.Array(
			[]jsonish.Node{
				jsonish.Value("true"),
				jsonish.Value("null"),
			},
		)},

		{"[0,]", jsonish.Array(
			[]jsonish.Node{jsonish.Value("0")},
		)},
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
		input  jsonish.Node
		output jsonish.Node
	}{
		{jsonish.Value("1"), jsonish.Value("1")},
		{
			jsonish.Object([]jsonish.Member{}),
			jsonish.Object([]jsonish.Member{}),
		},
		{
			jsonish.Object([]jsonish.Member{{"1", jsonish.Value("a")}}),
			jsonish.Object([]jsonish.Member{{"1", jsonish.Value("a")}}),
		},
		{
			jsonish.Object([]jsonish.Member{{"1", jsonish.Value("a")}, {"2", jsonish.Value("b")}}),
			jsonish.Object([]jsonish.Member{{"1", jsonish.Value("a")}, {"2", jsonish.Value("b")}}),
		},
		{
			jsonish.Object([]jsonish.Member{{"2", jsonish.Value("b")}, {"1", jsonish.Value("a")}}),
			jsonish.Object([]jsonish.Member{{"1", jsonish.Value("a")}, {"2", jsonish.Value("b")}}),
		},
		{
			jsonish.Object(
				[]jsonish.Member{
					{"2", jsonish.Value("b")},
					{"1", jsonish.Value("a")},
					{"3", jsonish.Object(
						[]jsonish.Member{
							{"1", jsonish.Value("one")},
							{"0", jsonish.Value("zero")}})}}),
			jsonish.Object(
				[]jsonish.Member{
					{"1", jsonish.Value("a")},
					{"2", jsonish.Value("b")},
					{"3", jsonish.Object(
						[]jsonish.Member{
							{"0", jsonish.Value("zero")},
							{"1", jsonish.Value("one")}})}}),
		},
		{
			jsonish.Object(
				[]jsonish.Member{
					{"2", jsonish.Value("b")},
					{"1", jsonish.Value("a")},
					{"3", jsonish.Array(
						[]jsonish.Node{
							jsonish.Object(
								[]jsonish.Member{
									{"1", jsonish.Value("one")},
									{"0", jsonish.Value("zero")}})})}}),
			jsonish.Object(
				[]jsonish.Member{
					{"1", jsonish.Value("a")},
					{"2", jsonish.Value("b")},
					{"3", jsonish.Array(
						[]jsonish.Node{
							jsonish.Object(
								[]jsonish.Member{
									{"0", jsonish.Value("zero")},
									{"1", jsonish.Value("one")}})})}}),
		},
		{
			jsonish.Array([]jsonish.Node{}),
			jsonish.Array([]jsonish.Node{}),
		},
		{
			jsonish.Array(
				[]jsonish.Node{
					jsonish.Object(
						[]jsonish.Member{
							{"1", jsonish.Value("one")},
							{"0", jsonish.Value("zero")}})}),
			jsonish.Array(
				[]jsonish.Node{
					jsonish.Object(
						[]jsonish.Member{
							{"0", jsonish.Value("zero")},
							{"1", jsonish.Value("one")}})}),
		},
		{
			jsonish.Array(
				[]jsonish.Node{
					jsonish.Object(
						[]jsonish.Member{
							{"1", jsonish.Value("one")},
							{"0", jsonish.Array(
								[]jsonish.Node{
									jsonish.Object(
										[]jsonish.Member{
											{"y", jsonish.Value("yy")},
											{"x", jsonish.Value("xx")}})})}})}),
			jsonish.Array(
				[]jsonish.Node{
					jsonish.Object(
						[]jsonish.Member{
							{"0", jsonish.Array(
								[]jsonish.Node{
									jsonish.Object(
										[]jsonish.Member{
											{"x", jsonish.Value("xx")},
											{"y", jsonish.Value("yy")}})})},
							{"1", jsonish.Value("one")}})}),
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
