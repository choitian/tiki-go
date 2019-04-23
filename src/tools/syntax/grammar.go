package syntax

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"tools/util"
)

const symbolEnd string = "__END__"
const symbolNull string = "__NULL__"
const symbolStart string = "__START__"

type StringSet map[string]bool
type Grammar struct {
	Productions []Production

	FST        map[string]StringSet
	IsTerminal map[string]bool
	Nullable   StringSet
}

type Production struct {
	Head   string
	Nodes  []string
	Script string
}

func (prod *Production) String() string {
	return fmt.Sprintf("%v --> %v %v", prod.Head, prod.Nodes, prod.Script)
}
func (prod *Production) IsNull() bool {
	return len(prod.Nodes) == 1 && prod.Nodes[0] == symbolNull
}
func (gram *Grammar) computeAttributes() {
	//initialize maps
	gram.FST = make(map[string]StringSet)
	gram.IsTerminal = make(map[string]bool)
	gram.Nullable = make(StringSet)
	//initialize IsTerminal & Nullable
	for _, prod := range gram.Productions {
		head := prod.Head
		gram.IsTerminal[head] = false
		gram.FST[head] = StringSet{}
	}
	//if not exist in IsTerminal as not being a head,then Is Terminal.
	for _, prod := range gram.Productions {
		for _, symbol := range prod.Nodes {
			if _, ok := gram.IsTerminal[symbol]; !ok {
				gram.IsTerminal[symbol] = true
				gram.FST[symbol] = StringSet{symbol: true}
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
					if fst, ok := gram.FST[symbol]; ok {
						for val, _ := range fst {
							if _, included := gram.FST[head][val]; !included {
								gram.FST[head][val] = true
								nothingChanged = false
							}
						}
					}

					if !gram.Nullable[symbol] {
						break
					}

					if i == len(prod.Nodes)-1 {
						if _, included := gram.Nullable[head]; !included {
							gram.Nullable[head] = true
							nothingChanged = false
						}
					}
				}
			} else {
				if _, included := gram.Nullable[head]; !included {
					gram.Nullable[head] = true
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
	var productions []Production
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
			productions = append(productions, Production{head, prodNodes, script})
		}
	}
	gm.Productions = productions
	gm.computeAttributes()
	return &gm
}
