package jsonish

import (
	"fmt"
	"io"
	"jsonsrt/lexer"
)

type NodeType int

const (
	Object NodeType = iota
	Array
	Value
)

type Node struct {
	Type   NodeType
	Array  []Node
	Object []Member
	Value  []byte
}

func (node Node) String() string {
	switch node.Type {
	case Object:
		return fmt.Sprintf("Object%s", node.Object)
	case Array:
		return fmt.Sprintf("Array%s", node.Array)
	case Value:
		return fmt.Sprintf("Value[%s]", node.Value)
	default:
		return "Unknown"
	}
}

type Member struct {
	Name  []byte
	Value Node
}

func (m Member) String() string {
	return fmt.Sprintf("(%s: %s)", string(m.Name), m.Value)
}

func Parse(data []byte) (Node, error) {
	lexer := lexer.New(data)

	node, err := parseNext(&lexer)
	if err != nil {
		return node, err
	}

	token, err := lexer.Next()
	if err != io.EOF {
		return Node{}, fmt.Errorf("expecting EOF at offset %d", token.Offset)
	}

	return node, nil
}

func parseNext(lex *lexer.Lexer) (Node, error) {
	token, err := lex.Next()
	if err != nil {
		return Node{}, err
	}
	return parseCurrent(lex, token)
}

func parseCurrent(lex *lexer.Lexer, token lexer.Token) (Node, error) {
	switch token.Type {
	case lexer.BeginObject:
		return parseObject(lex)
	case lexer.BeginArray:
		return parseArray(lex)
	case lexer.Value:
		return Node{Type: Value, Value: token.Value}, nil
	default:
		return Node{}, fmt.Errorf("unexpected token at offset %d", token.Offset)
	}
}

func parseArray(lex *lexer.Lexer) (Node, error) {
	array := make([]Node, 0)
	for {

		token, err := lex.Next()
		if err != nil {
			return Node{}, nil
		}
		if token.Type == lexer.EndArray {
			return Node{Type: Array, Array: array}, nil
		}

		value, err := parseCurrent(lex, token)
		if err != nil {
			return Node{}, err
		}

		array = append(array, value)

		token, err = lex.Next()
		if err != nil {
			return Node{}, err
		}
		if token.Type == lexer.EndArray {
			return Node{Type: Array, Array: array}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return Node{}, fmt.Errorf("expecting value separator at offset %d", token.Offset)
		}
	}
}

func parseObject(lex *lexer.Lexer) (Node, error) {
	object := make([]Member, 0)
	for {
		token, err := lex.Next()
		if err != nil {
			return Node{}, err
		}
		if token.Type == lexer.EndObject {
			return Node{Type: Object, Object: object}, nil
		}

		if token.Type != lexer.Value {
			return Node{}, fmt.Errorf("expected member name at offset %d", token.Offset)
		}
		name := token.Value

		token, err = lex.Next()
		if err != nil {
			return Node{}, err
		}
		if token.Type != lexer.NameSeparator {
			return Node{}, fmt.Errorf("expecting name separator at offset %d", token.Offset)
		}

		value, err := parseNext(lex)
		if err != nil {
			return Node{}, err
		}

		object = append(object, Member{name, value})

		token, err = lex.Next()
		if err != nil {
			return Node{}, err
		}
		if token.Type == lexer.EndObject {
			return Node{Type: Object, Object: object}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return Node{}, fmt.Errorf("expecting value separator at offset %d", token.Offset)
		}
	}
}
