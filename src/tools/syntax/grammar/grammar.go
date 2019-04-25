package grammar

import (
	"fmt"
	set "github.com/emirpasic/gods/sets/treeset"
	"io/ioutil"
	"regexp"
	"strings"
	"tools/util"
)

const SymbolEnd string = "__END__"
const SymbolNull string = "__NULL__"
const SymbolStart string = "__START__"

type Grammar struct {
	Productions []*Production

	FST        map[string]*set.Set
	IsTerminal map[string]bool
	Nullable   *set.Set
}

type Production struct {
	Head   string
	Nodes  []string
	Script string
}

func (prod *Production) String() string {
	return fmt.Sprintf("%v:%v %v", prod.Head, prod.Nodes, prod.Script)
}
func (prod *Production) IsNull() bool {
	return len(prod.Nodes) == 1 && prod.Nodes[0] == SymbolNull
}
func (prod *Production) IsInitial() bool {
	return prod.Head == SymbolStart
}
func (gram *Grammar) GetFst(symbol string) []string {
	return util.ToArrayString(gram.FST[symbol].Values())
}
func (gram *Grammar) CalcFst(symbols ...string) (firstSet []interface{}, nullable bool) {
	fst := set.NewWithStringComparator()
	nullable = false
	for i, symbol := range symbols {
		//merge it,skip SymbolNull
		if symbol != SymbolNull {
			vs := gram.FST[symbol].Values()
			fst.Add(vs...)
		}

		if !gram.Nullable.Contains(symbol) {
			break
		}
		if i == len(symbols)-1 {
			nullable = true
		}
	}
	firstSet = fst.Values()
	return
}
func (gram *Grammar) GetProductionsOfHead(head string) []*Production {
	var prods []*Production
	for _, prod := range gram.Productions {
		if prod.Head == head {
			prods = append(prods, prod)
		}
	}
	return prods
}
func (gram *Grammar) computeAttributes() {
	//initialize maps
	gram.FST = make(map[string]*set.Set)
	gram.IsTerminal = make(map[string]bool)
	gram.Nullable = set.NewWithStringComparator()

	//initialize special symbols
	gram.Nullable.Add(SymbolNull)
	gram.IsTerminal[SymbolEnd] = true
	gram.FST[SymbolEnd] = set.NewWithStringComparator(SymbolEnd)
	//initialize IsTerminal & Nullable
	for _, prod := range gram.Productions {
		head := prod.Head
		gram.IsTerminal[head] = false
		gram.FST[head] = set.NewWithStringComparator()
	}
	//if not exist in IsTerminal as not being a head,then Is Terminal.
	for _, prod := range gram.Productions {
		for _, symbol := range prod.Nodes {
			if _, ok := gram.IsTerminal[symbol]; !ok {
				gram.IsTerminal[symbol] = true
				gram.FST[symbol] = set.NewWithStringComparator(symbol)
			}
		}
	}
	//iterate util nothing changed
	for nothingChanged := false; !nothingChanged; {
		nothingChanged = true
		for _, prod := range gram.Productions {
			head := prod.Head
			fst, nullable := gram.CalcFst(prod.Nodes...)
			if nullable && !gram.Nullable.Contains(head) {
				gram.Nullable.Add(head)
				nothingChanged = false
			}

			oldSize := gram.FST[head].Size()
			gram.FST[head].Add(fst...)
			if gram.FST[head].Size() > oldSize {
				nothingChanged = false
			}
		}
	}
}
func NewGrammar(path string) *Grammar {
	var gm Grammar
	//read it as text
	byteValue, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	content := string(byteValue)

	//load productions
	var productions []*Production
	rSymbol := regexp.MustCompile(`(\s*(?P<Head>\w+)\s*)\:(?P<prods>((\s*\w+\s*)+(\s*\{.*\}\s*)?\|?)+)(\s*;)`)
	for _, mSymbol := range rSymbol.FindAllStringSubmatch(content, -1) {
		nm := util.MatchNamedMap(rSymbol, mSymbol)
		head := strings.TrimSpace(nm["Head"])
		prods := strings.TrimSpace(nm["prods"])
		rProd := regexp.MustCompile(`(?P<prod>(\s*\w+\s*)+)(\s*(?P<Script>\{.*\})\s*)?\|?`)
		for _, mProd := range rProd.FindAllStringSubmatch(prods, -1) {
			nm := util.MatchNamedMap(rProd, mProd)
			prod := strings.TrimSpace(nm["prod"])
			prodNodes := regexp.MustCompile(`\s+`).Split(prod, -1)
			script := strings.TrimSpace(nm["Script"])
			if script == "" {
				script = "{}"
			}
			productions = append(productions, &Production{head, prodNodes, script})
		}
	}
	if len(productions) != 0 {
		gm.Productions = []*Production{&Production{SymbolStart, []string{productions[0].Head}, "{}"}}
		gm.Productions = append(gm.Productions, productions...)
	}
	gm.computeAttributes()
	return &gm
}
