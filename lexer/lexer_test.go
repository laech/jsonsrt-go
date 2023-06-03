package lexer_test

import (
	"io"
	"jsonsrt/lexer"
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

func readAllTokens(input string) ([]lexer.Token, error) {
	lex := lexer.New(input)
	tokens := make([]lexer.Token, 0)
	for {
		token, err := lex.Next()
		if err == io.EOF {
			return tokens, nil
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, *token)
	}
}

var tests = []struct {
	input  string
	output []lexer.Token
}{
	{"{", []lexer.Token{beginObject(0)}},
	{"}", []lexer.Token{endObject(0)}},
	{"[", []lexer.Token{beginArray(0)}},
	{"]", []lexer.Token{endArray(0)}},
	{":", []lexer.Token{nameSeparator(0)}},
	{",", []lexer.Token{valueSeparator(0)}},
	{"\"\"", []lexer.Token{value(0, "\"\"")}},
	{" \"hello\"", []lexer.Token{value(1, "\"hello\"")}},
	{"\"he\tllo\"", []lexer.Token{value(0, "\"he\tllo\"")}},
	{"\"he\\\"llo\"", []lexer.Token{value(0, "\"he\\\"llo\"")}},
	{"\"he\\\tllo\"", []lexer.Token{value(0, "\"he\\\tllo\"")}},
	{"123", []lexer.Token{value(0, "123")}},
	{"123 ", []lexer.Token{value(0, "123")}},
	{"{}", []lexer.Token{beginObject(0), endObject(1)}},
	{"[]", []lexer.Token{beginArray(0), endArray(1)}},

	{"{\"a\": 1}", []lexer.Token{
		beginObject(0),
		value(1, "\"a\""),
		nameSeparator(4),
		value(6, "1"),
		endObject(7),
	}},

	{"[true, null]", []lexer.Token{
		beginArray(0),
		value(1, "true"),
		valueSeparator(5),
		value(7, "null"),
		endArray(11),
	}},
}

func beginObject(offset int) lexer.Token {
	return lexer.Token{lexer.BeginObject, "{", offset}
}

func endObject(offset int) lexer.Token {
	return lexer.Token{lexer.EndObject, "}", offset}
}

func beginArray(offset int) lexer.Token {
	return lexer.Token{lexer.BeginArray, "[", offset}
}

func endArray(offset int) lexer.Token {
	return lexer.Token{lexer.EndArray, "]", offset}
}

func nameSeparator(offset int) lexer.Token {
	return lexer.Token{lexer.NameSeparator, ":", offset}
}

func valueSeparator(offset int) lexer.Token {
	return lexer.Token{lexer.ValueSeparator, ",", offset}
}

func value(offset int, val string) lexer.Token {
	return lexer.Token{lexer.Value, val, offset}
}
