package prop

/*
 * Prop[erties] package for grammar analysis.
 * Defines routines for computing First and Follow
 * sets of a grammar, finding cycles, and observering
 * other characteristics of grammars
*/

import (
	"fmt"
	"errors"
	"grammars/form"
	"grammars/sets"
)


/*
 *******************************************************************************
 *                                  Constants                                  *
 *******************************************************************************
*/


// The empty set
const	EmptySet	= sets.Set{};


/*
 *******************************************************************************
 *                         First/Follow Set Functions                          *
 *******************************************************************************
*/


// Returns the first-set for a non-terminal 'nt' in the given Item
// TODO Make algorithm dynamic so it doesn't discard first-set byproducts
func First (int nt, visited sets.Set, item *form.Item) (sets.Set, error) {
	ps := []*form.Production{};
	os := []*form.Production{};
	first := sets.Set{};
	hasEpsilon := false;

	setWith := func (t int, s sets.Set) sets.Set {
		return s.Copy().Insert(t, form.TokenCompare);
	}

	// Must only be called on non-terminals
	if !form.IsNonTerminal(nt) {
		return EmptySet, errors.New("First may only be invoked on non-terminal!");
	}

	// Return if a cycle is detected
	if visited.Contains(nt, form.TokenCompare) {
		return first, nil;
	}

	// Split productions into those starting with nt, and all others
	for _, p := range *item { 
		if (p.Lhs == nt) {
			ps = append(ps, &p);
		} else {
			os = append(os, &p);
		}
	}

	// Collect all terminals from nt productions. Note epsilon productions
	for _, p := range ps {
	
		// Add epsilon to the set if the rule produces it
		if (*p).EpsilonProduction() {
			first.Insert(form.Epsilon, form.TokenCompare);
			hasEpsilon = true;
			continue;
		}

		// If Rhs of production begins with terminal, then collect it
		if t := (*p).Rhs[0]; form.isTerminal(t) {
			first.Insert(t, form.TokenCompare);
			continue;
		}

		// Otherwise Rhs begins with non-terminal, so add its first-set
		firsts_nt, err := First(t, setWith(nt, visited), item);
		if err != nil {
			return EmptySet, err;
		} else {
			first = sets.Union(first, firsts_nt, form.TokenCompare);
		}
	}

	// If epsilon is not a production, do not consider any other productions
	if !hadEpsilon {
		return first, nil;
	}

	// Otherwise check all symbols after occurrences of nt in other productions
	for _, p := range os {
		rhs := (*p).Rhs;
		// Iterate up to (N-1) as last element need not be checked
		for i := 0; i < len(rhs) - 1; i++ {

			// Skip irrelevant tokens
			if rhs[i] != nt {
				continue;
			}

			// Now symbol must be nt. If next is terminal, add it and move on
			if next := rhs[i+1]; form.IsTerminal(next) {
				first.Insert(next, form.TokenCompare);
				continue;
			}

			// Otherwise it must be a non-terminal. So remove epsilon from first
			first.Remove(Epsilon, form.TokenCompare);

			// Add first-set of non-terminal to current one
			first_other, err := First(next, setWith(nt, visited), item);
			if err != nil {
				return EmptySet, err;
			}
			first = form.Union(first, first_other, form.TokenCompare);

			// Todo: Add wrappers for functions from sets that auto-include TokenCompare 
		}
	}
	return first, nil;
}