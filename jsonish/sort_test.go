package jsonish

import (
	"fmt"
	"reflect"
	"testing"
)

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
			Object{{`"1 "`, Value("x")}, {`"1"`, Value("x")}},
			Object{{`"1"`, Value("x")}, {`"1 "`, Value("x")}},
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
		t.Run(fmt.Sprintf("%#v", test.input), func(t *testing.T) {
			test.input.SortByName()
			if !reflect.DeepEqual(test.input, test.output) {
				t.Fatalf("\nexpected: `%#v`\n     got: `%#v`\n",
					test.output, test.input)
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
				Object{{`"name"`, Value(`"x "`)}},
				Object{{`"name"`, Value(`"x"`)}},
			},
			Array{
				Object{{`"name"`, Value(`"x"`)}},
				Object{{`"name"`, Value(`"x "`)}},
			},
		},
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
		{
			"a",
			Array{
				Object{
					{`"a"`, Value("2")},
					{`"b"`, Array{
						Object{{`"a"`, Value("y")}},
						Object{{`"a"`, Value("x")}},
					}},
				},
				Object{{`"a"`, Value("0")}},
			},
			Array{
				Object{{`"a"`, Value("0")}},
				Object{
					{`"a"`, Value("2")},
					{`"b"`, Array{
						Object{{`"a"`, Value("x")}},
						Object{{`"a"`, Value("y")}},
					}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.input), func(t *testing.T) {
			test.input.SortByValue(test.name)
			if !reflect.DeepEqual(test.input, test.output) {
				t.Fatalf("\nexpected: `%#v`\n     got: `%#v`\n",
					test.output, test.input)
			}
		})
	}
}
