package lexer

import (
	"io"
	"unicode"
)

type TokenType string

const (
	BeginObject    TokenType = "{"
	EndObject      TokenType = "}"
	BeginArray     TokenType = "["
	EndArray       TokenType = "]"
	NameSeparator  TokenType = ":"
	ValueSeparator TokenType = ","
	Value          TokenType = "Value"
)

type Token struct {
	Type   TokenType
	Value  string
	Offset int
}

type Lexer struct {
	input  []rune
	offset int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (lexer *Lexer) Next() (*Token, error) {
	lexer.skipSpaces()

	if lexer.offset >= len(lexer.input) {
		return nil, io.EOF
	}

	b := lexer.input[lexer.offset]
	switch b {
	case '{', '}', '[', ']', ',', ':':
		lexer.offset++
		return &Token{TokenType(b), string(b), lexer.offset - 1}, nil
	}

	if b == '"' {
		return lexer.readString()
	} else {
		return lexer.readStringish(), nil
	}
}

func (lexer *Lexer) skipSpaces() {
	for ; lexer.offset < len(lexer.input); lexer.offset++ {
		if !unicode.IsSpace(rune(lexer.input[lexer.offset])) {
			return
		}
	}
}

func (lexer *Lexer) readString() (*Token, error) {
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
				return &Token{Value, string(lexer.input[offset:lexer.offset]), offset}, nil
			}
			escape = false
		}
	}
	return nil, io.ErrUnexpectedEOF
}

func (lexer *Lexer) readStringish() *Token {
	for i := lexer.offset; i < len(lexer.input); i++ {
		b := lexer.input[i]
		if unicode.IsSpace(rune(b)) ||
			b == '{' ||
			b == '}' ||
			b == '[' ||
			b == ']' ||
			b == ',' ||
			b == ':' {
			offset := lexer.offset
			lexer.offset = i
			return &Token{Value, string(lexer.input[offset:i]), offset}
		}
	}
	offset := lexer.offset
	lexer.offset = len(lexer.input)
	return &Token{Value, string(lexer.input[offset:]), offset}
}
