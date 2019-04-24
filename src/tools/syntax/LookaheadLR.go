package syntax

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	stack "github.com/emirpasic/gods/stacks/linkedliststack"
	"strconv"
	"strings"
)

type Item struct {
	prod *Production
	dot  int
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
func HashString(prod *Production, dot int) string {
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

type State struct {
	Items *treeset.Set

	GotoTable map[string]*State
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
	s.Items = treeset.NewWith(byHash)
	if len(values) > 0 {
		s.Items.Add(values...)
	}
	return &s
}
func (state *State) GetKernelItems() []*Item {
	var items []*Item
	for it := state.Items.Iterator(); it.Next(); {
		val := it.Value().(*Item)
		//only use kernels
		if (val.dot != 0) || val.prod.IsInitial() {
			items = append(items, val)
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
func (state *State) Closure(lalr *LookaheadLR) {
	uncheckedNonTerminal := stack.New()
	visited := treeset.NewWithStringComparator()
	for _, item := range state.GetItems() {
		dotRight := item.DotRight()
		if dotRight != "" && !lalr.gram.IsTerminal[dotRight] && !visited.Contains(dotRight) {
			uncheckedNonTerminal.Push(dotRight)
			visited.Add(dotRight)
		}
	}

	for !uncheckedNonTerminal.Empty() {
		val, _ := uncheckedNonTerminal.Pop()
		nonTerminal := val.(string)
		prods := lalr.gram.GetProductionsOfHead(nonTerminal)
		for _, prod := range prods {
			item := lalr.MakeLR0(prod, 0)
			if !state.Items.Contains(item) {
				state.Items.Add(item)

				dotRight := item.DotRight()
				if dotRight != "" && !lalr.gram.IsTerminal[dotRight] && !visited.Contains(dotRight) {
					uncheckedNonTerminal.Push(dotRight)
					visited.Add(dotRight)
				}
			}
		}
	}
}
func (state *State) GroupGOTOTable(lalr *LookaheadLR) map[string]*State {
	gotoTable := make(map[string]*State)
	for _, item := range state.GetItems() {
		dotRight := item.DotRight()
		if dotRight != "" {
			targetState, ok := gotoTable[dotRight]
			if !ok {
				targetState = NewState()
				gotoTable[dotRight] = targetState
			}
			peer := lalr.MakeLR0(item.prod, item.dot+1)
			targetState.Items.Add(peer)
		}
	}
	return gotoTable
}

type LookaheadLR struct {
	gram     *Grammar
	initial  *Item
	accept   *Item
	ItemPool map[string]*Item

	initialState *State
	States       map[string]*State
}

func NewLookaheadLR(gram *Grammar) *LookaheadLR {
	lalr := LookaheadLR{}
	lalr.gram = gram
	lalr.ItemPool = make(map[string]*Item)
	lalr.States = make(map[string]*State)
	return &lalr
}
func (lalr *LookaheadLR) MakeLR0(prod *Production, dot int) *Item {
	newed := &Item{prod, dot}
	hash := newed.HashString()
	if _, exist := lalr.ItemPool[hash]; !exist {
		lalr.ItemPool[hash] = newed
	}
	return lalr.ItemPool[hash]
}
func (lalr *LookaheadLR) AddState(state *State) (result *State, added bool) {
	hash := state.HashString()
	var exist bool
	if _, exist = lalr.States[hash]; !exist {
		lalr.States[hash] = state
	}
	return lalr.States[hash], !exist
}
func (lalr *LookaheadLR) ConstructCanonicalCollection() {
	lalr.initial = lalr.MakeLR0(lalr.gram.Productions[0], 0)
	lalr.accept = lalr.MakeLR0(lalr.gram.Productions[0], 1)

	lalr.initialState = NewState(lalr.initial)
	lalr.AddState(lalr.initialState)

	uncheckedState := stack.New()
	uncheckedState.Push(lalr.initialState)
	for !uncheckedState.Empty() {
		val, _ := uncheckedState.Pop()
		state := val.(*State)
		state.Closure(lalr)
		gotoTable := state.GroupGOTOTable(lalr)

		state.GotoTable = make(map[string]*State)
		for onSymbol, targetState := range gotoTable {
			result, added := lalr.AddState(targetState)
			if added {
				uncheckedState.Push(result)
			}
			state.GotoTable[onSymbol] = result
		}
	}
}
