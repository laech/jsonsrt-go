package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type TokenType byte

const (
	BeginObject    TokenType = '{'
	EndObject      TokenType = '}'
	BeginArray     TokenType = '['
	EndArray       TokenType = ']'
	NameSeparator  TokenType = ':'
	ValueSeparator TokenType = ','
	Value          TokenType = 'v'
)

func (t TokenType) String() string {
	switch t {
	case BeginObject:
		return "BeginObject"
	case EndObject:
		return "EndObject"
	case BeginArray:
		return "BeginArray"
	case EndArray:
		return "EndArray"
	case NameSeparator:
		return "NameSeparator"
	case ValueSeparator:
		return "ValueSeparator"
	case Value:
		return "Value"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type   TokenType
	Value  []byte
	Offset int
}

func (token Token) String() string {
	return fmt.Sprintf(
		"Token{Type: %s, Value: \"%s\", Offset: %d}",
		token.Type, string(token.Value), token.Offset)
}

type Lexer struct {
	reader bufio.Reader
	buf    bytes.Buffer
	offset int
}

func New(reader bufio.Reader) Lexer {
	return Lexer{reader: reader}
}

func (lexer *Lexer) Next() (Token, error) {
	buf := &lexer.buf
	reader := &lexer.reader
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return Token{}, err
		}
		if !unicode.IsSpace(r) {
			if err := reader.UnreadRune(); err != nil {
				return Token{}, err
			}
			break
		}
		lexer.offset++
	}

	b, err := reader.ReadByte()
	if err != nil {
		return Token{}, err
	}

	switch b {
	case '{', '}', '[', ']', ',', ':':
		offset := lexer.offset
		lexer.offset++
		return Token{TokenType(b), []byte{b}, offset}, nil
	}

	buf.Reset()
	buf.WriteByte(b)

	if b == '"' {
		escape := false
		for {
			b, err = reader.ReadByte()
			if err != nil {
				return Token{}, err
			}
			buf.WriteByte(b)

			if b == '\\' {
				escape = !escape
			} else if !escape && b == '"' {
				offset := lexer.offset
				lexer.offset += buf.Len()
				return Token{Value, bytes.Clone(buf.Bytes()), offset}, nil
			}
		}
	}

	for {
		b, err = reader.ReadByte()
		if err == io.EOF {
			offset := lexer.offset
			lexer.offset += buf.Len()
			return Token{Value, bytes.Clone(buf.Bytes()), offset}, nil
		}
		if err != nil {
			return Token{}, err
		}
		switch b {
		case '{', '}', '[', ']', ',', ':':
			if err := reader.UnreadByte(); err != nil {
				return Token{}, err
			}
			offset := lexer.offset
			lexer.offset += buf.Len()
			return Token{Value, bytes.Clone(buf.Bytes()), offset}, nil
		}
		buf.WriteByte(b)
	}
}
