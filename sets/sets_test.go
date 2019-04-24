package sets

import (
	"fmt"
	"testing"
)


/*
 *******************************************************************************
 *                              Support Functions                              *
 *******************************************************************************
*/


func intsToSet (ns []int) Set {
	var s Set = make(Set, len(ns));
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

func TestSetNotContains (t *testing.T) {
	elements := []int{0,1,2,3};
	set      := Set{};

	for _, k := range elements {
		if set.Contains(k, IntCompare) == true {
			t.Errorf("The set does not contain: %d but says it not!", k);
		}
	}
}

func TestSetInsert (t *testing.T) {
	elements := []int{0,1,2,3};
	set      := Set{};

	// Insert the elements once
	for _, k := range elements {
		set.Insert(k, IntCompare);
	}

	// Record the size after insertion
	pre_size := set.Len();

	// Attempt inserting them again
	for _, k := range elements {
		set.Insert(k, IntCompare);
	}

	// Record size after duplicate insertion
	post_size := set.Len();

	// Check that duplicates don't result in additional entries
	if !(pre_size == len(elements) && post_size == pre_size) {
		t.Errorf("Insertion of duplicates should not increase set size!");
	}
}

func TestSetCopy (t *testing.T) {
	elements := []int{0,1,2,3};
	set := intsToSet(elements);

	// Make a copy of the set
	copy := set.Copy();

	// Remove an element from the copy
	copy.Remove(2, IntCompare);

	// Check that the original is unchanged
	if set.Contains(2, IntCompare) == false {
		t.Errorf("Copy cross-contaminates sets!");
	}
}

func TestSetRemove (t *testing.T) {
	elements := []int{0,1,2,3};
	set := intsToSet(elements);

	// Remove odd elements
	for i := 1; i < len(elements); i += 2 {
		set.Remove(elements[i], IntCompare);
	}

	// Ensure those elements don't exist
	if set.Len() != 2 {
		t.Errorf("Incorrect set size after removing elements");
	}

	for i := 1; i < len(elements); i += 2 {
		if set.Contains(elements[i], IntCompare) {
			t.Errorf("Removed %d but it remains in the set!", elements[i]);
		}
	}

	// Ensure removing things twice doesn't alter the set
	old_size := set.Len();
	for i := 1; i < len(elements); i += 2 {
		set.Remove(elements[i], IntCompare);
	}
	new_size := set.Len();

	if old_size != new_size {
		t.Errorf("Removing elements twice should not change set size!");
	}
}

func TestSetUnion (t *testing.T) {
	as := []int{0,2,4,6,8};
	bs := []int{1,3,5,7,9};
	cs := []int{};

	set_as := intsToSet(as);
	set_bs := intsToSet(bs);
	set_cs := intsToSet(cs);

	// Check that the union of cs and as is as
	cas := Union(&set_cs, &set_as, IntCompare);
	if cas.Len() != set_as.Len() {
		t.Errorf("Failure when performing union: {} U {...}!");
	}
	for _, k := range as {
		if !cas.Contains(k, IntCompare) {
			t.Errorf("Set should contain %d but does not!", k);
		}
	}

	// Check that union of as and bs contains all elements of both
	abs := Union(&set_as, &set_bs, IntCompare);
	if abs.Len() != (set_as.Len() + set_bs.Len()) {
		t.Errorf("Union of two distinct sets requires size be sum of both sizes!"); 
	}
	for _, k := range as {
		if !abs.Contains(k, IntCompare) {
			t.Errorf("Set should contain %d but does not!", k);
		}
	}
	for _, k := range bs {
		if !abs.Contains(k, IntCompare) {
			t.Errorf("Set should contain %d but does not!", k);
		}
	}

	// Check that union of bs and bs is just bs
	bbs := Union(&set_bs, &set_bs, IntCompare);
	if bbs.Len() != set_bs.Len() {
		t.Errorf("Set union with itself should not increase size!");
	}
	for _, k := range bs {
		if !bbs.Contains(k, IntCompare) {
			t.Errorf("Set should contain %d but does not!", k);
		}
	}

	// Check that the union of two empty sets is also an empty set
	ccs := Union(&set_cs, &set_cs, IntCompare);
	if ccs.Len() != 0 {
		t.Errorf("Union of empty sets should also be empty!");
	}
}