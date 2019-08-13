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
	"grammars/symtab"
	"grammars/parse"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Describes rewrite where a non-terminal is replaced by its definition
type Rewrite struct {
	Rule	int;				// The index of the production being considered
	Index	int;				// The index within the source production
}


// Describes clash between two productions, each represented by grammar index
type Clash struct {
	P1		int;
	P2		int;
}


/*
 *******************************************************************************
 *                         First/Follow Set Functions                          *
 *******************************************************************************
*/


// Returns a mapping of tokens to their first-sets 
func FirstSets (g *form.Item) (map[int]*sets.Set, error) {
	var firstSets map[int]*sets.Set = make(map[int]*sets.Set);

	// For each production NT, install epsilon rules before any others
	for _, p := range *g {
		if p.EpsilonProduction() {
			set := sets.Set{form.Epsilon};
			firstSets[p.Lhs] = &set;
		}
	}

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


// Returns true if non-terminal 'tok' is left-recursive. Rewrites describe why
func IsLeftRecursive (tok int, g *form.Item, firsts *map[int]*sets.Set) (bool, []Rewrite) {
	rws := isLeftRecursive(tok, tok, g, firsts, sets.Set{});
	fmt.Println("isLeftRecursive() ...");
	if len(rws) != 0 {
		for i := 0; i < len(rws); i++ {
			fmt.Println(rws[i].String(g, nil));
		}
	}
	return len(rws) != 0, rws;
}

// Determines if a given 'nt' has a cycle. Returns nonempty rewrite sequence if true
func isLeftRecursive (rule, nt int, g *form.Item, firsts *map[int]*sets.Set, visited sets.Set) []Rewrite {

	// Assume no recursion (empty rewrites)
	rws := []Rewrite{};

	// Combined copy-insert closure
	setWith := func (i int, s sets.Set) sets.Set {
		cpy := s.Copy();
		cpy.Insert(i, form.TokenCompare);
		return cpy;
	}

	// For all productions of the currently searched rule, look for occurrence of nt
	for r, p := range *g {

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
			if p.Rhs[i] == nt {
				rws = append(rws, Rewrite{r, i});
				break;
			}

			// Otherwise different non-terminal. If not already searched - search it
			if !visited.Contains(p.Rhs[i], form.TokenCompare) {
				got := isLeftRecursive(p.Rhs[i], nt, g, firsts, setWith(rule, visited));
				if len(got) != 0 {
					rws = append(got, Rewrite{r, i});
				}
				// result = result || isLeftRecursive(p.Rhs[i], nt, g, firsts, setWith(rule, visited));
			}

			// If the first-set doesn't contain epsilon, then stop
			if f, ok := (*firsts)[p.Rhs[i]]; ok && !f.Contains(form.Epsilon, form.TokenCompare) {
				break;
			}	
		}
	}

	// Return outcome
	return rws;

}


/*
 *******************************************************************************
 *                    First-Set-Clash Function Definitions                     *
 *******************************************************************************
*/


// Returns true if a grammar contains first-set clashes
func IsFirstSetClash (g *form.Item, firsts *map[int]*sets.Set) (bool, Clash) {

	// Collect unique productions
	uniqueProductions := sets.Set{};
	for _, p := range (*g) {
		uniqueProductions.Insert(p.Lhs, form.TokenCompare);
	}

	// Check unique productions for clashes
	for _, nt := range uniqueProductions {
		if isClash, clash := isFirstSetClash(nt.(int), g, firsts); isClash {
			return isClash, clash;
		}
	}

	return false, Clash{};
}


