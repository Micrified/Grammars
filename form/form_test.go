package form

import (
	"testing"
	"grammars/symtab"
	"grammars/parse"
	"strings"
)


func TestProductionAll (t *testing.T) {
	var a Production = Production{-1, []int{1,-2,2}, 0};
	var b Production = Production{-2, []int{3,4}, 0};
	var c Production = Production{-3, []int{}, 0};

	// Testing IsNonTerminal (positive and negative test)
	if !IsNonTerminal(a.Lhs) {
		t.Errorf("Value %d incorrectly not considered a non-terminal!", a.Lhs);
	}
	if IsNonTerminal(a.Rhs[0]) {
		t.Errorf("Value %d incorrectly considered a non-terminal!", a.Rhs[0]);
	}

	// Testing IsTerminal (positive and negative test)
	if !IsTerminal(b.Rhs[0]) {
		t.Errorf("Value %d incorrectly not considered a terminal!", b.Rhs[0]);
	}
	if IsTerminal(b.Lhs) {
		t.Errorf("Value %d incorrectly considered a terminal!", b.Lhs);
	}

	// Testing Epsilon production (positive and negative test)
	if c.EpsilonProduction() == false {
		t.Errorf("Set with zero productions incorrectly not considered epsilon!");
	}
	if a.EpsilonProduction() == true {
		t.Errorf("Set with nonzero productions incorrectly considered epsilon!");
	}

	// Create a simple symbol table
	s := symtab.SymTab{[]string{}, []string{}};

	symtab.RegisterNonTerminal("A", &s);
	symtab.RegisterNonTerminal("B", &s);
	symtab.RegisterTerminal("a", &s);
	symtab.RegisterTerminal("b", &s);

	// Expected string form of production a: 
	a_str := "A " + parse.DefineOperator + " a B b";

	if strings.Compare(a.String(false, &s), a_str) != 0 {
		t.Errorf("The outputted string form:\n\"%s\"\nDoes not match:\n\"%s\"\n", 
			a.String(false, &s), a_str);
	}
	
}

func TestItemAll (t *testing.T) {
	var a Production = Production{-1, []int{1,-2,2}, 0};
	var b Production = Production{-2, []int{3,4}, 0};
	var d Item = []Production{a,b};
	var e Item = []Production{};

	// Test IsEmpty (positive and negative)
	if d.IsEmpty() == true {
		t.Errorf("Nonempty production was reported empty!");
	}
	if e.IsEmpty() == false {
		t.Errorf("Empty production reported nonempty!");
	}
	
	// Not testing the string form of Items since its just productions
 	// In a list. 
}