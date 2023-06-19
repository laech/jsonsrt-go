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
