package syntax

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"io/ioutil"
	"regexp"
	"strings"
	"tools/util"
)

const symbolEnd string = "__END__"
const symbolNull string = "__NULL__"
const symbolStart string = "__START__"

type Grammar struct {
	Productions []*Production

	FST        map[string]*treeset.Set
	IsTerminal map[string]bool
	Nullable   *treeset.Set
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
	return len(prod.Nodes) == 1 && prod.Nodes[0] == symbolNull
}
func (prod *Production) IsInitial() bool {
	return prod.Head == symbolStart
}
func (gram *Grammar) GetFst(symbol string) []string {
	var keys []string
	for _, v := range gram.FST[symbol].Values() {
		keys = append(keys, v.(string))
	}
	return keys
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
	gram.FST = make(map[string]*treeset.Set)
	gram.IsTerminal = make(map[string]bool)
	gram.Nullable = treeset.NewWithStringComparator()
	//initialize IsTerminal & Nullable
	for _, prod := range gram.Productions {
		head := prod.Head
		gram.IsTerminal[head] = false
		gram.FST[head] = treeset.NewWithStringComparator()
	}
	//if not exist in IsTerminal as not being a head,then Is Terminal.
	for _, prod := range gram.Productions {
		for _, symbol := range prod.Nodes {
			if _, ok := gram.IsTerminal[symbol]; !ok {
				gram.IsTerminal[symbol] = true
				gram.FST[symbol] = treeset.NewWithStringComparator(symbol)
			}
		}
	}
	//iterate util nothing changed
	for nothingChanged := false; !nothingChanged; {
		nothingChanged = true
		for _, prod := range gram.Productions {
			head := prod.Head
			if !prod.IsNull() {
				for i, symbol := range prod.Nodes {
					//add FST of symbol to head's
					for it := gram.FST[symbol].Iterator(); it.Next(); {
						if !gram.FST[head].Contains(it.Value()) {
							gram.FST[head].Add(it.Value())
							nothingChanged = false
						}
					}
					if !gram.Nullable.Contains(symbol) {
						break
					}

					if i == len(prod.Nodes)-1 {
						if !gram.Nullable.Contains(head) {
							gram.Nullable.Add(head)
							nothingChanged = false
						}
					}
				}
			} else {
				if !gram.Nullable.Contains(head) {
					gram.Nullable.Add(head)
					nothingChanged = false
				}
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
		gm.Productions = []*Production{&Production{symbolStart, []string{productions[0].Head}, "{}"}}
		gm.Productions = append(gm.Productions, productions...)
	}
	gm.computeAttributes()
	return &gm
}
