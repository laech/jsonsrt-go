package jsonish_test

import (
	"jsonsrt/jsonish"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	for _, test := range tests {
		node, err := jsonish.Parse(test.input)
		if err != nil {
			t.Fatalf("\nfailed: %s\n input: %s", err, string(test.input))
		}
		if !reflect.DeepEqual(node, test.output) {
			t.Fatalf("\nexpected: %s\n     got: %s\n   input: %s",
				test.output, node, string(test.input))
		}
	}
}

var tests = []struct {
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
