package jsonish

type Node interface {
	Format() string
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
