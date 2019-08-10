package prop

import (
	"fmt"
	"testing"
	"grammars/form"
	"grammars/sets"
)


/*
 *******************************************************************************
 *                              Helper Functions                               *
 *******************************************************************************
*/


// Converts a token (an integer) into a string format
func show (t interface{}) string {
	return fmt.Sprintf("%d", t.(int));
}


/*
 *******************************************************************************
 *                         First-Set Testing Functions                         *
 *******************************************************************************
*/


// Checks that independent productions map directly to their own terminals
func TestFirstSetsDisjointTerminals (t *testing.T) {

	// Creating rules
	p1 := form.Production{-1, []int{1}, 0};
	p2 := form.Production{-2, []int{2}, 0};
	p3 := form.Production{-3, []int{3}, 0};

	// Create the grammar
	g := form.Item{p1,p2,p3};

	// Compute the first-sets for all productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf(fmt.Sprintf("Error computing first-sets: %s", err));
	}

	// Ensure the first-sets all contain terminals assigned previously
	var set_1 *sets.Set = first_sets[-1];
	var set_2 *sets.Set = first_sets[-2];
	var set_3 *sets.Set = first_sets[-3];
	if set_1 == nil || (*set_1).Len() != 1 || (*set_1)[0] != 1 {
		t.Errorf("First set of P1 should be {1} but is %s!", set_1.String(show));
	}
	if set_2 == nil || (*set_2).Len() != 1 || (*set_2)[0] != 2 {
		t.Errorf("First set of P2 should be {2} but is %s!", set_2.String(show));
	}
	if set_3 == nil || (*set_3).Len() != 1 || (*set_3)[0] != 3 {
		t.Errorf("First set of P3 should be {3} but is %s!", set_3.String(show));
	}
}


// Checks that non-terminal productions inherit the first-sets of their sub-productions
func TestFirstSetIndirection (t *testing.T) {
	// Creating rules
	p1 := form.Production{-1, []int{-2}, 0};
	p2 := form.Production{-2, []int{-3}, 0};
	p3 := form.Production{-3, []int{1}, 0};

	// Create the grammar
	g := form.Item{p1,p2,p3};

	// Compute the first-sets for all productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf(fmt.Sprintf("Error computing first-sets: %s", err));
	}

	// Each first-set should be the same and contain 1
	for _, p := range g {
		fsp := first_sets[p.Lhs];

		if fsp == nil || fsp.Len() != 1 || (*fsp)[0] != 1 {
			t.Errorf("First(%d) should be {1} but is %s", p.Lhs, fsp.String(show));
		}
	}
}

// Verifies that a production with a nullable nonterminal includes succeeding terminals in the First set but not epsilon
func TestFirstSetEpsilon (t *testing.T) {
	// Creating rules
	p1 := form.Production{-1, []int{-2,1}, 0};
	p2 := form.Production{-2, []int{2}, 0};
	p3 := form.Production{-2, []int{}, 0};

	// Create the grammar
	g := form.Item{p1,p2,p3};

	// Compute the first-sets for all productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// First sets of both rules should contain 1,2, but not epsilon (0)
	set_1 := first_sets[-1];
	set_2 := first_sets[-2];

	if ((set_1.Len() != 2) || (set_1.Contains(1, form.TokenCompare) == false) || (set_1.Contains(2, form.TokenCompare) == false)) {
		t.Errorf("set_1 should be {1,2} but is %s!", set_1.String(show));
	}

	if ((set_2.Len() != 2) || (set_2.Contains(2, form.TokenCompare) == false) || (set_2.Contains(0, form.TokenCompare) == false)) {
		t.Errorf("set_2 should be {1,2} but is %s!", set_2.String(show));
	}

	if (set_1.Contains(form.Epsilon, form.TokenCompare)) {
		t.Errorf("set_1 must not contain epsilon!");
	}
}

