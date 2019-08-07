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


// Returns a mapping of tokens to their first-sets 
func FirstSets (g *form.Item) (map[int]*sets.Set, error) {
	var firstSets map[int]*sets.Set = make(map[int]*sets.Set);

	// For each production NT, set its first-set if it isn't already
	for _, p := range *g {

		// First install all terminals into the map (if any)
		for _, t := range p.Rhs {
			if form.IsTerminal(t) {
				set := sets.Set{t}; 
				firstSets[t] = &set;
			}
		}

		if fs := firstSets[p.Lhs]; fs == nil {
			// First() will install other first-sets it is forced to discover
			set, err := First(p.Lhs, sets.Set{}, g, &firstSets);

			// Bubble up any errors
			if (err != nil) {
				return make(map[int]*sets.Set), err;
			}

			// Otherwise map the new first-set to the non-terminal
			firstSets[p.Lhs] = &set;
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

	// Otherwise it is a non-terminal. Check for any cycles
	if visited.Contains(tok, form.TokenCompare) {
		return first, nil;
	}

	// Filter productions into those starting with tok
	for i, p := range (*g) {
		if p.Lhs == tok {
			ps_tok = append(ps_tok, &((*g)[i]));
		}
	}

	// For all tok productions, collect first-symbols
	for _, p := range ps_tok {
		var i int = 0;

		// If the production is empty, add epsilon to first-set
		if p.EpsilonProduction() {
			first.Insert(form.Epsilon, form.TokenCompare);
			continue;
		}

		// Else until next T or (NT not producing eps), add First sets
		rhs := (*p).Rhs;
		for i = 0; i < (len(rhs) - 1); i++ {

			// If terminal, stop and return it
			if form.IsTerminal(rhs[i]) {
				first.Insert(rhs[i], form.TokenCompare);
				break;
			}

			// If a non-terminal, get the first-set
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

		// If last is terminal add it, else merge first-set of non-terminal
		if form.IsTerminal(rhs[i]) {
			first.Insert(rhs[i], form.TokenCompare);
		} else {
			set, err := getFirst(rhs[i]);
			if err != nil {
				return sets.Set{}, err;
			}
			first = sets.Union(&first, &set, form.TokenCompare);
		}
	}
	
	return first, nil;
}


// Returns a mapping of non-terminals to their follow-sets 
func FollowSets (g *form.Item, firstSets *map[int]*sets.Set) (map[int]*sets.Set, error) {
	var followSets map[int]*sets.Set = make(map[int]*sets.Set);

	// For each production, set its follow-set if it isn't already
	for _, p := range *g {
		if fs := followSets[p.Lhs]; fs == nil {

			// Follow() will install other first-sets it is forced to discover
			set, err := Follow(p.Lhs, sets.Set{}, g, firstSets, &followSets);

			// Bubble up any errors
			if (err != nil) {
				return make(map[int]*sets.Set), err;
			}

			// Otherwise map the new first-set to the non-terminal
			followSets[p.Lhs] = &set;
		}
	}

	return followSets, nil;
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
	for _, p := range *g {
		
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
			for k := i + 1; k < length && !done; k, j = k + 1, j + 1 {

				// Extract the first-set for the next token
				include := ((*firsts)[rhs[k]]).Copy();

				// If the set doesn't contain epsilon, mark to stop
				done = !include.Contains(form.Epsilon, form.TokenCompare);

				// Remove epsilon from the include set (if it was ever there)
				include.Remove(form.Epsilon, form.TokenCompare);

				// Merge the include set with the follow-set
				follow = sets.Union(&follow, &include, form.TokenCompare);
			}

			// If nothing after 'tok' or spanned till end: add follow of production LHS
			if !done || j == 0 {
				include, err := getFollow(lhs)
				if err != nil {
					return sets.Set{}, err;
				}
				follow = sets.Union(&follow, &include, form.TokenCompare);
				break;
			}

			// Update iterator
			i += j;
		}
	}

	return follow, nil;
}


/*
 *******************************************************************************
 *      Left-to-right Leftmost-derivation (LL) Grammar Property Functions      *
 *******************************************************************************
*/


// Returns true if the given non-terminal 'tok' is left-recursive
func IsLeftRecursive (tok int, g *form.Item, firsts *map[int]*sets.Set, visited sets.Set) bool {
	return isLeftRecursive(tok, tok, g, firsts, sets.Set{});
}


// Searches for a cycle in a rule
func isLeftRecursive (rule, target int, g *form.Item, firsts *map[int]*sets.Set, visited sets.Set) bool {

	// Assume not found
	result := false;

	// Combined copy-insert closure
	setWith := func (i int, s sets.Set) sets.Set {
		cpy := s.Copy();
		cpy.Insert(i, form.TokenCompare);
		return cpy;
	}


	// For all productions of the currently searched rule, look for occurrence of tok
	for _, p := range *g {

		// Ignore irrelevant rules and epsilon productions of relevant rules
		if p.Lhs != rule || p.EpsilonProduction() {
			continue;
		}

		// For each component of the production ... 
		for i := 0; i < len(p.Rhs); i++ {

			// Exit upon discovering a terminal as no cycle is possible now
			if form.IsTerminal(p.Rhs[i]) {
				break;
			}

			// If non-terminal is sought one - mark true and break
			if p.Rhs[i] == target {
				result = true;
				break;
			}

			// Otherwise different non-terminal. If not already searched - search it
			if !visited.Contains(p.Rhs[i], form.TokenCompare) {
				result = result || isLeftRecursive(p.Rhs[i], target, g, firsts, setWith(rule, visited));
			}

			// If the first-set doesn't contain epsilon, then stop
			if f, ok := (*firsts)[p.Rhs[i]]; ok && !f.Contains(form.Epsilon, form.TokenCompare) {
				break;
			}	
		}
	}

	// Return outcome
	return result;

}