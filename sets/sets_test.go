package sets

import (
	"fmt"
	"testing"
	"strings"
)


/*
 *******************************************************************************
 *                              Support Functions                              *
 *******************************************************************************
*/


func intsToSet (ns []int) sets.Set {
	var s sets.Set = make(sets.Set, len(set));
	for i, n := range ns {
		s[i] = n;
	}
	return s;
}


func IntCompare (a, b interface{}) bool {
	return (a.(int) == b.(int));
}

func IntStringify (a interface{}) string {
	return fmt.Sprintf("%d", a.(int));
}

/*
 *******************************************************************************
 *                              Testing Functions                              *
 *******************************************************************************
*/
 

func TestSetContains (t *testing.T) {
	elements := []int{0,1,2,3};
	set      := intsToSet(elements);

	for _, k := range elements {
		if set.Contains(k, IntCompare) == false {
			t.Errorf("The set should contain: %d but says it does not!", k);
		}
	}
	
}