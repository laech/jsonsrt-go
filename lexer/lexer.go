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
	buffer *bytes.Buffer
	offset int
}

func New(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		buffer: new(bytes.Buffer),
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

	lexer.buffer.Reset()
	lexer.buffer.WriteByte(b)
	if b == '"' {
		return lexer.readString()
	} else {
		return lexer.readValue()
	}
}

func (lexer *Lexer) skipSpaces() error {
	for {
		b, err := lexer.reader.ReadByte()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(rune(b)) {
			if err := lexer.reader.UnreadByte(); err != nil {
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
		lexer.buffer.WriteByte(b)

		if b == '\\' {
			escape = !escape
		} else if !escape && b == '"' {
			offset := lexer.offset
			lexer.offset += lexer.buffer.Len()
			return &Token{Value, bytes.Clone(lexer.buffer.Bytes()), offset}, nil
		}
	}
}

func (lexer *Lexer) readValue() (*Token, error) {
	for {
		b, err := lexer.reader.ReadByte()
		if err == io.EOF {
			offset := lexer.offset
			lexer.offset += lexer.buffer.Len()
			return &Token{Value, bytes.Clone(lexer.buffer.Bytes()), offset}, nil
		}
		if err != nil {
			return nil, err
		}
		if unicode.IsSpace(rune(b)) ||
			b == '{' ||
			b == '}' ||
			b == '[' ||
			b == ']' ||
			b == ',' ||
			b == ':' {
			if err := lexer.reader.UnreadByte(); err != nil {
				return nil, err
			}
			offset := lexer.offset
			lexer.offset += lexer.buffer.Len()
			return &Token{Value, bytes.Clone(lexer.buffer.Bytes()), offset}, nil
		}
		lexer.buffer.WriteByte(b)
	}
}
