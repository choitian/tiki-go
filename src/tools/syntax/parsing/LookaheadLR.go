package parsing

import (
	set "github.com/emirpasic/gods/sets/treeset"
	stack "github.com/emirpasic/gods/stacks/linkedliststack"
	"tools/syntax/grammar"
	"tools/util"
)

type LookaheadLR struct {
	gram     *grammar.Grammar
	initial  *Item
	accept   *Item
	ItemPool map[string]*Item

	initialState *State
	States       map[string]*State
}

func NewLookaheadLR(gram *grammar.Grammar) *LookaheadLR {
	lalr := LookaheadLR{}
	lalr.gram = gram
	lalr.ItemPool = make(map[string]*Item)
	lalr.States = make(map[string]*State)
	return &lalr
}
func (lalr *LookaheadLR) MakeLR0(prod *grammar.Production, dot int) *Item {
	newed := newItem(prod, dot)
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
func (lalr *LookaheadLR) BuildCanonicalCollection() {
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

		for onSymbol, targetState := range gotoTable {
			result, added := lalr.AddState(targetState)
			if added {
				uncheckedState.Push(result)
			}
			state.GotoTable[onSymbol] = result
		}
	}
}
func (lalr *LookaheadLR) BuildPropagateAndSpontaneouTable() {
	for _, kernel := range lalr.ItemPool {
		if kernel.IsKernel() {
			kernel.SpontaneousTable = make(map[string]*set.Set)
			kernel.propagateTable = set.NewWith(byHash)

			dummyState := NewState(kernel)
			dummyState.AddLookahead(kernel, grammar.SymbolEnd)
			dummyState.ClosureWithLookahead(lalr)

			for itemHash, lookaheadSet := range dummyState.LookaheadTable {
				item := lalr.ItemPool[itemHash]
				//skip item who's dot at end
				if item.DotRight() != "" {
					lookaheadValues := util.ToArrayString(lookaheadSet.Values())
					for _, lookahead := range lookaheadValues {
						if lookahead != grammar.SymbolEnd {
							hash := item.HashString()
							if lookaheadSet, exist := kernel.SpontaneousTable[hash]; !exist {
								kernel.SpontaneousTable[hash] = set.NewWithStringComparator(lookahead)
							} else {
								lookaheadSet.Add(lookahead)
							}
						} else {
							kernel.propagateTable.Add(item)
						}
					}
				}
			}
		}
	}
}
func (lalr *LookaheadLR) doPropagation() {

}
