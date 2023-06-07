package lexer

import (
	"io"
	"unicode"
)

type Token interface {
	Offset() int
}

type BeginObject int
type EndObject int
type BeginArray int
type EndArray int
type NameSeparator int
type ValueSeparator int
type Value struct {
	offset  int
	Content string
}

func (token BeginObject) Offset() int    { return int(token) }
func (token EndObject) Offset() int      { return int(token) }
func (token BeginArray) Offset() int     { return int(token) }
func (token EndArray) Offset() int       { return int(token) }
func (token NameSeparator) Offset() int  { return int(token) }
func (token ValueSeparator) Offset() int { return int(token) }
func (token Value) Offset() int          { return token.offset }

type Lexer struct {
	input  []rune
	offset int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (lexer *Lexer) Next() (Token, error) {
	lexer.skipSpaces()

	if lexer.offset >= len(lexer.input) {
		return nil, io.EOF
	}

	var token Token
	switch lexer.input[lexer.offset] {
	case '{':
		token = BeginObject(lexer.offset)
	case '}':
		token = EndObject(lexer.offset)
	case '[':
		token = BeginArray(lexer.offset)
	case ']':
		token = EndArray(lexer.offset)
	case ',':
		token = ValueSeparator(lexer.offset)
	case ':':
		token = NameSeparator(lexer.offset)
	}
	if token != nil {
		lexer.offset++
		return token, nil
	}

	if lexer.input[lexer.offset] == '"' {
		return lexer.readString()
	} else {
		return lexer.readStringish(), nil
	}
}

func (lexer *Lexer) skipSpaces() {
	for ; lexer.offset < len(lexer.input); lexer.offset++ {
		if !unicode.IsSpace(lexer.input[lexer.offset]) {
			return
		}
	}
}

func (lexer *Lexer) readString() (Token, error) {
	if lexer.input[lexer.offset] != '"' {
		panic("not at double quote")
	}
	escape := false
	for i := lexer.offset + 1; i < len(lexer.input); i++ {
		b := lexer.input[i]
		if b == '\\' {
			escape = !escape
		} else {
			if !escape && b == '"' {
				offset := lexer.offset
				lexer.offset = i + 1
				return Value{offset, string(lexer.input[offset:lexer.offset])}, nil
			}
			escape = false
		}
	}
	return nil, io.ErrUnexpectedEOF
}

func (lexer *Lexer) readStringish() Token {
	for i := lexer.offset; i < len(lexer.input); i++ {
		b := lexer.input[i]
		if unicode.IsSpace(b) ||
			b == '{' ||
			b == '}' ||
			b == '[' ||
			b == ']' ||
			b == ',' ||
			b == ':' {
			offset := lexer.offset
			lexer.offset = i
			return Value{offset, string(lexer.input[offset:i])}
		}
	}
	offset := lexer.offset
	lexer.offset = len(lexer.input)
	return Value{offset, string(lexer.input[offset:])}
}
