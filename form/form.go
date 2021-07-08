package form


/*
 * Form[at] package for grammar analysis
 * Defines internal productions, items and
 * item-sets (grammars). Supplies common
 * functions for general use. Designed
 * to work with the symtab and parse 
 * package
*/


import (
	"fmt"
	"strings"
	"grammars/parse"
	"grammars/symtab"
)


/*
 *******************************************************************************
 *                                  Constants                                  *
 *******************************************************************************
*/


// The epsilon token value
const Epsilon	=	0


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Production. Has a Lhs token index and a Rhs token index set
type Production struct {
	Lhs int;
	Rhs []int;
	Off int;
}


// Item. A set of productions
type Item []Production;


/*
 *******************************************************************************
 *                               Token Functions                               *
 *******************************************************************************
*/


// True if token is a non-terminal (t < 0)
func IsNonTerminal (t int) bool {
	return (t < 0);
}


// True if token is a terminal (s > 0)
func IsTerminal (t int) bool {
	return (t > 0);
}


// True if token is epsilon (s == 0)
func IsEpsilon (t int) bool {
	return (t == 0);
}


// General token equality function (for use with sets package)
func TokenCompare (a, b interface{}) bool {
	return (a.(int) == b.(int));
}


/*
 *******************************************************************************
 *                            Production Functions                             *
 *******************************************************************************
*/


// True if Production is an epsilon production
func (p *Production) EpsilonProduction () bool {
	return len(p.Rhs) == 0;
}


// String form of a production with optional dot notation
func (p *Production) String (dot bool, tab *symtab.SymTab) string {

	// Ignore error here, since it still returns placeholder lhs.
	lhs, _ := symtab.LookupID(p.Lhs, tab); 
	s := fmt.Sprintf("%s %s ", lhs, parse.DefineOperator);

	// Use epsilon symbol for empty productions
	if p.EpsilonProduction() {
		return s + "ε";
	}

	// Concatenate all tokens, including optional dot separator
	for i, j := range p.Rhs {
		if dot && i == p.Off {
			s += ".";
		}
		tok, _ := symtab.LookupID(j, tab);
		s = s + tok + " ";
	}
	
	return strings.TrimSuffix(s, " ");
}


/*
 *******************************************************************************
 *                               Item Functions                                *
 *******************************************************************************
*/


// True if Item is empty
func (i *Item) IsEmpty () bool {
	return len(*i) == 0;
}


// String form of an Item with optional dot notation
func (i *Item) String (dot bool, tab *symtab.SymTab) string {
	if i.IsEmpty() {
		return "Ø";
	}
	s := (*i)[0].String(dot, tab);
	for j := 1; j < len(*i); j++ {
		s = s + (*i)[j].String(dot, tab) + "\n";
	}
	return s;
}

