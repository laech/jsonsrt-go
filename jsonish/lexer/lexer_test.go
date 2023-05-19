package lexer_test

import (
	"jsonsort/jsonish/lexer"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	for _, test := range tests {
		lexer := lexer.New(test.input)
		tokens, err := lexer.GetAll()
		if err != nil {
			t.Fatalf("Lexer failed: %s", err)
		}
		if !reflect.DeepEqual(tokens, test.output) {
			t.Fatalf("\nexpected: %s\n     got: %s\n   input: %s",
				test.output, tokens, string(test.input))
		}
	}
}

var tests = []struct {
	input  []byte
	output []lexer.Token
}{
	{[]byte("{"), []lexer.Token{beginObject(0)}},
	{[]byte("}"), []lexer.Token{endObject(0)}},
	{[]byte("["), []lexer.Token{beginArray(0)}},
	{[]byte("]"), []lexer.Token{endArray(0)}},
	{[]byte(":"), []lexer.Token{nameSeparator(0)}},
	{[]byte(","), []lexer.Token{valueSeparator(0)}},
	{[]byte("\"\""), []lexer.Token{value(0, []byte("\"\""))}},
	{[]byte(" \"hello\""), []lexer.Token{value(1, []byte("\"hello\""))}},
	{[]byte("123"), []lexer.Token{value(0, []byte("123"))}},
	{[]byte("{}"), []lexer.Token{beginObject(0), endObject(1)}},
	{[]byte("[]"), []lexer.Token{beginArray(0), endArray(1)}},

	{[]byte("{\"a\": 1}"), []lexer.Token{
		beginObject(0),
		value(1, []byte("\"a\"")),
		nameSeparator(4),
		value(6, []byte{'1'}),
		endObject(7),
	}},

	{[]byte("[true, null]"), []lexer.Token{
		beginArray(0),
		value(1, []byte("true")),
		valueSeparator(5),
		value(7, []byte("null")),
		endArray(11),
	}},
}

func beginObject(offset int) lexer.Token {
	return lexer.Token{lexer.BeginObject, []byte{'{'}, offset}
}

func endObject(offset int) lexer.Token {
	return lexer.Token{lexer.EndObject, []byte{'}'}, offset}
}

func beginArray(offset int) lexer.Token {
	return lexer.Token{lexer.BeginArray, []byte{'['}, offset}
}

func endArray(offset int) lexer.Token {
	return lexer.Token{lexer.EndArray, []byte{']'}, offset}
}

func nameSeparator(offset int) lexer.Token {
	return lexer.Token{lexer.NameSeparator, []byte{':'}, offset}
}

func valueSeparator(offset int) lexer.Token {
	return lexer.Token{lexer.ValueSeparator, []byte{','}, offset}
}

func value(offset int, val []byte) lexer.Token {
	return lexer.Token{lexer.Value, val, offset}
}
