package form

/*
 * Form[at] package for grammar analysis.
 * Defines productions, items, item-sets, and
 * common functions that are applied to them.
 * Designed to work with the symtab package.
*/


import (
	"fmt"
	"io"
	"bufio"
	"regexp"
	"strings"
	"symtab"
	"unicode"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Production: Lhs ::= Rhs
type Prod struct {
	Off int;
	Lhs int;
	Rhs []int;
}

// Item: Set of productions, each with individual offsets
type Item struct {
	Ps []Prod;
}


/*
 *******************************************************************************
 *                     Functions on Symbol and Symbol Sets                     *
 *******************************************************************************
*/


// True if symbol is a non-terminal (s < 0)
func IsNonTerminal (s int) bool {
	return (s < 0);
}


// True if symbol is a terminal (s > 0)
func IsTerminal (s int) bool {
	return (s > 0);
}


// True if symbol is contained in slice
func SetContains (set []int, s int) bool {
	for _, x := range set {
		if x == s {
			return true;
		}
	}
	return false;
}


// Inserts symbol into set. Discards if already present
func SetInsert (set []int, s int) []int {
	if SetContains(set, s) {
		return set;
	}
	return append(set, s);
}


// Removes a symbol from a set (checks all elements in case)
func SetRemove (set []int, s int) []int {
	var filtered []int = []int{};
	for _, x := range set {
		if x == s {
			continue;
		}
		filtered = append(filtered, x);
	}
	return filtered;
}


// Returns the union of two sets
func SetUnion (a, b []int) []int {
	for _, x := range b {
		a = SetInsert(a, x);
	}
	return a;
}


// Returns string form of a set
func SetToString (set []int, tab *symtab.SymTab) string {
	s := "{";
	l := len(set);
	i := 0;
	if l == i {
		goto end;
	}
	for {
		id, err := symtab.LookupID(set[i], tab);
		s += id;
		i++;
		if (i >= l || err != nil) {
			break;
		}
		s += ",";
	}
	end:
	s += "}";
	return s;
}


/*
 *******************************************************************************
 *                          Functions on Productions                           *
 *******************************************************************************
*/


// True if Production is epsilon
func (p *Prod) Epsilon () bool {
	return len(p.Rhs) == 0;
}


// String form of a production. If dot is true, it is shown in production
func (p *Prod) String (dot bool, tab *symtab.SymTab) string {
	lhs, _ := symtab.LookupID(p.Lhs, tab); // If err, use placeholder returned
	s := fmt.Sprintf("%s -> ", lhs);
	if p.Epsilon() {
		s += "ε";
		goto end;
	}
	for i, r := range p.Rhs {
		if (dot && i == p.Off) {
			s += ".";
		}
		sym, _ := symtab.LookupID(r, tab);
		s += fmt.Sprintf("%s ", sym);
	}

	end:
	return s;
}


/*
 *******************************************************************************
 *                       Functions on Item and Item Sets                       *
 *******************************************************************************
*/


// True if Item has an empty set of productions
func (item *Item) IsEmpty () bool {
	return len(item.Ps) == 0;
}


// String form of an Item. If dot is true, it is shown in all productions
func (item *Item, tab *symtab.SymTab) String (dot bool) string {
	if item.IsEmpty() {
		return "Ø";
	}
	s := "";
	for _, p := range item.Ps {
		s += fmt.Sprintf("%s\n", p.String(dot, tab));
	}
	return s;
}


/*
 *******************************************************************************
 *            Functions for parsing Items and Productions from text            *
 *******************************************************************************
*/


// Parse a production from a string. Registers symbols in the given table
func ParseProduction (line string, tab *symtab.SymTab) (Prod, error) {
	var p Prod;

	// Validate format
	format := `^[ \t]*[A-Z][a-zA-Z']*[ \t]*->[ \t]*[-+*/^a-zA-Z0-9()' ]*[ \t]*[$]?[\n]?$`;
	match, err := regexp.MatchString(format, line);
	if err != nil {
		return p, err;
	}
	if match == false {
		return p, fmt.Errorf("Invalid production format: %q", line);
	}

	// Returns next word, and remaining runes. If no word, ptr is nil
	nextword := func(rs []rune) (string *, []rune) {
		i := 0; j := 0; l := len(rs);

		// drop whitespace until next non-whitespace character
		for _, r := range rs {
			if !unicode.IsSpace(r) {
				break;
			}
			i++;
		}
	
		if (i == l) {
			return nil, []rune{};
		} else {
			rs = rs[i:];
		}

		// collect characters until next whitespace character
		for _, r := range rs {
			if unicode.IsSpace(r) {
				break;
			}
			j++;
		}
		s := string(rs[:j]);
		return &s, rs[j:];
	}

	// Returns next word and runes read. If no word, ptr is nil. 
	nextword := func(rs []rune) (string *, int) {
		i := 0;

		// drop whitespace until next non-whitespace character
		for _, r := range rs {
			if !unicode.IsSpace(r) {
				break;
			}
			i++;
		}

		
		// grab until next whitespace character
		for _, r := range 
	};
	
	// Remove whitespace
	ws := " \t\n";
	for i := range ws {
		line = strings.Replace(line, string(ws[i]), "", -1);
	}

	// Parse tokens

	// Convert line to runes
	runes := []rune(line);
	return Prod{Off: 0, Lhs: runes[0], Rhs: runes[3:]}, nil;
}


// Parse an item from a readable source
func ParseItem (r *bufio.Reader) (Item, error) {
	var line string;
	var err error;
	var p Prod;
	var ps []Prod;

	for {
		if line, err = r.ReadString('\n'); err != nil {
			break;
		}
		if p, err = ParseProduction(line); err != nil {
			break;
		}
		ps = append(ps, p);
	}
	i := Item{Ps: ps};
	if err != io.EOF {
		return i, err;
	} else {
		return i, nil;
	}
}
