package parsing

import (
	set "github.com/emirpasic/gods/sets/treeset"
	"strings"
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
