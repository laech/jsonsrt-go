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
	input  []byte
	output jsonish.Node
}{
	{[]byte("\"\""), jsonish.Value{[]byte("\"\"")}},
	{[]byte(" \"hello\""), jsonish.Value{[]byte("\"hello\"")}},
	{[]byte("123"), jsonish.Value{[]byte("123")}},
	{[]byte("{}"), jsonish.Object{[]jsonish.Member{}, false}},
	{[]byte("{}"), jsonish.Object{[]jsonish.Member{}, false}},
	{[]byte("[]"), jsonish.Array{[]jsonish.Node{}, false}},

	{[]byte("{\"a\": 1}"), jsonish.Object{
		[]jsonish.Member{
			{[]byte("\"a\""), jsonish.Value{[]byte("1")}},
		},
		false,
	}},

	{[]byte("{\"a b\": null,}"), jsonish.Object{
		[]jsonish.Member{
			{[]byte("\"a b\""), jsonish.Value{[]byte("null")}},
		},
		true,
	}},

	{[]byte("[true, null]"), jsonish.Array{
		[]jsonish.Node{
			jsonish.Value{[]byte("true")},
			jsonish.Value{[]byte("null")},
		},
		false,
	}},

	{[]byte("[0,]"), jsonish.Array{
		[]jsonish.Node{jsonish.Value{[]byte("0")}},
		true,
	}},
}
