package parsing

import (
	"encoding/xml"
	set "github.com/emirpasic/gods/sets/treeset"
	stack "github.com/emirpasic/gods/stacks/linkedliststack"
	"log"
	"os"
	"strconv"
	"strings"
	"tools/syntax/grammar"
	"tools/util"
)

type LookaheadLR struct {
	gram     *grammar.Grammar
	initial  *Item
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
func (lalr *LookaheadLR) visitItemWithLookahead(unchecked *stack.Stack, visited *set.Set, item *Item, lookaheadValues ...string) {
	dotRight := item.DotRight()
	if dotRight != "" && !lalr.gram.IsTerminal[dotRight] {
		for _, lookahead := range lookaheadValues {
			fst, _ := lalr.gram.CalcFst(append(item.DotRightTailingNodes(), lookahead)...)
			for _, val := range fst {
				lookahead := val.(string)
				key := dotRight + "/" + lookahead
				if !visited.Contains(key) {
					unit := [2]string{dotRight, lookahead}
					unchecked.Push(unit)
					visited.Add(key)
				}
			}
		}
	}
}
func (lalr *LookaheadLR) ClosureWithLookahead(state *State) {
	unchecked := stack.New()
	visited := set.NewWithStringComparator()
	for itemHash, lookaheadSet := range state.LookaheadTable {
		item := lalr.ItemPool[itemHash]
		lookaheadValues := util.ToArrayString(lookaheadSet.Values())
		lalr.visitItemWithLookahead(unchecked, visited, item, lookaheadValues...)
	}
	for !unchecked.Empty() {
		val, _ := unchecked.Pop()
		unit := val.([2]string)
		nonTerminal := unit[0]
		lookahead := unit[1]

		prods := lalr.gram.GetProductionsOfHead(nonTerminal)
		for _, prod := range prods {
			item := lalr.MakeLR0(prod, 0)
			if state.AddLookahead(item, lookahead) {
				lalr.visitItemWithLookahead(unchecked, visited, item, lookahead)
			}
		}
	}
}
func (lalr *LookaheadLR) visitItem(uncheckedNonTerminal *stack.Stack, visited *set.Set, item *Item) {
	dotRight := item.DotRight()
	if dotRight != "" && !lalr.gram.IsTerminal[dotRight] && !visited.Contains(dotRight) {
		uncheckedNonTerminal.Push(dotRight)
		visited.Add(dotRight)
	}
}
func (lalr *LookaheadLR) Closure(state *State) {
	uncheckedNonTerminal := stack.New()
	visited := set.NewWithStringComparator()
	for _, item := range state.GetItems() {
		lalr.visitItem(uncheckedNonTerminal, visited, item)
	}

	for !uncheckedNonTerminal.Empty() {
		val, _ := uncheckedNonTerminal.Pop()
		nonTerminal := val.(string)
		prods := lalr.gram.GetProductionsOfHead(nonTerminal)
		for _, prod := range prods {
			item := lalr.MakeLR0(prod, 0)
			if !state.Items.Contains(item) {
				state.Items.Add(item)
				lalr.visitItem(uncheckedNonTerminal, visited, item)
			}
		}
	}
}
func (lalr *LookaheadLR) AddState(state *State) (result *State, added bool) {
	hash := state.HashString()
	var exist bool
	if _, exist = lalr.States[hash]; !exist {
		state.Id = len(lalr.States)
		lalr.States[hash] = state
	}
	return lalr.States[hash], !exist
}
func (lalr *LookaheadLR) GroupGOTOTable(state *State) map[string]*State {
	gotoTable := make(map[string]*State)
	for _, item := range state.GetItems() {
		dotRight := item.DotRight()
		if dotRight != "" {
			if _, ok := gotoTable[dotRight]; !ok {
				gotoTable[dotRight] = NewState()
			}
			peer := lalr.MakeLR0(item.prod, item.dot+1)
			gotoTable[dotRight].Items.Add(peer)
		}
	}
	return gotoTable
}
func (lalr *LookaheadLR) BuildCanonicalCollection() {
	lalr.initial = lalr.MakeLR0(lalr.gram.Productions[0], 0)

	lalr.initialState = NewState(lalr.initial)
	lalr.AddState(lalr.initialState)

	uncheckedState := stack.New()
	uncheckedState.Push(lalr.initialState)
	for !uncheckedState.Empty() {
		val, _ := uncheckedState.Pop()
		state := val.(*State)
		lalr.Closure(state)
		gotoTable := lalr.GroupGOTOTable(state)

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
			lalr.ClosureWithLookahead(dummyState)

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

		lalr.ClosureWithLookahead(state)
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
func (lalr *LookaheadLR) ToXml() {
	xmll := xmlLALR{}

	for _, prod := range lalr.gram.Productions {
		xmlp := xmlProduction{Head: prod.Head, Nodes: strings.Join(prod.Nodes, "|"), Script: prod.Script, Len: len(prod.Nodes), Id: prod.Id}
		xmll.Productions = append(xmll.Productions, xmlp)
	}

	for _, state := range lalr.States {
		xmls := xmlState{Id: state.Id}
		for on, target := range state.GotoTable {
			xmlgt := &xmlGoto{On: on, State: target.Id}
			xmls.Gotos = append(xmls.Gotos, xmlgt)
		}
		for on, action := range state.ParsingActionTable {
			actName := action[0].(string)
			if actName == "reduce" {
				reduceProd := action[1].(*grammar.Production)
				actName += strconv.Itoa(reduceProd.Id)
			}
			xmlact := &xmlAction{On: on, Do: actName}
			xmls.Actions = append(xmls.Actions, xmlact)
		}
		xmll.States = append(xmll.States, xmls)
	}

	var f *os.File

	// Check if thet file exists, err != nil if the file does not exist
	_, err := os.Stat("lalr.xml")
	if err != nil {
		// if the file doesn't exist, open it with write and create flags
		f, err = os.OpenFile("lalr.xml", os.O_WRONLY|os.O_CREATE, 0666)
	} else {
		// if the file does exist, open it with append and write flags
		f, err = os.OpenFile("lalr.xml", os.O_WRONLY|os.O_TRUNC, 0666)
	}
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := xml.NewEncoder(f)

	err = e.Encode(xmll)
	if err != nil {
		panic(err)
	}
}
