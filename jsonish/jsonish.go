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
	SortByValue(name string)
}

type Value string
type Array []Node
type Object []Member

type Member struct {
	Name  string
	Value Node
}

func format(node Node) string {
	builder := strings.Builder{}
	node.format(&builder, "  ", 0, false)
	return builder.String()
}

func (node Value) String() string  { return format(node) }
func (node Array) String() string  { return format(node) }
func (node Object) String() string { return format(node) }

func (node Value) SortByName() {}

func (node Array) SortByName() {
	nodes := []Node(node)
	for i := range nodes {
		nodes[i].SortByName()
	}
}

func (node Object) SortByName() {
	members := []Member(node)
	for i := range members {
		members[i].Value.SortByName()
	}
	sort.Slice(members, func(i, j int) bool {
		return members[i].Name < members[j].Name
	})
}

func (node Value) SortByValue(name string) {}

func (node Array) SortByValue(name string) {
	nodes := []Node(node)
	sort.Slice(nodes, func(i, j int) bool {
		a, aOk := nodes[i].(Object)
		b, bOk := nodes[j].(Object)
		if !aOk || !bOk {
			return false
		}
		x := a.findValue(name)
		y := b.findValue(name)
		if x != nil && b != nil {
			return string(*x) < string(*y)
		}
		return false
	})
}

func (node Object) SortByValue(name string) {
	members := []Member(node)
	for i := range members {
		members[i].Value.SortByValue(name)
	}
}

func (node Value) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString(string(node))
}

func (node Array) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString("[")

	nodes := []Node(node)
	if len(nodes) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range nodes {
		child.format(builder, indent, level+1, true)
		if i < len(nodes)-1 {
			builder.WriteString(",\n")
		}
	}

	if len(nodes) > 0 {
		builder.WriteString("\n")
		printIndent(builder, indent, level)
	}

	builder.WriteString("]")
}

func (node Object) format(builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	builder.WriteString("{")

	members := []Member(node)
	if len(members) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range members {
		printIndent(builder, indent, level+1)
		builder.WriteString(child.Name)
		builder.WriteString(": ")
		child.Value.format(builder, indent, level+1, false)
		if i < len(members)-1 {
			builder.WriteString(",\n")
		}
	}

	if len(members) > 0 {
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

func (node Object) findValue(name string) *Value {
	members := []Member(node)
	for i := range members {
		if members[i].Name == name {
			if val, ok := members[i].Value.(Value); ok {
				return &val
			}
			return nil
		}
	}
	return nil
}

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