// Tests that a simplified arithmetic expression grammar computes the proper first-sets
func TestSimpleArithmeticGrammarFirstSets (t *testing.T) {

	// Allows sets to be displayed as sets of integers
	strfy := func (val interface{}) string {
		intval := val.(int);
		return fmt.Sprintf("%d", intval);
	}

	// Allows for a simpler comparison
	in := func (s *sets.Set, v int) bool {
		return s.Contains(v, form.TokenCompare);
	}

	// Create productions
	p0 := form.Production{-9, []int{-1, 1}, 0};		// ? -> S $ (custom rule to introduce end marker)
	p1 := form.Production{-1, []int{-2}, 0};		// S -> E
	p2 := form.Production{-2, []int{-3, -4}, 0};	// E -> T E'
	p3 := form.Production{-4, []int{2, -3, -4}, 0};	// E' -> + T E'
	p4 := form.Production{-4, []int{3, -3, -4}, 0};	// E' -> - T E'
	p5 := form.Production{-4, []int{}, 0};			// E' -> 
	p6 := form.Production{-3, []int{-5, -6}, 0};	// T -> F T'
	p7 := form.Production{-6, []int{4, -5, -6}, 0};	// T' -> * F T'
	p8 := form.Production{-6, []int{5, -5, -6}, 0}; // T' -> / F T'
	p9 := form.Production{-6, []int{}, 0};			// T' -> 
	pa := form.Production{-5, []int{6}, 0};			// F -> x
	pb := form.Production{-5, []int{7}, 0}; 		// F -> y

	// Create the grammar
	g := form.Item{p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, pa, pb};

	// Compute first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);

	if e1 != nil {
		t.Errorf("Failed to compute first-sets: %s", e1);
	}

	// First(S[-1]) should contain {6,7}
	if s := first_sets[-1]; s.Len() != 2 || !in(s,6) || !in(s,7) {
		t.Errorf("Production S[-1] should have first-set {6,7} but has: %s", s.String(strfy));
	} 

	// First(E[-2]) should contain {6,7}
	if s := first_sets[-2]; s.Len() != 2 || !in(s,6) || !in(s,7) {
		t.Errorf("Production E[-2] should have first-set {6,7} but has: %s", s.String(strfy));
	}

	// First(T[-3]) should contain {6,7}
	if s := first_sets[-3]; s.Len() != 2 || !in(s,6) || !in(s,7) {
		t.Errorf("Production T[-3] should have first-set {6,7} but has: %s", s.String(strfy));
	}

	// First(F[-5]) should contain {6,7}
	if s := first_sets[-5]; s.Len() != 2 || !in(s,6) || !in(s,7) {
		t.Errorf("Production F[-5] should have first-set {6,7} but has: %s", s.String(strfy));
	}

	// First(E'[-4]) should contain {2,3,0}
	if s := first_sets[-4]; s.Len() != 3 || !in(s,2) || !in(s,3) || !in(s, form.Epsilon) {
		t.Errorf("Production E'[-4] should have first-set {2,3,0} but has: %s", s.String(strfy));
	}

	// First(T'[-6]) should contain {4,5,0}
	if s := first_sets[-6]; s.Len() != 3 || !in(s,4) || !in(s,5) || !in(s, form.Epsilon) {
		t.Errorf("Production T'[-6] should have first-set {4,5,0} but has: %s", s.String(strfy));
	}
}


/*
 *******************************************************************************
 *                          Follow-Set Test Functions                          *
 *******************************************************************************
*/


// Tests that a terminal after a non-terminal is included in the follow-set
func TestFollowTerminalAfterNonTerminal (t *testing.T) {
	
	// Create our productions
	p1 := form.Production{-1, []int{-2,1}, 0}; 	// S -> Ex
	p2 := form.Production{-2, []int{2}, 0};		// E -> y

	// Create the grammar
	g := form.Item{p1, p2};

	// Compute the first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);
	
	if e1 != nil {
		t.Errorf("Error computing first-sets: %s", e1);
	}

	// Compute follow-set
	follow_sets, e2 := FollowSets(&g, &first_sets);

	if e2 != nil {
		t.Errorf("Error computing follow-sets: %s", e2);
	}

	// Ensure that the follow-set for E contains only x
	if set := follow_sets[-2]; set.Len() != 1 || !set.Contains(1, form.TokenCompare) {
		t.Errorf("set should have length 1 and contain {1} but has length %d!", set.Len());
	}
}

