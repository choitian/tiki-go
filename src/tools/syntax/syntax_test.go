package syntax

import (
	"log"
	"os"
	"strings"
	"testing"
	"tools/syntax/grammar"
	"tools/syntax/parsing"
)

func assertEqual(t *testing.T, name string, expect int, value int) {
	if expect != value {
		t.Fatalf("%v size is wrong(expect %v,but is %v)", name, expect, value)
	}
}
func Test_Parsing_BuildCanonicalCollection(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/dnf0.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		assertEqual(t, "state", 249, len(lalr.States))

		kernelSum := 0
		gotoSum := 0
		for _, state := range lalr.States {
			kernelSum += len(state.GetKernelItems())
			gotoSum += len(state.GotoTable)
		}

		assertEqual(t, "kernelSum", 383, kernelSum)
		assertEqual(t, "gotoSum", 1465, gotoSum)

		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()

		TestLookaheadSum := 0
		for _, state := range lalr.States {
			for _, lookaheadSet := range state.LookaheadTable {
				TestLookaheadSum += lookaheadSet.Size()
			}
		}

		assertEqual(t, "TestLookaheadSum", 4825, TestLookaheadSum)

		lalr.BuildParsingActionTable()
	}
	{
		gram := grammar.NewGrammar("test/dnf1.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		assertEqual(t, "state", 245, len(lalr.States))

		kernelSum := 0
		gotoSum := 0
		for _, state := range lalr.States {
			kernelSum += len(state.GetKernelItems())
			gotoSum += len(state.GotoTable)
		}

		assertEqual(t, "kernelSum", 374, kernelSum)
		assertEqual(t, "gotoSum", 1511, gotoSum)

		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()

		TestLookaheadSum := 0
		for _, state := range lalr.States {
			for _, lookaheadSet := range state.LookaheadTable {
				TestLookaheadSum += lookaheadSet.Size()
			}
		}

		assertEqual(t, "TestLookaheadSum", 4791, TestLookaheadSum)

		lalr.BuildParsingActionTable()
	}
	{
		gram := grammar.NewGrammar("test/re_grammar_0.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		assertEqual(t, "state", 19, len(lalr.States))

		kernelSum := 0
		gotoSum := 0
		for _, state := range lalr.States {
			kernelSum += len(state.GetKernelItems())
			gotoSum += len(state.GotoTable)
		}
		assertEqual(t, "kernelSum", 24, kernelSum)
		assertEqual(t, "TestLookaheadSum", 42, gotoSum)

		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()

		TestLookaheadSum := 0
		for _, state := range lalr.States {
			for _, lookaheadSet := range state.LookaheadTable {
				TestLookaheadSum += lookaheadSet.Size()
			}
		}
		assertEqual(t, "TestLookaheadSum", 154, TestLookaheadSum)
		lalr.BuildParsingActionTable()

		lalr.ToXml()
	}

}
func Test_Grammer(t *testing.T) {
	//gram :=NewGrammar("test/re_grammar.txt")
	gram := grammar.NewGrammar("test/dnf0.txt")
	/*
		for _, p := range gram.Productions {
			log.Printf("%v\n", &p)
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
	if 4 != gram.Nullable.Size() {
		t.Fatalf("Nullable.Size is wrong")
	}
	if !gram.Nullable.Contains("ini_exp") {
		t.Fatalf("ini_exp is not nullable")
	}
	if !gram.Nullable.Contains("reini_exp") {
		t.Fatalf("reini_exp is not nullable")
	}
	if !gram.Nullable.Contains("test_exp") {
		t.Fatalf("test_exp is not nullable")
	}
}
func TestMain(m *testing.M) {
	log.SetOutput(os.Stderr)

	runTests := m.Run()
	os.Exit(runTests)
}
