package parsing

import (
	"fmt"
	set "github.com/emirpasic/gods/sets/treeset"
	"strconv"
	"strings"
	"tools/syntax/grammar"
)

type State struct {
	Items     *set.Set
	GotoTable map[string]*State

	LookaheadTable     map[string]*set.Set
	ParsingActionTable map[string][2]interface{}
	Id                 int
}

func (state *State) String() string {
	str := ""
	for _, kernel := range state.GetKernelItems() {
		str += kernel.String() + "\n"
	}
	return str
}
func NewState(values ...interface{}) *State {
	s := State{}
	s.Items = set.NewWith(byHash)
	if len(values) > 0 {
		s.Items.Add(values...)
	}
	s.LookaheadTable = make(map[string]*set.Set)
	s.GotoTable = make(map[string]*State)
	return &s
}
func (state *State) AddLookahead(item *Item, lookahead string) (added bool) {
	hash := item.HashString()
	if val, exist := state.LookaheadTable[hash]; !exist {
		state.LookaheadTable[hash] = set.NewWithStringComparator(lookahead)
		return true
	} else {
		if !val.Contains(lookahead) {
			val.Add(lookahead)
			return true
		}
		return false
	}
}
func (state *State) GetKernelItems() []*Item {
	var items []*Item
	for it := state.Items.Iterator(); it.Next(); {
		item := it.Value().(*Item)
		//only use kernels
		if item.IsKernel() {
			items = append(items, item)
		}
	}
	return items
}
func (state *State) GetItems() []*Item {
	var items []*Item
	for it := state.Items.Iterator(); it.Next(); {
		val := it.Value().(*Item)
		items = append(items, val)
	}
	return items
}
func (state *State) HashString() string {
	var items []string
	for _, kernel := range state.GetKernelItems() {
		items = append(items, kernel.HashString())
	}
	return strings.Join(items, "@")
}

type Item struct {
	prod *grammar.Production
	dot  int

	SpontaneousTable map[string]*set.Set
	PropagateTable   *set.Set
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
func (item *Item) DotRightTailingNodes() []string {
	var nodes []string

	if item.prod.IsNull() {
		return nodes
	}
	if (item.dot + 1) < len(item.prod.Nodes) {
		nodes = item.prod.Nodes[(item.dot + 1):]
	}
	return nodes
}
func (item *Item) IsKernel() bool {
	if item.prod.IsNull() {
		return false
	}
	if item.prod.Head == grammar.SymbolStart {
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