// Tests that a non-terminal at the end of a rule includes the follow-set of the LHS production
func TestFollowIncludeLHS (t *testing.T) {
	
	// Create productions
	p1 := form.Production{-1, []int{-2,1}, 0}; 		// S -> E x
	p2 := form.Production{-2, []int{-3}, 0};		// E -> Y
	p3 := form.Production{-3, []int{2}, 0};			// Y -> z

	// Create the grammar
	g := form.Item{p1, p2, p3};
	
	// Compute the first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);
	
	if e1 != nil {
		t.Errorf("Error computing first-sets: %s", e1);
	}

	// Compute follow-set
	follow_sets, e2 := FollowSets(&g, &first_sets);

	if e2 != nil {
		t.Errorf("Error computing follow-sets: %s", e2);
	}

	// Ensure that production Y (-3) only includes (1) that it got from rule S (-1)
	if set_y := follow_sets[-3]; set_y.Len() != 1 || !set_y.Contains(1, form.TokenCompare) {
		t.Errorf("set_y should have length 1 and contain {1} but has length %d!", set_y.Len());
	}

	// Ensure that production E (-2) also only includes (1)
	if set_e := follow_sets[-2]; set_e.Len() != 1 || !set_e.Contains(1, form.TokenCompare) {
		t.Errorf("set_e should have length 1 and contain {1} but has length %d!", set_e.Len());
	}
}

// Tests that a non-terminal after another non-terminal with a non-null first-set is in the follow-set
func TestFollowIncludeNonTerminalFirstSet (t *testing.T) {
	
	// Create productions
	p1 := form.Production{-1, []int{-2, -3}, 0}; // S -> E F
	p2 := form.Production{-3, []int{1}, 0};		 // F -> x
	p3 := form.Production{-2, []int{2}, 0};		 // E -> y 

	// Create the grammar
	g := form.Item{p1, p2, p3};
	
	// Compute the first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);
	
	if e1 != nil {
		t.Errorf("Error computing first-sets: %s", e1);
	}

	// Compute follow-set
	follow_sets, e2 := FollowSets(&g, &first_sets);

	if e2 != nil {
		t.Errorf("Error computing follow-sets: %s", e2);
	}


	// Ensure that the production E (-2) contains x (1) that it got from rule F (-3)
	if set_e := follow_sets[-2]; set_e.Len() != 1 || !set_e.Contains(1, form.TokenCompare) {
		t.Errorf("Rule E (-2) should contain x (1) since that is in the first-set of NT after it!");
	} 

}

