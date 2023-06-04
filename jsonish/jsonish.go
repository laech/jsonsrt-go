package jsonish

import (
	"fmt"
	"io"
	"jsonsrt/lexer"
	"sort"
	"strings"
)

type Node interface {
	fmt.Stringer
	format(builder *strings.Builder, indent string, level int, applyInitalIndent bool)
	SortByName()
}

func Format(node Node) string {
	builder := strings.Builder{}
	node.format(&builder, "  ", 0, false)
	return builder.String()
}

type Value struct {
	Value string
}

func (val Value) String() string {
	return Format(val)
}

func (val Value) SortByName() {}

func (node Value) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString(node.Value)
}

type Array struct {
	Value []Node
}

func (node Array) String() string {
	return Format(node)
}

func (node Array) SortByName() {
	for _, v := range node.Value {
		v.SortByName()
	}
}

func (node Array) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString("[")

	if len(node.Value) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range node.Value {
		child.format(builder, indent, level+1, true)
		if i < len(node.Value)-1 {
			builder.WriteString(",\n")
		}
	}

	if len(node.Value) > 0 {
		builder.WriteString("\n")
		printIndent(builder, indent, level)
	}

	builder.WriteString("]")
}

type Member struct {
	Name  string
	Value Node
}

type Object struct {
	Value []Member
}

func (node Object) String() string {
	return Format(node)
}

func (node Object) SortByName() {
	for _, v := range node.Value {
		v.Value.SortByName()
	}
	sort.Slice(node.Value, func(i, j int) bool {
		return node.Value[i].Name < node.Value[j].Name
	})
}

func (node Object) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString("{")

	if len(node.Value) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range node.Value {
		printIndent(builder, indent, level+1)
		builder.WriteString(child.Name)
		builder.WriteString(": ")
		child.Value.format(builder, indent, level+1, false)
		if i < len(node.Value)-1 {
			builder.WriteString(",\n")
		}
	}

	if len(node.Value) > 0 {
		builder.WriteString("\n")
		printIndent(builder, indent, level)
	}

	builder.WriteString("}")
}

func printIndent(builder *strings.Builder, indent string, level int) {
	for i := 0; i < level; i++ {
		builder.WriteString(indent)
	}
}

func Parse(input string) (Node, error) {
	lex := lexer.New(input)
	node, err := parseNext(lex)
	if err != nil {
		return node, err
	}

	token, err := lex.Next()
	if err != io.EOF {
		return nil, fmt.Errorf("expecting EOF at offset %d, got `%s`", token.Offset, token.Value)
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
		return nil, fmt.Errorf("unexpected token at offset %d: %s", token.Offset, token.Value)
	}
}

func parseArray(lex *lexer.Lexer) (Node, error) {
	nodes := make([]Node, 0)
	for {

		token, err := lex.Next()
		if err != nil {
			return nil, nil
		}
		if token.Type == lexer.EndArray {
			return Array{nodes}, nil
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
			return Array{nodes}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return nil, fmt.Errorf("expecting value separator at offset %d, got `%s`", token.Offset, token.Value)
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
		if token.Type == lexer.EndObject {
			return Object{members}, nil
		}

		if token.Type != lexer.Value {
			return nil, fmt.Errorf("expecting member name at offset %d, got `%s`", token.Offset, token.Value)
		}
		name := token.Value

		token, err = lex.Next()
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.NameSeparator {
			return nil, fmt.Errorf("expecting name separator at offset %d, got `%s`", token.Offset, token.Value)
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
			return Object{members}, nil
		}
		if token.Type != lexer.ValueSeparator {
			return nil, fmt.Errorf("expecting value separator at offset %d, got `%s`", token.Offset, token.Value)
		}
	}
}
