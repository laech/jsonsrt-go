package jsonish

import (
	"sort"
)

type Node interface {
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
		return unquote(members[i].Name) < unquote(members[j].Name)
	})
}

func (node Value) SortByValue(string) {}

func (node Array) SortByValue(name string) {
	nodes := []Node(node)
	for i := range nodes {
		nodes[i].SortByValue(name)
	}
	sort.Slice(nodes, func(i, j int) bool {
		a, aOk := nodes[i].(Object)
		b, bOk := nodes[j].(Object)
		if !aOk || !bOk {
			return false
		}
		x := a.findValue(name)
		y := b.findValue(name)
		if x != nil && b != nil {
			return unquote(string(*x)) < unquote(string(*y))
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

func unquote(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	} else {
		return str
	}
}

func (node Object) findValue(name string) *Value {
	name = `"` + name + `"`
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
