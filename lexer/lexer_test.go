package lexer_test

import (
	"io"
	. "jsonsrt/lexer"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			tokens, err := readAllTokens(test.input)
			if err != nil {
				t.Fatalf("Lexer failed: %s", err)
			}
			if !reflect.DeepEqual(tokens, test.output) {
				t.Fatalf("\nexpected: %#v\n     got: %#v", test.output, tokens)
			}
		})
	}
}

func readAllTokens(input string) ([]Token, error) {
	lex := New(input)
	tokens := make([]Token, 0)
	for {
		token, err := lex.Next()
		if err == io.EOF {
			return tokens, nil
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
}

var tests = []struct {
	input  string
	output []Token
}{
	{"{", []Token{BeginObject(0)}},
	{"}", []Token{EndObject(0)}},
	{"[", []Token{BeginArray(0)}},
	{"]", []Token{EndArray(0)}},
	{":", []Token{NameSeparator(0)}},
	{",", []Token{ValueSeparator(0)}},
	{"\"\"", []Token{Value{0, "\"\""}}},
	{" \"hello\"", []Token{Value{1, "\"hello\""}}},
	{"\"he\tllo\"", []Token{Value{0, "\"he\tllo\""}}},
	{"\"he\\\"llo\"", []Token{Value{0, "\"he\\\"llo\""}}},
	{"\"he\\\tllo\"", []Token{Value{0, "\"he\\\tllo\""}}},
	{"123", []Token{Value{0, "123"}}},
	{"123 ", []Token{Value{0, "123"}}},
	{"{}", []Token{BeginObject(0), EndObject(1)}},
	{"[]", []Token{BeginArray(0), EndArray(1)}},

	{"{\"a\": 1}", []Token{
		BeginObject(0),
		Value{1, "\"a\""},
		NameSeparator(4),
		Value{6, "1"},
		EndObject(7),
	}},

	{"[true, null]", []Token{
		BeginArray(0),
		Value{1, "true"},
		ValueSeparator(5),
		Value{7, "null"},
		EndArray(11),
	}},
}
