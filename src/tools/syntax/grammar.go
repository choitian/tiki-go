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

type FstSet map[string]bool
type Grammar struct {
	Productions []Production

	FST        map[string]FstSet
	IsTerminal map[string]bool
	IsNullable map[string]bool
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
	gram.FST = make(map[string]FstSet)
	gram.IsTerminal = make(map[string]bool)
	gram.IsNullable = make(map[string]bool)
	//initialize IsTerminal & IsNullable
	for _, p := range gram.Productions {
		head := p.Head
		gram.IsTerminal[head] = false
		gram.IsNullable[head] = p.IsNull()
		gram.FST[head] = FstSet{}
	}
	//if not exist in IsTerminal as not being a head,then Is Terminal.
	for _, p := range gram.Productions {
		for _, symbol := range p.Nodes {
			if _, ok := gram.IsTerminal[symbol]; !ok {
				gram.IsTerminal[symbol] = true
				gram.FST[symbol] = FstSet{symbol: true}
			}
		}
	}
	//iterate util nothing changed
	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, p := range gram.Productions {
			head := p.Head
			oldFSTSize := 0
			if fst, ok := gram.FST[head]; ok {
				oldFSTSize = len(fst)
			}
			oldNullable := gram.IsNullable[head]

			if !p.IsNull() {
				for i, symbol := range p.Nodes {
					//add FST of symbol to head's
					if fst, ok := gram.FST[symbol]; ok {
						for val, _ := range fst {
							gram.FST[head][val] = true
						}
					}

					if !gram.IsNullable[symbol] {
						break
					}
					if i == len(p.Nodes)-1 {
						gram.IsNullable[head] = true
					}
				}
			}

			if (oldFSTSize != len(gram.FST[head])) || (oldNullable != gram.IsNullable[head]) {
				somethingChanged = true
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
