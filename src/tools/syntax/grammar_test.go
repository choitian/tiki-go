package syntax

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestLALRConstructCanonicalCollection(t *testing.T) {
	{
		gram := NewGrammar("test/re_grammar_0.txt")
		lalr := NewLookaheadLR(gram)
		lalr.ConstructCanonicalCollection()

		if len(lalr.States) != 19 {
			t.Fatalf("state size is wrong")
		}
		kernelSum := 0
		gotoSum := 0
		for _, state := range lalr.States {
			kernelSum += len(state.GetKernelItems())
			gotoSum += len(state.GotoTable)
		}

		if kernelSum != 24 {
			t.Fatalf("kernelSum is wrong")
		}
		if gotoSum != 42 {
			t.Fatalf("gotoSum is wrong")
		}
	}
	{
		gram := NewGrammar("test/dnf.txt")
		lalr := NewLookaheadLR(gram)
		lalr.ConstructCanonicalCollection()

		if len(lalr.States) != 249 {
			t.Fatalf("state size is wrong")
		}
		kernelSum := 0
		gotoSum := 0
		for _, state := range lalr.States {
			kernelSum += len(state.GetKernelItems())
			gotoSum += len(state.GotoTable)
		}

		if kernelSum != 383 {
			t.Fatalf("kernelSum is wrong")
		}
		if gotoSum != 1465 {
			t.Fatalf("gotoSum is wrong")
		}
	}
}
func TestGrammerDescription(t *testing.T) {
	//gram :=NewGrammar("test/re_grammar.txt")
	gram := NewGrammar("test/dnf.txt")
	/*
		for _, p := range gram.Productions {
			log.Printf("%val\n", &p)
		}
	*/
	for key, _ := range gram.FST {
		fst := gram.GetFst(key)
		//t.Logf("%v: %v\n", key, fst)
		if key == "exp" {
			if strings.Join(fst, " ") != "DEC FALSE FLOATING ID INC INTEGER LPAREN NEW NOT NULL STRING SUB TRUE" {
				t.Fatalf("Fst of 'exp' is wrong")
			}
		}
		if key == "ini_exp" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR COMMA DEC FALSE FLOAT FLOATING ID INC INT INTEGER LPAREN NEW NOT NULL STATIC STRING SUB TRUE VOID" {
				t.Fatalf("Fst of 'ini_exp' is wrong")
			}
		}
		if key == "postfix_exp" {
			if strings.Join(fst, " ") != "FALSE FLOATING ID INTEGER LPAREN NULL STRING TRUE" {
				t.Fatalf("Fst of 'postfix_exp' is wrong")
			}
		}
		if key == "method_definition" {
			if strings.Join(fst, " ") != "BOOLEAN CHAR FLOAT ID INT STATIC VOID" {
				t.Fatalf("Fst of 'method_definition' is wrong")
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