// Tests that a simplified arithmetic expression grammar computes the proper follow-sets
func TestSimpleArithmeticGrammarFollowSets (t *testing.T) {

	// Allows sets to be displayed as sets of integers
	strfy := func (val interface{}) string {
		intval := val.(int);
		return fmt.Sprintf("%d", intval);
	}

	// Allows for a simpler comparison
	in := func (s *sets.Set, v int) bool {
		return s.Contains(v, form.TokenCompare);
	}

	// Create productions
	p0 := form.Production{-9, []int{-1, 1}, 0};		// ? -> S $ (custom rule to introduce end marker)
	p1 := form.Production{-1, []int{-2}, 0};		// S -> E
	p2 := form.Production{-2, []int{-3, -4}, 0};	// E -> T E'
	p3 := form.Production{-4, []int{2, -3, -4}, 0};	// E' -> + T E'
	p4 := form.Production{-4, []int{3, -3, -4}, 0};	// E' -> - T E'
	p5 := form.Production{-4, []int{}, 0};			// E' -> 
	p6 := form.Production{-3, []int{-5, -6}, 0};	// T -> F T'
	p7 := form.Production{-6, []int{4, -5, -6}, 0};	// T' -> * F T'
	p8 := form.Production{-6, []int{5, -5, -6}, 0}; // T' -> / F T'
	p9 := form.Production{-6, []int{}, 0};			// T' -> 
	pa := form.Production{-5, []int{6}, 0};			// F -> x
	pb := form.Production{-5, []int{7}, 0}; 		// F -> y

	// Create the grammar
	g := form.Item{p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, pa, pb};

	// Compute first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);

	if e1 != nil {
		t.Errorf("Error computing first-sets: %s", e1);
	}

	// Compute follow-set
	follow_sets, e2 := FollowSets(&g, &first_sets);

	if e2 != nil {
		t.Errorf("Error computing follow-sets: %s", e2);
	}


	// Follow(S[-1]) and Follow(E[-2]) and Follow(E'[-4]) should have follow-set {1}
	if s := follow_sets[-1]; s.Len() != 1 || !in(s, 1) {
		t.Errorf("Production S[-1] should have follow-set {1} but has: %s", s.String(strfy));
	}
	if s := follow_sets[-2]; s.Len() != 1 || !in(s, 1) {
		t.Errorf("Production E[-2] should have follow-set {1} but has: %s", s.String(strfy));
	}
	if s := follow_sets[-4]; s.Len() != 1 || !in(s, 1) {
		t.Errorf("Production E'[-4] should have follow-set {1} but has: %s", s.String(strfy));
	}


	// Follow(T[-3]) and Follow(T'[-6]) should have follow-set {2,3,1}
	if s := follow_sets[-3]; s.Len() != 3 || !in(s,1) || !in(s,2) || !in(s,3) {
		t.Errorf("Follow(T[-3]) should have follow-set {2,3,1} but has: %s", s.String(strfy));
	}
	if s := follow_sets[-6]; s.Len() != 3 || !in(s,1) || !in(s,2) || !in(s,3) {
		t.Errorf("Follow(T'[-6]) should have follow-set {2,3,1} but has: %s", s.String(strfy));
	}

	
	// Follow(F[-5]) should have follow-set {4,5,2,3,1}
	if s := follow_sets[-5]; s.Len() != 5 || !in(s,1) || !in(s,2) || !in(s,3) || !in(s,4) || !in(s,5) {
		t.Errorf("Follow(F[-5]) should have follow-set {4,5,2,3,1} but has: %s", s.String(strfy));
	}

}


// Tests combined follow-set properties
func TestFollowSetPropertiesCombined (t *testing.T) {

	// Allows sets to be displayed as sets of integers
	strfy := func (val interface{}) string {
		intval := val.(int);
		return fmt.Sprintf("%d", intval);
	}

	// Allows for a simpler comparison
	in := func (s *sets.Set, v int) bool {
		return s.Contains(v, form.TokenCompare);
	}

	// Create productions
	p1 := form.Production{-1, []int{-2, 1}, 0};								// P -> S w
	p2 := form.Production{-2, []int{-3, 2, -3, 3, -3, 4, -3, -4, -5}, 0};	// S -> E x E y E z E Q R
	p3 := form.Production{-4, []int{5}, 0};									// Q -> a
	p4 := form.Production{-4, []int{}, 0};									// Q -> 
	p5 := form.Production{-5, []int{6}, 0};									// R -> b
	p6 := form.Production{-5, []int{}, 0};									// R -> 
	p7 := form.Production{-3, []int{7}, 0};									// E -> e

	// Create the grammar
	g := form.Item{p1, p2, p3, p4, p5, p6, p7};

	// Compute first-sets for all productions (needed for the follow-set)
	first_sets, e1 := FirstSets(&g);

	if e1 != nil {
		t.Errorf("Error computing first-sets: %s", e1);
	}

	// Compute follow-set
	follow_sets, e2 := FollowSets(&g, &first_sets);

	if e2 != nil {
		t.Errorf("Error computing follow-sets: %s", e2);
	}

	// Follow(-3) should be set: {2,3,4,5,6,1}
	if s := follow_sets[-3]; s.Len() != 6 || !in(s,1) || !in(s,2) || !in(s,3) || !in(s,4) || !in(s,5) || !in(s,6) {
		t.Errorf("Follow(-3) should be {2,3,4,5,6,1} but is: %s\n", s.String(strfy));
	}

}


