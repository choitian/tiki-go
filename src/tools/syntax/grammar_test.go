package syntax

import (
	"log"
	"os"
	"strings"
	"testing"
	"tools/util"
)

func TestGrammerDescription(t *testing.T) {
	//gram :=NewGrammar("test/re_grammar.txt")
	gram := NewGrammar("test/dnf.txt")
	/*
		for _, p := range gram.Productions {
			log.Printf("%v\n", &p)
		}
	*/
	for k, v := range gram.FST {
		fst := util.StringBoolMapKeys(v)
		//t.Logf("%v: %v\n", k, fst)
		if k == "exp" {
			if strings.Join(fst, " ") != "DEC FALSE FLOATING ID INC INTEGER LPAREN NEW NOT NULL STRING SUB TRUE" {
				t.Fatalf("Fst of 'exp' is wrong!")
			}
		}
		if k == "ini_exp" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR COMMA DEC FALSE FLOAT FLOATING ID INC INT INTEGER LPAREN NEW NOT NULL STATIC STRING SUB TRUE VOID" {
				t.Fatalf("Fst of 'ini_exp' is wrong!")
			}
		}
		if k == "postfix_exp" {
			if strings.Join(fst, " ") != "FALSE FLOATING ID INTEGER LPAREN NULL STRING TRUE" {
				t.Fatalf("Fst of 'postfix_exp' is wrong!")
			}
		}
		if k == "method_definition" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR FLOAT ID INT STATIC VOID" {
				t.Fatalf("Fst of 'method_definition' is wrong!")
			}
		}
	}
	/*
		for k, v := range gram.IsNullable {
			if v {
				t.Logf("%v: %v\n", k, v)
			}
		}
	*/
}
func TestMain(m *testing.M) {
	log.SetOutput(os.Stderr)

	runTests := m.Run()
	os.Exit(runTests)
}
