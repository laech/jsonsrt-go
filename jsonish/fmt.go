package jsonish

import (
	"fmt"
	"io"
	"reflect"
)

func print(node Node, writer io.Writer, indent []byte, level int, applyInitalIndent bool) error {
	if applyInitalIndent {
		if err := printIndent(writer, indent, level); err != nil {
			return err
		}
	}
	switch t := node.(type) {
	case Value:
		if _, err := writer.Write(t.Value); err != nil {
			return err
		}
	case Array:
		if err := printArray(t, writer, indent, level); err != nil {
			return nil
		}
	case Object:
		if err := printObject(t, writer, indent, level); err != nil {
			return nil
		}
	default:
		panic(fmt.Sprintf("Unknown type: %s", reflect.TypeOf(t)))
	}
	return nil
}

func printArray(arr Array, writer io.Writer, indent []byte, level int) error {
	if _, err := writer.Write([]byte("[")); err != nil {
		return err
	}

	if len(arr.Value) > 0 {
		if _, err := writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	for i, child := range arr.Value {
		if err := print(child, writer, indent, level+1, true); err != nil {
			return err
		}
		if i < len(arr.Value)-1 {
			if _, err := writer.Write([]byte(",")); err != nil {
				return err
			}
		}
	}

	if arr.TrailingSep {
		if _, err := writer.Write([]byte(",")); err != nil {
			return err
		}
	}

	if len(arr.Value) > 0 {
		if _, err := writer.Write([]byte("\n")); err != nil {
			return err
		}
		if err := printIndent(writer, indent, level); err != nil {
			return err
		}
	}

	if _, err := writer.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}

func printObject(obj Object, writer io.Writer, indent []byte, level int) error {
	if _, err := writer.Write([]byte("{")); err != nil {
		return err
	}

	if len(obj.Value) > 0 {
		if _, err := writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	for i, child := range obj.Value {
		if err := printIndent(writer, indent, level+1); err != nil {
			return err
		}
		if _, err := writer.Write(child.Name); err != nil {
			return err
		}
		if _, err := writer.Write([]byte(": ")); err != nil {
			return err
		}
		if err := print(child.Value, writer, indent, level+1, false); err != nil {
			return err
		}
		if i < len(obj.Value)-1 {
			if _, err := writer.Write([]byte(",")); err != nil {
				return err
			}
		}
	}

	if obj.TrailingSep {
		if _, err := writer.Write([]byte(",")); err != nil {
			return err
		}
	}

	if len(obj.Value) > 0 {
		if _, err := writer.Write([]byte("\n")); err != nil {
			return err
		}
		if err := printIndent(writer, indent, level); err != nil {
			return err
		}
	}

	if _, err := writer.Write([]byte("}")); err != nil {
		return err
	}

	return nil
}

func printIndent(writer io.Writer, indent []byte, level int) error {
	for i := 0; i < level; i++ {
		if _, err := writer.Write(indent); err != nil {
			return err
		}
	}
	return nil
}
