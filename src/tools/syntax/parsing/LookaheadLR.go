package parsing

import (
	set "github.com/emirpasic/gods/sets/treeset"
	stack "github.com/emirpasic/gods/stacks/linkedliststack"
	"log"
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
func (lalr *LookaheadLR) BuildPropagateAndSpontaneousTable() {
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
func tryAddLookahead(lalr *LookaheadLR, unpropagated *stack.Stack, fromState *State, byItem *Item, lookaheads ...string) {
	dotRight := byItem.DotRight()
	peer := lalr.MakeLR0(byItem.prod, byItem.dot+1)
	targetState := fromState.GotoTable[dotRight]

	for _, lookahead := range lookaheads {
		if targetState.AddLookahead(peer, lookahead) {
			unit := [3]interface{}{targetState, peer, lookahead}
			unpropagated.Push(unit)
		}
	}
}
func (lalr *LookaheadLR) DoPropagation() {
	unpropagated := stack.New()
	//initialize spontaneous lookahead
	lalr.initialState.AddLookahead(lalr.initial, grammar.SymbolEnd)
	unit := [3]interface{}{lalr.initialState, lalr.initial, grammar.SymbolEnd}
	unpropagated.Push(unit)
	for _, state := range lalr.States {
		for _, kernel := range state.GetKernelItems() {
			for itemHash, LookaheadSet := range kernel.SpontaneousTable {
				item := lalr.ItemPool[itemHash]
				LookaheadValues := util.ToArrayString(LookaheadSet.Values())
				tryAddLookahead(lalr, unpropagated, state, item, LookaheadValues...)
			}
		}
	}
	//propagate away
	for !unpropagated.Empty() {
		val, _ := unpropagated.Pop()
		unit := val.([3]interface{})
		fromState := unit[0].(*State)
		fromItem := unit[1].(*Item)
		byLookahead := unit[2].(string)
		for _, val := range fromItem.propagateTable.Values() {
			byItem := val.(*Item)
			tryAddLookahead(lalr, unpropagated, fromState, byItem, byLookahead)
		}
	}
}
func (lalr *LookaheadLR) BuildParsingActionTable() {
	for _, state := range lalr.States {

		state.ClosureWithLookahead(lalr)
		state.ParsingActionTable = make(map[string][2]interface{})

		for itemHash, lookaheadSet := range state.LookaheadTable {
			item := lalr.ItemPool[itemHash]
			lookaheadValues := util.ToArrayString(lookaheadSet.Values())
			dotRight := item.DotRight()
			for _, lookahead := range lookaheadValues {
				if dotRight != "" {
					//shift
					if lalr.gram.IsTerminal[dotRight] {
						action := [2]interface{}{"shift"}
						if oldAction, exist := state.ParsingActionTable[dotRight]; exist {
							actName := oldAction[0].(string)
							switch actName {
							case "reduce":
								reduceProd := oldAction[1].(*grammar.Production)
								log.Printf("Warning Conflicting(S/R),perfer shift: shift %v / reduce %v", dotRight, reduceProd)
								state.ParsingActionTable[dotRight] = action
							}
						} else {
							state.ParsingActionTable[dotRight] = action
						}
					}
				} else {
					//accept
					if item.prod.Head == grammar.SymbolStart && lookahead == grammar.SymbolEnd {
						action := [2]interface{}{"accept"}
						state.ParsingActionTable[lookahead] = action
					} else //reduce
					{
						action := [2]interface{}{"reduce", item.prod}
						if oldAction, exist := state.ParsingActionTable[lookahead]; exist {
							actName := oldAction[0].(string)
							switch actName {
							case "reduce":
								reduceProd := oldAction[1].(*grammar.Production)
								log.Fatalf("Error Conflicting(R/R): reduce %v / reduce %v", reduceProd, item.prod)
							case "shift":
								log.Printf("Warning Conflicting(S/R),perfer shift: shift %v / reduce %v", lookahead, item.prod)
							}
						} else {
							state.ParsingActionTable[lookahead] = action
						}
					}
				}
			}
		}
	}
}