/*
 *******************************************************************************
 *                        Left-Recursion Test Functions                        *
 *******************************************************************************
*/

// Tests trivial Left-Recursion
func TestLeftRecursionImmediateCase (t *testing.T) {
	
	// Create Production
	p := form.Production{-1, []int{-1}, 0};	// A -> A

	// Create grammar
	g := form.Item{p};

	// Compute first-sets for productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// Check that a cycle exists for A
	if isCycleA, _ := IsLeftRecursive(-1, &g, &first_sets); !isCycleA {
		t.Errorf("Production A[-1] has a cycle that was not found!");
	}
}

// Tests trivial case of no Left-Recursion
func TestLeftRecursionImmediateCaseNone (t *testing.T) {
	
	// Create Production
	p := form.Production{-1, []int{1}, 0};	// A -> a

	// Create grammar
	g := form.Item{p};

	// Compute first-sets for productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// Check that a cycle exists for A
	if isCycleA, _ := IsLeftRecursive(-1, &g, &first_sets); isCycleA {
		t.Errorf("Production A[-1] does not have a cycle - yet one was found!");
	}
}

// Tests case where left-recursion occurs through two non-terminal rewrites
func TestLeftRecursionIndirectRewrites (t *testing.T) {
	
	// Create Production
	p1 := form.Production{-1, []int{-2}, 0};	// A -> B
	p2 := form.Production{-2, []int{-1}, 0};	// B -> A

	// Create grammar
	g := form.Item{p1, p2};

	// Compute first-sets for productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// Check that a cycle exists for A
	if isCycleA, _ := IsLeftRecursive(-1, &g, &first_sets); !isCycleA {
		t.Errorf("Production A[-1] has a cycle!");
	}
}

// Tests case where left-recursion is obfuscated by another nonterminal with an epsilon production
func TestLeftRecursionEpsilonObfuscation (t *testing.T) {
	
	// Create Production
	p1 := form.Production{-1, []int{-2, -1}, 0};	// A -> BA
	p2 := form.Production{-2, []int{1}, 0};			// B -> a
	p3 := form.Production{-2, []int{}, 0};			// B -> 

	// Create grammar
	g := form.Item{p1, p2, p3};

	// Compute first-sets for productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// Check that a cycle exists for A
	if isCycleA, _ := IsLeftRecursive(-1, &g, &first_sets); !isCycleA {
		t.Errorf("Production A[-1] has a cycle!");
	}
}

// Tests Left-Recursion on a grammar with nullable nonterminals and blocking terminals
func TestLeftRecursionCombined (t *testing.T) {

	// Create productions
	p1 := form.Production{-1, []int{1, -2, 2}, 0};	// A -> aBb
	p2 := form.Production{-1, []int{-3, -2, 3}, 0};	// A -> CBc
	p3 := form.Production{-2, []int{4}, 0};			// B -> d
	p4 := form.Production{-2, []int{-3, -1}, 0};	// B -> CA
	p5 := form.Production{-3, []int{5}, 0};			// C -> e
	p6 := form.Production{-3, []int{}, 0};			// C ->

	// Create the grammar
	g := form.Item{p1, p2, p3, p4, p5, p6};

	// Compute first-sets for all productions
	first_sets, err := FirstSets(&g);

	if err != nil {
		t.Errorf("Error computing first-sets: %s", err);
	}

	// Check that A has a cycle
	if isCycleA, _ := IsLeftRecursive(-1, &g, &first_sets); !isCycleA {
		t.Errorf("Production A[-1] has a cycle that was not found!");
	}

	// Check that B has a cycle
	if isCycleB, _ := IsLeftRecursive(-2, &g, &first_sets); !isCycleB {
		t.Errorf("Production B[-2] has a cycle that was not found!");
	}

	// Check that C does not have a cycle
	if isCycleC, _ := IsLeftRecursive(-3, &g, &first_sets); isCycleC {
		t.Errorf("Production C[-3] does not have a cycle - yet one was found!");
	}
	
}