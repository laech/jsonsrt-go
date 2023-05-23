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
	reader *bufio.Reader
	buf    *bytes.Buffer
	offset int
}

func New(reader *bufio.Reader) *Lexer {
	return &Lexer{
		reader: reader,
		buf:    new(bytes.Buffer),
	}
}

func (lexer *Lexer) Next() (*Token, error) {
	if err := lexer.skipSpaces(); err != nil {
		return nil, err
	}

	b, err := lexer.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case '{', '}', '[', ']', ',', ':':
		offset := lexer.offset
		lexer.offset++
		return &Token{TokenType(b), []byte{b}, offset}, nil
	}

	lexer.buf.Reset()
	lexer.buf.WriteByte(b)
	if b == '"' {
		return lexer.readString()
	} else {
		return lexer.readValue()
	}
}

func (lexer *Lexer) skipSpaces() error {
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			if err := lexer.reader.UnreadRune(); err != nil {
				return err
			}
			return nil
		}
		lexer.offset++
	}
}

func (lexer *Lexer) readString() (*Token, error) {
	escape := false
	for {
		b, err := lexer.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		lexer.buf.WriteByte(b)

		if b == '\\' {
			escape = !escape
		} else if !escape && b == '"' {
			offset := lexer.offset
			lexer.offset += lexer.buf.Len()
			return &Token{Value, bytes.Clone(lexer.buf.Bytes()), offset}, nil
		}
	}
}

func (lexer *Lexer) readValue() (*Token, error) {
	for {
		b, err := lexer.reader.ReadByte()
		if err == io.EOF {
			offset := lexer.offset
			lexer.offset += lexer.buf.Len()
			return &Token{Value, bytes.Clone(lexer.buf.Bytes()), offset}, nil
		}
		if err != nil {
			return nil, err
		}
		switch b {
		case '{', '}', '[', ']', ',', ':':
			if err := lexer.reader.UnreadByte(); err != nil {
				return nil, err
			}
			offset := lexer.offset
			lexer.offset += lexer.buf.Len()
			return &Token{Value, bytes.Clone(lexer.buf.Bytes()), offset}, nil
		}
		lexer.buf.WriteByte(b)
	}
}
