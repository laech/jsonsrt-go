package format

import (
	. "jsonsrt/jsonish"
	"strings"
)

func Print(node Node) string {
	builder := strings.Builder{}
	print(node, &builder, "  ", 0, false)
	return builder.String()
}

func print(node Node, builder *strings.Builder, indent string, level int, applyInitialIndent bool) {
	if applyInitialIndent {
		printIndent(builder, indent, level)
	}
	switch node := node.(type) {
	case Value:
		builder.WriteString(string(node))
	case Array:
		printArray(node, builder, indent, level)
	case Object:
		printObject(node, builder, indent, level)
	}
}

func printArray(node Array, builder *strings.Builder, indent string, level int) {
	builder.WriteString("[")
	nodes := []Node(node)
	if len(nodes) > 0 {
		builder.WriteString("\n")
	}
	for i, child := range nodes {
		print(child, builder, indent, level+1, true)
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

func printObject(node Object, builder *strings.Builder, indent string, level int) {
	builder.WriteString("{")

	members := []Member(node)
	if len(members) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range members {
		printIndent(builder, indent, level+1)
		builder.WriteString(child.Name)
		builder.WriteString(": ")
		print(child.Value, builder, indent, level+1, false)
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
