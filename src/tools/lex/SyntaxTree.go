package lex

const (
	Leaf string = "Leaf"
	Or   string = "Or"
	Cat  string = "Cat"
	Star string = "Star"
)

type SyntaxTree struct {
	Operator       string
	Value          byte //for Leaf
	Child0, Child1 *SyntaxTree
}

func NewLeaf(v byte) *SyntaxTree {
	var node SyntaxTree
	node.Operator = Leaf
	node.Value = v
	return &node
}
func NewOr(c0 *SyntaxTree, c1 *SyntaxTree) *SyntaxTree {
	var node SyntaxTree
	node.Operator = Or
	node.Child0 = c0
	node.Child1 = c1
	return &node
}
func NewCat(c0 *SyntaxTree, c1 *SyntaxTree) *SyntaxTree {
	var node SyntaxTree
	node.Operator = Cat
	node.Child0 = c0
	node.Child1 = c1
	return &node
}
func NewStar(c0 *SyntaxTree) *SyntaxTree {
	var node SyntaxTree
	node.Operator = Star
	node.Child0 = c0
	return &node
}
