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
			log.Printf("%val\n", &p)
		}
	*/
	for key, val := range gram.FST {
		fst := util.StringBoolMapKeys(val)
		//t.Logf("%val: %val\n", key, fst)
		if key == "exp" {
			if strings.Join(fst, " ") != "DEC FALSE FLOATING ID INC INTEGER LPAREN NEW NOT NULL STRING SUB TRUE" {
				t.Fatalf("Fst of 'exp' is wrong!")
			}
		}
		if key == "ini_exp" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR COMMA DEC FALSE FLOAT FLOATING ID INC INT INTEGER LPAREN NEW NOT NULL STATIC STRING SUB TRUE VOID" {
				t.Fatalf("Fst of 'ini_exp' is wrong!")
			}
		}
		if key == "postfix_exp" {
			if strings.Join(fst, " ") != "FALSE FLOATING ID INTEGER LPAREN NULL STRING TRUE" {
				t.Fatalf("Fst of 'postfix_exp' is wrong!")
			}
		}
		if key == "method_definition" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR FLOAT ID INT STATIC VOID" {
				t.Fatalf("Fst of 'method_definition' is wrong!")
			}
		}
	}
	/*
		for key, val := range gram.Nullable {
			if val {
				t.Logf("%val: %val\n", key, val)
			}
		}
	*/
}
func TestMain(m *testing.M) {
	log.SetOutput(os.Stderr)

	runTests := m.Run()
	os.Exit(runTests)
}
