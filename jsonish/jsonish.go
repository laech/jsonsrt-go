package jsonish

import (
	"bytes"
	"fmt"
	"io"
	"jsonsrt/lexer"
)

type Node interface {
	fmt.Stringer
}

func nodeString(node Node) string {
	buf := bytes.Buffer{}
	print(node, &buf, []byte("  "), 0, false)
	return buf.String()
}

type Value struct {
	Value []byte
}

type Array struct {
	Value       []Node
	TrailingSep bool
}

type Object struct {
	Value       []Member
	TrailingSep bool
}

type Member struct {
	Name  []byte
	Value Node
}

func (val Value) String() string {
	return nodeString(val)
}

func (arr Array) String() string {
	return nodeString(arr)
}

func (obj Object) String() string {
	return nodeString(obj)
}

func Parse(reader io.Reader) (Node, error) {
	lex := lexer.New(reader)
	node, err := parseNext(lex)
	if err != nil {
		return node, err
	}

	token, err := lex.Next()
	if err != io.EOF {
		return nil, fmt.Errorf("expecting EOF at offset %d", token.Offset)
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

func parseCurrent(lex *lexer.Lexer, token *lexer.Token) (Node, error) {
	switch token.Type {
	case lexer.BeginObject:
		return parseObject(lex)
	case lexer.BeginArray:
		return parseArray(lex)
	case lexer.Value:
		return Value{token.Value}, nil
	default:
		return nil, fmt.Errorf("unexpected token at offset %d", token.Offset)
	}
}

func parseArray(lex *lexer.Lexer) (Node, error) {
	nodes := make([]Node, 0)
	for i := 0; ; i++ {

		token, err := lex.Next()
		if err != nil {
			return nil, nil
		}
		if token.Type == lexer.EndArray {
			return Array{nodes, i > 0}, nil
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
		if token.Type == lexer.EndArray {
			return Array{nodes, false}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return nil, fmt.Errorf("expecting value separator at offset %d", token.Offset)
		}
	}
}

func parseObject(lex *lexer.Lexer) (Node, error) {
	members := make([]Member, 0)
	for i := 0; ; i++ {
		token, err := lex.Next()
		if err != nil {
			return nil, err
		}
		if token.Type == lexer.EndObject {
			return Object{members, i > 0}, nil
		}

		if token.Type != lexer.Value {
			return nil, fmt.Errorf("expecting member name at offset %d", token.Offset)
		}
		name := token.Value

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.NameSeparator {
			return nil, fmt.Errorf("expecting name separator at offset %d", token.Offset)
		}

		value, err := parseNext(lex)
		if err != nil {
			return nil, err
		}

		members = append(members, Member{name, value})

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if token.Type == lexer.EndObject {
			return Object{members, false}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return nil, fmt.Errorf("expecting value separator at offset %d", token.Offset)
		}
	}
}