// Returns true if a production has a first-set clash. 
func isFirstSetClash (nt int, g *form.Item, firsts *map[int]*sets.Set) (bool, Clash) {

	// Stores indexes of productions in a grammar
	var p_indices []int = []int{};

	// Performs a lookup and returns a set for a given first-set map. Auto-panics
	firstOf := func (nt int, firsts *map[int]*sets.Set) *sets.Set {
		set, ok := (*firsts)[nt];
		if !ok {
			panic(fmt.Errorf("Invalid first-set lookup. Nothing for %d in map!", nt));
		}
		return set;
	}

	// Possible to only need one set, but tracking sets of sets for better feedback
	var p_firsts []sets.Set = []sets.Set{}; 

	// Collect first-sets of all productions in the grammar beginning with 'nt'
	for idx, p := range (*g) {
		var first sets.Set = sets.Set{};
		var i int;

		// Ignore irrelevant productions
		if p.Lhs != nt {
			continue
		}

		// Store the index of the production
		p_indices = append(p_indices, idx);

		// If production first-set is empty, store epsilon and move on
		if p.EpsilonProduction() {
			first = append(first, form.Epsilon);

		} else {

			// Until T or NT not producing epsilon, add first-set
			for i = 0; i < (len(p.Rhs) - 1); i++ {
			
				// If terminal, stop and store set
				if form.IsTerminal(p.Rhs[i]) {
					first = append(first, p.Rhs[i]);
					break;
				}

				// If non-terminal, get first-set of that rule
				set := firstOf(p.Rhs[i], firsts);

				// If set contains epsilon, remove since more follows
				hasEpsilon := set.Contains(form.Epsilon, form.TokenCompare);
				set.Remove(form.Epsilon, form.TokenCompare);

				// Combine first-set with current one
				first = sets.Union(&first, set, form.TokenCompare);

				// If epsilon was not produced, break
				if !hasEpsilon {
					break;
				}
			}

			// Case: Last element: If terminal, add it. Else merge first-set of NT
			if form.IsTerminal(p.Rhs[i]) {
				first.Insert(p.Rhs[i], form.TokenCompare);
			} else {
				set := firstOf(p.Rhs[i], firsts);
				first = sets.Union(&first, set, form.TokenCompare);
			}
		}

		// Check if any other set so far intersects the current set
		for j := 0; j < len(p_firsts); j++ {
			if isect := sets.Intersect(&first, &(p_firsts[j]), form.TokenCompare); isect.Len() != 0 {
				return true, Clash{idx, p_indices[j]};
			}
		}

		
		// Otherwise simply insert the set and move on
		p_firsts = append(p_firsts, first);
	}

	// Nothing was found
	return false, Clash{};	
}


// Returns true if grammar has first-set clash. First-sets must be valid
/*func IsFirstSetClash (g *form.Item, firsts *map[int]*sets.Set) (bool, Clash) {
	for i := 0; i < len(*g) - 1; i++ {
		for j := i + 1; j < len(*g); j++ {

			// Extract the productions to compare
			p1 := ((*g)[i]).Lhs;
			p2 := ((*g)[j]).Lhs;

			// Don't compare productions from the same non-terminal
			if p1 == p2 {
				continue;
			}

			// Recover the first-sets
			s1, ok1 := (*firsts)[p1];
			s2, ok2 := (*firsts)[p2];

			// Panic if any error occurred
			if !ok1 || !ok2 {
				panic("Invalid first-sets provided to first set clash!");
			}

			// If the sets intersect, report the clash
			if s := sets.Intersect(s1, s2, form.TokenCompare); s.Len() > 0 {
				return true, Clash{i, j};
			}
		}
	}
	return false, Clash{0,0};
}*/



/*
 *******************************************************************************
 *                        Rewrite Function Definitions                         *
 *******************************************************************************
*/


// Returns a string describing a rewrite for a given grammar and symbol table
func (r *Rewrite) String (g *form.Item, tab *symtab.SymTab) string {
	var s string;

	// Returns a string describing the given token, depending on state of tab
	tokStr := func (tok int) string {
		if tab == nil {
			return fmt.Sprintf("%d", tok);
		}
		str, err := symtab.LookupID(tok, tab);
		if err != nil {
			return fmt.Sprintf("<err: unregistered symbol %d>", tok);
		}
		return str;
	}

	// Return placeholder if the rule cannot be logically indexed in 'g'
	if n := len(*g); r.Rule >= n || r.Rule < 0 {
		return "<err: bad rule>";
	}

	// Ready the production and rewrite index
	p := (*g)[r.Rule]; i := r.Index;

	// Write the LHS and definition operator
	s = fmt.Sprintf("%s %s ", tokStr(p.Lhs), parse.DefineOperator);

	// Write epsilons until the index 
	for j := 0; j < i; j++ {
		s = s + "ε ";
	}

	// Write the rewritten token
	s = s + fmt.Sprintf("{{{ %s }}} ", tokStr(p.Rhs[i]));

	// Write remaining tokens
	for j := i + 1; j < len(p.Rhs); j++ {
		s = s + tokStr(p.Rhs[j]) + " ";
	}
	
	// Return string but trim last space
	return s[:len(s) - 1];
}
