package parsing

import (
	set "github.com/emirpasic/gods/sets/treeset"
	stack "github.com/emirpasic/gods/stacks/linkedliststack"
	"strings"
	"tools/util"
)

type State struct {
	Items          *set.Set
	LookaheadTable map[string]*set.Set

	GotoTable          map[string]*State
	ParsingActionTable map[string][2]interface{}
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

func visitItem(lalr *LookaheadLR, uncheckedNonTerminal *stack.Stack, visited *set.Set, item *Item) {
	dotRight := item.DotRight()
	if dotRight != "" && !lalr.gram.IsTerminal[dotRight] && !visited.Contains(dotRight) {
		uncheckedNonTerminal.Push(dotRight)
		visited.Add(dotRight)
	}
}
func (state *State) Closure(lalr *LookaheadLR) {
	uncheckedNonTerminal := stack.New()
	visited := set.NewWithStringComparator()
	for _, item := range state.GetItems() {
		visitItem(lalr, uncheckedNonTerminal, visited, item)
	}

	for !uncheckedNonTerminal.Empty() {
		val, _ := uncheckedNonTerminal.Pop()
		nonTerminal := val.(string)
		prods := lalr.gram.GetProductionsOfHead(nonTerminal)
		for _, prod := range prods {
			item := lalr.MakeLR0(prod, 0)
			if !state.Items.Contains(item) {
				state.Items.Add(item)
				visitItem(lalr, uncheckedNonTerminal, visited, item)
			}
		}
	}
}
func visitItemWithLookahead(lalr *LookaheadLR, unchecked *stack.Stack, visited *set.Set, item *Item, lookaheadValues ...string) {
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
func (state *State) ClosureWithLookahead(lalr *LookaheadLR) {
	unchecked := stack.New()
	visited := set.NewWithStringComparator()
	for itemHash, lookaheadSet := range state.LookaheadTable {
		item := lalr.ItemPool[itemHash]
		lookaheadValues := util.ToArrayString(lookaheadSet.Values())
		visitItemWithLookahead(lalr, unchecked, visited, item, lookaheadValues...)
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
				visitItemWithLookahead(lalr, unchecked, visited, item, lookahead)
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
