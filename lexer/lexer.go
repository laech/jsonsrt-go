package lexer

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
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
	data   []byte
	offset int
	peeked bool
	next   Token
}

func New(data []byte) Lexer {
	return Lexer{
		data:   data,
		offset: 0,
	}
}

func (lexer *Lexer) GetAll() ([]Token, error) {
	tokens := make([]Token, 0)
	for lexer.offset < len(lexer.data) {
		token, err := lexer.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (lexer *Lexer) Peek() (Token, error) {
	if lexer.peeked {
		return lexer.next, nil
	}

	next, err := lexer.Next()
	if err != nil {
		return Token{}, err
	}
	lexer.peeked = true
	lexer.next = next
	return next, nil
}

func (lexer *Lexer) Next() (Token, error) {
	if lexer.peeked {
		next := lexer.next
		lexer.next = Token{}
		lexer.peeked = false
		return next, nil
	}

	for lexer.offset < len(lexer.data) {

		r, size := utf8.DecodeRune(lexer.data[lexer.offset:])
		if unicode.IsSpace(r) {
			lexer.offset += size
			continue
		}

		b := lexer.data[lexer.offset]
		switch b {
		case '{', '}', '[', ']', ',', ':':
			offset := lexer.offset
			lexer.offset++
			return Token{TokenType(b), []byte{b}, offset}, nil
		}

		if b == '"' {
			escape := false
			for i := lexer.offset + 1; i < len(lexer.data); i++ {
				if lexer.data[i] == '\\' {
					escape = !escape
				} else if !escape && lexer.data[i] == '"' {
					val := lexer.data[lexer.offset : i+1]
					offset := lexer.offset
					lexer.offset = i + 1
					return Token{Value, val, offset}, nil
				}
			}
			return Token{}, fmt.Errorf("expecting end of string quote, got EOF")
		}

		for i := lexer.offset + 1; i < len(lexer.data); i++ {
			switch lexer.data[i] {
			case '{', '}', '[', ']', ',', ':':
				val := lexer.data[lexer.offset:i]
				offset := lexer.offset
				lexer.offset = i
				return Token{Value, val, offset}, nil
			}
			if i == len(lexer.data)-1 {
				val := lexer.data[lexer.offset:]
				offset := lexer.offset
				lexer.offset = len(lexer.data)
				return Token{Value, val, offset}, nil
			}
		}
	}

	return Token{}, io.EOF
}
