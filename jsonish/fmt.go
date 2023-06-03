package jsonish

import (
	"fmt"
	"reflect"
	"strings"
)

func print(node Node, builder *strings.Builder, indent string, level int, applyInitalIndent bool) {
	if applyInitalIndent {
		printIndent(builder, indent, level)
	}
	switch t := node.(type) {
	case Value:
		builder.WriteString(t.Value)
	case Array:
		printArray(t, builder, indent, level)
	case Object:
		printObject(t, builder, indent, level)
	default:
		panic(fmt.Sprintf("Unknown type: %s", reflect.TypeOf(t)))
	}
}

func printArray(arr Array, builder *strings.Builder, indent string, level int) {
	builder.WriteString("[")
	if len(arr.Value) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range arr.Value {
		print(child, builder, indent, level+1, true)
		if i < len(arr.Value)-1 {
			builder.WriteString(",")
		}
	}

	if arr.TrailingSep {
		builder.WriteString(",")
	}

	if len(arr.Value) > 0 {
		builder.WriteString("\n")
		printIndent(builder, indent, level)
	}

	builder.WriteString("]")
}

func printObject(obj Object, builder *strings.Builder, indent string, level int) {
	builder.WriteString("{")

	if len(obj.Value) > 0 {
		builder.WriteString("\n")
	}

	for i, child := range obj.Value {
		printIndent(builder, indent, level+1)
		builder.WriteString(child.Name)
		builder.WriteString(": ")
		print(child.Value, builder, indent, level+1, false)
		if i < len(obj.Value)-1 {
			builder.WriteString(",")
		}
	}

	if obj.TrailingSep {
		builder.WriteString(",")
	}

	if len(obj.Value) > 0 {
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
