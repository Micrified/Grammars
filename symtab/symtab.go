package symtab

/*
 * Defines a symbol table with registration and lookup routines. 
 * Returns an index for registered identifiers with O(1) lookup.
 * The routines allows two kinds of symbols to be registered:
 * (1) Terminals (returning index with value > 0)
 * (2) Non-terminals (returning index with value < 0)
 * They are stored in individual tables in the given SymTab table.
*/

import (
	"fmt"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// A symbol table (stores identifiers in two classes)
type SymTab struct {
	Ts  []string;		// Terminals
	NTs []string;		// Non-terminals
}


/*
 *******************************************************************************
 *                         Table Registration Routines                         *
 *******************************************************************************
*/


// Registers terminal identifier; returns index
func RegisterTerminal (id string, tab *SymTab) int {
	idx := len((*tab).Ts);
	(*tab).Ts = append((*tab).Ts, id);
	return (idx + 1);
}


// Registers non-terminal identifier; returns index
func RegisterNonTerminal (id string, tab *SymTab) int {
	idx := len((*tab).NTs);
	(*tab).NTs = append((*tab).NTs, id);
	return -(idx + 1);
}


/*
 *******************************************************************************
 *                           Table Lookup Routines                             *
 *******************************************************************************
*/


// Returns identifier for given index. Returns placeholder on error
func LookupID (idx int, tab *SymTab) (string, error) {
	str := "<NULL>";

	if (idx == 0) {
		goto bad_bounds;
	}

	if (idx > 0) {
		if (idx > len((*tab).Ts)) {
			goto bad_bounds;
		}
		return (*tab).Ts[idx - 1], nil;
	}

	if ((-idx) > len((*tab).NTs)) {
		goto bad_bounds;
	}

	return (*tab).NTs[(-idx) - 1], nil;
	
	bad_bounds:
	return str, fmt.Errorf("Index %d is invalid!", idx);
}


/*
 *******************************************************************************
 *                          Table Control Operations                           *
 *******************************************************************************
*/


// Purges the tables
func ResetSymTab (tab *SymTab) {
	(*tab).Ts = []string{};
	(*tab).NTs = []string{};
}
