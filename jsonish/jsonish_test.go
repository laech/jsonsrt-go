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
	{[]byte("\"\""), jsonish.Node{Type: jsonish.Value, Value: []byte("\"\"")}},
	{[]byte(" \"hello\""), jsonish.Node{Type: jsonish.Value, Value: []byte("\"hello\"")}},
	{[]byte("123"), jsonish.Node{Type: jsonish.Value, Value: []byte("123")}},
	{[]byte("{}"), jsonish.Node{Type: jsonish.Object, Object: []jsonish.Member{}}},
	{[]byte("[]"), jsonish.Node{Type: jsonish.Array, Array: []jsonish.Node{}}},

	{[]byte("{\"a\": 1}"), jsonish.Node{
		Type: jsonish.Object,
		Object: []jsonish.Member{
			{[]byte("\"a\""), jsonish.Node{Type: jsonish.Value, Value: []byte("1")}},
		},
	}},

	{[]byte("[true, null]"), jsonish.Node{
		Type: jsonish.Array,
		Array: []jsonish.Node{
			{Type: jsonish.Value, Value: []byte("true")},
			{Type: jsonish.Value, Value: []byte("null")},
		},
	}},
}
