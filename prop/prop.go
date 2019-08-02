package prop

/*
 * Prop[erties] package for grammar analysis.
 * Defines routines for computing First and Follow
 * sets of a grammar, finding cycles, and observing
 * other characteristics of grammars
*/

import (
	"fmt"
	"grammars/form"
	"grammars/sets"
)



/*
 *******************************************************************************
 *                         First/Follow Set Functions                          *
 *******************************************************************************
*/


// Returns a mapping of non-terminals to their first-sets 
func FirstSets (g *form.Item) (map[int]*sets.Set, error) {
	var firstSets map[int]*sets.Set = make(map[int]*sets.Set);

	// For each production, set its first-set if it isn't already
	for _, p := range *g {
		fmt.Printf("FirstSets: %d: ", p.Lhs);
		if fs := firstSets[p.Lhs]; fs == nil {
			fmt.Printf("nil -> Finding First(%d)\n", p.Lhs);
			// First() will install other first-sets it is forced to discover
			set, err := First(p.Lhs, sets.Set{}, g, &firstSets);

			// Bubble up any errors
			if (err != nil) {
				return make(map[int]*sets.Set), err;
			}

			// Otherwise map the new first-set to the non-terminal
			firstSets[p.Lhs] = &set;
		} else {
			fmt.Printf("Already installed!\n");
		}
	}

	return firstSets, nil;
}


// Returns the first-set for token 'tok' in grammar 'g'
func First (tok int, visited sets.Set, g *form.Item, store *map[int]*sets.Set) (sets.Set, error) {
	ps_tok  := []*form.Production{};
	first   := sets.Set{};

	// Combined copy-insert closure
	setWith := func (i int, s sets.Set) sets.Set {
		cpy := s.Copy();
		cpy.Insert(i, form.TokenCompare);
		return cpy;
	}

	// Performs lookup for first-set in store. Else creates and saves it
	getFirst := func (t int) (sets.Set, error) {
		if ptr := (*store)[t]; ptr != nil {
			return (*ptr), nil;
		}
		set, err := First(t, setWith(tok, visited), g, store);
		if err != nil {
			return sets.Set{}, err;
		}
		(*store)[t] = &set;
		return set, nil;
	}

	// If token is a terminal, return it immediately
	if form.IsTerminal(tok) {
		return sets.Set{tok}, nil;
	}

	// Otherwise it is a non-terminal. Check for any cycles
	if visited.Contains(tok, form.TokenCompare) {
		return first, nil;
	}

	// Filter productions into those starting with tok
	for i, p := range (*g) {
		if p.Lhs == tok {
			fmt.Printf("First(%d) concerns rule %d\n", tok, i);
			ps_tok = append(ps_tok, &((*g)[i]));
		} else {
			fmt.Printf("First(%d) doesn't concern %d\n", tok, i);	
		}
	}

	// For all tok productions, collect first-symbols
	for j, p := range ps_tok {
		var i int = 0;

		// If the production is empty, add epsilon to first-set
		if p.EpsilonProduction() {
			fmt.Printf("First(%d) includes epsilon!\n", tok);
			first.Insert(form.Epsilon, form.TokenCompare);
			continue;
		}

		// Else until next T or (NT not producing eps), add First sets
		rhs := (*p).Rhs;
		for i = 0; i < (len(rhs) - 1); i++ {

			// If terminal, stop and return it
			if form.IsTerminal(rhs[i]) {
				fmt.Printf("First(%d) includes terminal %d in production %d\n", tok, rhs[i], j);
				first.Insert(rhs[i], form.TokenCompare);
				break;
			}

			// If a non-terminal, get the first-set
			fmt.Printf("First(%d) includes First(%d)\n", tok, rhs[i]);
			set, err := getFirst(rhs[i]);
			if err != nil {
				return sets.Set{}, err;
			}

			// If the set produces epsilon, remove it since more follows
			eps := set.Contains(form.Epsilon, form.TokenCompare);
			set.Remove(form.Epsilon, form.TokenCompare);

			// Combine the first-set with current one
			first = sets.Union(&first, &set, form.TokenCompare);

			// If that did not produce epsilon, then its over
			if !eps {
				break;
			}
		}

		// Add last first-set to first
		set, err := getFirst(rhs[i]);
		if err != nil {
			return sets.Set{}, err;
		}
		first = sets.Union(&first, &set, form.TokenCompare); 
	}
	
	return first, nil;
}


// Returns the follow-set for the non-terminal 'tok' in grammar 'g' with first-set 'fs' 
func Follow (tok int, visited sets.Set, g *form.Item, firsts *map[int]*sets.Set, store *map[int]*sets.Set) (sets.Set, error) {
	follow := sets.Set{};

	// Combined copy-insert closure
	setWith := func (i int, s sets.Set) sets.Set {
		cpy := s.Copy();
		cpy.Insert(i, form.TokenCompare);
		return cpy;
	}

	// Performs lookup for follow-set in store. Else creates and saves it
	getFollow := func (t int) (sets.Set, error) {
		if ptr := (*store)[t]; ptr != nil {
			return (*ptr), nil;
		}
		set, err := Follow(t, setWith(tok, visited), g, firsts, store);
		if err != nil {
			return sets.Set{}, err;
		}
		(*store)[t] = &set;
		return set, nil;
	}

	// Return error if not invoked on a non-terminal
	if form.IsNonTerminal(tok) == false {
		return follow, fmt.Errorf("Follow may only be computed for non-terminals, got: %d", tok);
	}

	// If the token has already been visited, return empty set but no error
	if visited.Contains(tok, form.TokenCompare) {
		return follow, nil;
	}

	// Check each production in grammar 'g' for occurrences of 'tok'
	for _, p := range g {
		
		// Ignore epsilon productions
		if p.EpsilonProduction() {
			continue;
		}

		// For each occurrence of 'tok' in the production, update follow-set
		i := 0; rhs := p.Rhs; lhs := p.Lhs; length := len(rhs);
		for	{
			
			// Move to next occurrence of token
			for ; i < length && rhs[i] != tok; i++ {
			}

			// If no occurrences - end now
			if i >= length {
				break;
			}

			// Since there was an occurrence - collect first-follow sets
			j := 0; done := false;
			for k := i + 1; k < length && !done; k++, j++ {

				// Extract the first-set for the next token
				include := firsts[rhs[k]];

				// If the set doesn't contain epsilon, mark to stop
				done = include.Contains(form.Epsilon, form.TokenCompare);

				// Remove epsilon from the include set (if it was ever there)
				include.Remove(form.Epsilon, form.TokenCompare);

				// Merge the include set with the follow-set
				follow = sets.Union(&follow, &include, form.TokenCompare);
			}

			// If nothing after 'tok' or spanned till end: add follow of production LHS
			if j == 0 || i + j >= length {
				include, err := getFollow(lhs);
				if err != nil {
					return sets.Set{}, err;
				}
				follow = sets.Union(&follow, &include, form.TokenCompare);
			}

			// Update iterator
			i += j;
		}
	}


	return follow, nil;
}