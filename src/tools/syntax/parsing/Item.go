package parsing

import (
	"fmt"
	set "github.com/emirpasic/gods/sets/treeset"
	"strconv"
	"strings"
	"tools/syntax/grammar"
)

type Item struct {
	prod *grammar.Production
	dot  int

	SpontaneousTable map[string]*set.Set
	propagateTable   *set.Set
}

func newItem(prod *grammar.Production, dot int) *Item {
	return &Item{prod, dot, nil, nil}
}
func (item *Item) String() string {
	return fmt.Sprintf("%v:%v /%v", item.prod.Head, item.prod.Nodes, item.dot)
}
func (item *Item) DotRight() string {
	if item.prod.IsNull() {
		return ""
	}
	if item.dot < len(item.prod.Nodes) {
		return item.prod.Nodes[item.dot]
	}
	return ""
}
func (item *Item) DotRightNodes(dot int) []string {
	var nodes []string

	if item.prod.IsNull() {
		return nodes
	}
	if dot < len(item.prod.Nodes) {
		nodes = item.prod.Nodes[item.dot:]
	}
	return nodes
}
func (item *Item) IsKernel() bool {
	if item.prod.IsNull() {
		return false
	}
	if item.prod.IsInitial() {
		return true
	}
	return item.dot != 0
}
func HashString(prod *grammar.Production, dot int) string {
	hash := prod.Head + "/" + strings.Join(prod.Nodes, "+") + "/" + strconv.Itoa(dot)
	return hash
}
func (item *Item) HashString() string {
	return HashString(item.prod, item.dot)
}
func byHash(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	item1 := a.(*Item)
	item2 := b.(*Item)
	return strings.Compare(item1.HashString(), item2.HashString())
}
