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
func Test_BuildCanonicalCollection(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")
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
	}
	{
		gram := grammar.NewGrammar("test/dnf.txt")
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
	}
}
func Test_ClosureWithLookahead(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		kernel := lalr.InitialItem
		dummyState := parsing.NewState(kernel)
		dummyState.AddLookahead(kernel, grammar.SymbolEnd)
		lalr.ClosureWithLookahead(dummyState)

		DummyStateLookaheadSum := 0
		for _, lookaheadSet := range dummyState.LookaheadTable {
			DummyStateLookaheadSum += lookaheadSet.Size()
		}
		assertEqual(t, "DummyStateLookaheadSum", 86, DummyStateLookaheadSum)
	}

	{
		gram := grammar.NewGrammar("test/dnf.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		kernel := lalr.InitialItem
		dummyState := parsing.NewState(kernel)
		dummyState.AddLookahead(kernel, grammar.SymbolEnd)
		lalr.ClosureWithLookahead(dummyState)

		DummyStateLookaheadSum := 0
		for _, lookaheadSet := range dummyState.LookaheadTable {
			DummyStateLookaheadSum += lookaheadSet.Size()
		}
		assertEqual(t, "DummyStateLookaheadSum", 32, DummyStateLookaheadSum)
	}
}
func Test_BuildPropagateAndSpontaneousTable(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()
		lalr.BuildPropagateAndSpontaneousTable()

		spontaneousSum, propagateSum := 0, 0
		for _, kernel := range lalr.ItemPool {
			if kernel.IsKernel() {
				spontaneousSum += len(kernel.SpontaneousTable)
				propagateSum += len(kernel.PropagateTable.Values())
			}
		}
		assertEqual(t, "spontaneousSum", 45, spontaneousSum)
		assertEqual(t, "propagateSum", 47, propagateSum)
	}
	{
		gram := grammar.NewGrammar("test/dnf.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()
		lalr.BuildPropagateAndSpontaneousTable()

		spontaneousSum, propagateSum := 0, 0
		for _, kernel := range lalr.ItemPool {
			if kernel.IsKernel() {
				spontaneousSum += len(kernel.SpontaneousTable)
				propagateSum += len(kernel.PropagateTable.Values())
			}
		}
		assertEqual(t, "spontaneousSum", 2289, spontaneousSum)
		assertEqual(t, "propagateSum", 1240, propagateSum)
	}
}
func Test_DoPropagation(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()
		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()

		TestLookaheadSum := 0
		for _, state := range lalr.States {
			for _, lookaheadSet := range state.LookaheadTable {
				TestLookaheadSum += lookaheadSet.Size()
			}
		}
		assertEqual(t, "TestLookaheadSum", 154, TestLookaheadSum)
	}
	{
		gram := grammar.NewGrammar("test/dnf.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()

		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()

		TestLookaheadSum := 0
		for _, state := range lalr.States {
			for _, lookaheadSet := range state.LookaheadTable {
				TestLookaheadSum += lookaheadSet.Size()
			}
		}

		assertEqual(t, "TestLookaheadSum", 4791, TestLookaheadSum)
	}
}
func Test_BuildParsingActionTable(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()
		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()
		lalr.BuildParsingActionTable()

		TestParsingActionSum := 0
		for _, state := range lalr.States {
			TestParsingActionSum += len(state.ParsingActionTable)
		}
		assertEqual(t, "TestParsingActionSum", 125, TestParsingActionSum)

		lalr.ToXml("re_grammar.xml")
	}
	{
		gram := grammar.NewGrammar("test/dnf.txt")
		lalr := parsing.NewLookaheadLR(gram)
		lalr.BuildCanonicalCollection()
		lalr.BuildPropagateAndSpontaneousTable()
		lalr.DoPropagation()
		lalr.BuildParsingActionTable()

		TestParsingActionSum := 0
		for _, state := range lalr.States {
			TestParsingActionSum += len(state.ParsingActionTable)
		}
		assertEqual(t, "TestParsingActionSum", 2780, TestParsingActionSum)

		lalr.ToXml("dnf.xml")
	}
}
func Test_Grammer(t *testing.T) {
	{
		gram := grammar.NewGrammar("test/re_grammar.txt")

		if 1 != gram.Nullable.Size() {
			t.Fatalf("Nullable.Size is wrong")
		}

		for key, _ := range gram.FST {
			fst := gram.GetFst(key)
			//t.Logf("%v: %v\n", key, fst)
			if key == "exp" {
				fstStr := strings.Join(fst, " ")
				expectFst := "ARRAY CHAR DOT LITERAL LPAREN"
				if fstStr != expectFst {
					t.Fatalf("Fst of 'exp' is wrong:get: %v /expect %v", fstStr, expectFst)
				}
			}
		}
	}
	{
		gram := grammar.NewGrammar("test/dnf.txt")

		if 4 != gram.Nullable.Size() {
			t.Fatalf("Nullable.Size is wrong")
		}
		if !gram.Nullable.Contains("ini_exp") {
			t.Fatalf("ini_exp is not nullable")
		}

		for key, _ := range gram.FST {
			fst := gram.GetFst(key)
			//t.Logf("%v: %v\n", key, fst)
			if key == "exp" {
				fstStr := strings.Join(fst, " ")
				expectFst := "DEC FALSE FLOATING ID INC INTEGER LPAREN NEW NOT NULL SCOPE_ID STRING SUB TRUE"
				if fstStr != expectFst {
					t.Fatalf("Fst of 'exp' is wrong:get: %v /expect %v", fstStr, expectFst)
				}
			}
			if key == "ini_exp" {
				fstStr := strings.Join(fst, " ")
				expectFst := "BOOLEAN CHAR COMMA DEC FALSE FLOAT FLOATING ID INC INT INTEGER LPAREN NEW NOT NULL SCOPE_ID STATIC STRING SUB TRUE VOID"
				if fstStr != expectFst {
					t.Fatalf("Fst of 'ini_exp' is wrong:get: %v /expect %v", fstStr, expectFst)
				}
			}
			if key == "method_definition" {
				fstStr := strings.Join(fst, " ")
				expectFst := "BOOLEAN CHAR FLOAT ID INT SCOPE_ID STATIC VOID"
				if fstStr != expectFst {
					t.Fatalf("Fst of 'method_definition' is wrong:get: %v /expect %v", fstStr, expectFst)
				}
			}
		}
	}
}
func TestMain(m *testing.M) {
	log.SetOutput(os.Stderr)

	runTests := m.Run()
	os.Exit(runTests)
}
