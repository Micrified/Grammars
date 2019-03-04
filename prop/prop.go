package prop

/*
 * Prop[erties] package for grammar analysis.
 * Defines routines for computing the First and Follow
 * sets of a grammar, finding cycles, and observing
 * other properties used in identifying a grammar class
 *
*/

import (
	"fmt"
	"grammars/form"
)


/*
 *******************************************************************************
 *                       Functions for First/Follow Sets                       *
 *******************************************************************************
*/


// Returns the first-set for non-terminal nt in grammar g
// [Note: Wasteful since discards first-sets it is not looking for ]
func FirstSet (nt rune, visited []rune, g form.Item) ([]rune, error) {
	var nt_ps []*form.Prod		= []*form.Prod{};
	var ot_ps []*form.Prod		= []*form.Prod{};
	var first []rune			= []rune{};
	var hasEpsilon bool			= false;

	//fmt.Printf("FirstSet invoked on %c\n", nt);

	// Fail if nt is not in fact a non-terminal
	if !form.IsNonTerminal(nt) {
		return first, fmt.Errorf("First-Set: %c is not a non-terminal!", nt);
	}

	// Return if a cycle is detected
	for _, s := range visited {
		if s == nt {
			fmt.Printf("FirstSet: Cycle -> Returning...\n");
			return first, nil;
		}
	};

	// Sort productions into those of nt, and others
	for i := range g.Ps {
		if (g.Ps[i]).Lhs == nt {
			nt_ps = append(nt_ps, &(g.Ps[i]));
		} else {
			ot_ps = append(ot_ps, &(g.Ps[i]));
		}
	}

	// Collect all terminals. Note if epsilon is a production
	for _, p := range nt_ps {
		rhs := (*p).Rhs;
		//fmt.Printf("Analysis of production: %s\n", p.String(false));
		// Note if production is empty (epsilon)
		if p.Epsilon() {
			//fmt.Printf("Rule produces epsilon!\n");
			first = form.SetInsert(first, 'ε');
			hasEpsilon = true;
			continue;
		}

		// If production rhs begins with terminal, collect it
		if form.IsTerminal(rhs[0]) {
			//fmt.Printf("Appending %c to set!\n", rhs[0]);
			first = form.SetInsert(first, rhs[0]);
			continue;
		}
	
		// If rhs begins with non-terminal, add its first-set
		//fmt.Printf("First(%c) is in First(%c) ...\n", rhs[0], nt);
		ot_first, err := FirstSet(rhs[0], append(visited, nt), g);
		if err != nil {
			return first, err;
		}
		first = form.SetUnion(first, ot_first);
	}

	// If epsilon not produced, may immediately return
	if !hasEpsilon {
		//fmt.Printf("First(%c) Returning early...\n", nt);
		return first, nil;
	}

	// Else check all symbols after which nt occurs in other productions
	for _, p := range ot_ps {
		rhs := (*p).Rhs;

		// Iterate to (N-2) as last element need not be checked
		for i := 0; i < len(rhs) - 1; i++ {
			if rhs[i] != nt {
				continue;
			}

			// If next symbol is terminal, add it and move on
			if next := rhs[i+1]; form.IsTerminal(next) {
				//fmt.Printf("%c is being added to First(%c)\n", rhs[i+1], nt);
				first = form.SetInsert(first, next);
				continue;
			}

			// Otherwise it is a non-terminal, so remove epsilon
			first = form.SetRemove(first, 'ε');

			// Then add first-set of non-terminal to current set
			//fmt.Printf("First(%c) also contains First(%c)...\n", nt, rhs[i+1]);
			ot_first, err := FirstSet(rhs[i+1], append(visited, nt), g);
			if err != nil {
				return first, err;
			}
			first = form.SetUnion(first, ot_first);
		}
	}
	return first, nil;
}


// Returns the follow-set for non terminal nt in grammar g
// [Note: Wasteful since discards follow sets it is not looking for ]
func FollowSet (nt rune, visited []rune, firsts[][]rune, g form.Item) ([]rune, error) {
	var ps []form.Prod		= g.Ps;
	var follow []rune		= []rune{};
	var j int;

	// If nt is not a non-terminal, return error
	if !form.IsNonTerminal(nt) {
		return follow, fmt.Errorf("FollowSet: %c is not a non-terminal!", nt);
	}

	// If a cycle is detected, return
	for _, s := range visited {
		if s == nt {
			return follow, nil;
		}
	}

	for _, p := range ps {
		length := len(p.Rhs);

		// Skip empty productions
		if length == 0 {
			continue;
		}

		// For all runes up the the second-last
		for j = 0; j < (length - 1); j++ {
			
			// Skip runes until occurrence of nt
			if p.Rhs[j] != nt {
				continue;
			}

			// If next rune is a terminal, append and cont
			if form.IsTerminal(p.Rhs[j+1]) {
				follow = form.SetInsert(follow, p.Rhs[j+1]);
				continue;
			}

			// Otherwise must be non-terminal, so fetch first-set
			ot_first := firsts[p.Rhs[j+1] % 26];

			// If doesn't contain epsilon, merge set and cont
			if !form.SetContains(ot_first, 'ε') {
				follow = form.SetUnion(follow, ot_first);
				continue;
			}

			// Otherwise has epsilon.
			// 1) Remove it from the first-set.
			// 2) Merge with next follow set, but only if
			// 3) No terminals succeed it.
			ot_first = form.SetRemove(ot_first, 'ε');
			ot_follow, err := FollowSet(p.Rhs[j+1], append(visited, nt), firsts, g);
			if err != nil {
				return follow, err;
			}
			follow = form.SetUnion(follow, form.SetUnion(ot_first, ot_follow));
		}

		// At last element: If non-terminal and rhs[j] == nt, add follow set too
		if form.IsNonTerminal(p.Rhs[j]) && p.Rhs[j] == nt && p.Lhs != nt {
			ot_follow, err := FollowSet(p.Lhs, append(visited, nt), firsts, g);
			if err != nil {
				return follow, err;
			}
			follow = form.SetUnion(follow, ot_follow);
		}
	}
	return follow, nil;
}