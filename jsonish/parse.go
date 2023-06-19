package jsonish

import (
	"fmt"
	"io"
	"jsonsrt/lexer"
)

func Parse(input string) (Node, error) {
	lex := lexer.New(input)
	node, err := parseNext(lex)
	if err != nil {
		return node, err
	}

	token, err := lex.Next()
	if err != io.EOF {
		return nil, fmt.Errorf("expecting EOF at offset %d", token.Offset())
	}

	return node, nil
}

func parseNext(lex *lexer.Lexer) (Node, error) {
	token, err := lex.Next()
	if err != nil {
		return nil, err
	}
	return parseCurrent(lex, token)
}

func parseCurrent(lex *lexer.Lexer, token lexer.Token) (Node, error) {
	switch token := token.(type) {
	case lexer.BeginObject:
		return parseObject(lex)
	case lexer.BeginArray:
		return parseArray(lex)
	case lexer.Value:
		return Value(token.Content), nil
	default:
		return nil, fmt.Errorf("unexpected token at offset %d", token.Offset())
	}
}

func parseArray(lex *lexer.Lexer) (Node, error) {
	nodes := make([]Node, 0)
	for {

		token, err := lex.Next()
		if err != nil {
			return nil, nil
		}
		if _, ok := token.(lexer.EndArray); ok {
			return Array(nodes), nil
		}

		value, err := parseCurrent(lex, token)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, value)

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.EndArray); ok {
			return Array(nodes), nil
		}
		if _, ok := token.(lexer.ValueSeparator); !ok {
			return nil, fmt.Errorf("expecting value separator at offset %d", token.Offset())
		}
	}
}

func parseObject(lex *lexer.Lexer) (Node, error) {
	members := make([]Member, 0)
	for {
		token, err := lex.Next()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.EndObject); ok {
			return Object(members), nil
		}

		name, ok := token.(lexer.Value)
		if !ok {
			return nil, fmt.Errorf("expecting member name at offset %d", token.Offset())
		}

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.NameSeparator); !ok {
			return nil, fmt.Errorf("expecting name separator at offset %d", token.Offset())
		}

		value, err := parseNext(lex)
		if err != nil {
			return nil, err
		}

		members = append(members, Member{name.Content, value})

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.EndObject); ok {
			return Object(members), nil
		}
		if _, ok := token.(lexer.ValueSeparator); !ok {
			return nil, fmt.Errorf("expecting value separator at offset %d", token.Offset())
		}
	}
}
