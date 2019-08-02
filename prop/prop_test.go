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
 *                              Testing Functions                              *
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