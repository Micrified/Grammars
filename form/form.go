package form

/*
 * Form[at] package for grammar analysis.
 * Defines productions, items, item-sets, and
 * common functions that are applied to them.
*/


import (
	"fmt"
	"io"
	"bufio"
	"regexp"
	"strings"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Production: Lhs ::= Rhs
type Prod struct {
	Off int;
	Lhs rune;
	Rhs []rune;
}

// Item: Set of productions, each with individual offsets
type Item struct {
	Ps []Prod;
}


/*
 *******************************************************************************
 *                       Functions on Rune and Rune Sets                       *
 *******************************************************************************
*/


// True if rune is a non-terminal
func IsNonTerminal (r rune) bool {
	return (r >= 'A' && r <= 'Z');
}


// True if rune is a terminal
func IsTerminal (r rune) bool {
	return !IsNonTerminal(r) && (r >= '!' && r <= '~');
}


// True if rune is contained in slice
func SetContains (set []rune, r rune) bool {
	for _, x := range set {
		if x == r {
			return true;
		}
	}
	return false;
}


// Inserts rune into set. Discards if already present
func SetInsert (set []rune, r rune) []rune {
	if SetContains(set, r) {
		return set;
	}
	return append(set, r);
}


// Removes a rune from a set (checks all elements in case)
func SetRemove (set []rune, r rune) []rune {
	var filtered []rune = []rune{};
	for _, x := range set {
		if x == r {
			continue;
		}
		filtered = append(filtered, x);
	}
	return filtered;
}


// Returns the union of two rune sets
func SetUnion (a, b []rune) []rune {
	for _, x := range b {
		a = SetInsert(a, x);
	}
	return a;
}


// Returns string form of a set
func SetToString (set []rune) string {
	s := "{";
	l := len(set);
	i := 0;
	if l == i {
		goto end;
	}
	for {
		s += fmt.Sprintf("%c", set[i]);
		i++;
		if i >= l {
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
func (p *Prod) String (dot bool) string {
	s := fmt.Sprintf("%c -> ", p.Lhs);
	if p.Epsilon() {
		s += "ε";
		goto end;
	}
	for i, r := range p.Rhs {
		if (dot && i == p.Off) {
			s += ".";
		}
		s += fmt.Sprintf("%c ", r);
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
func (item *Item) String (dot bool) string {
	if item.IsEmpty() {
		return "Ø";
	}
	s := "";
	for _, p := range item.Ps {
		s += fmt.Sprintf("%s\n", p.String(dot));
	}
	return s;
}


/*
 *******************************************************************************
 *            Functions for parsing Items and Productions from text            *
 *******************************************************************************
*/


// Parse a production from a string
func ParseProduction (line string) (Prod, error) {
	var p Prod;

	// Validate format
	format := `^[ \t]*[A-Z][ \t]*->[ \t]*[-+*/a-zA-Z0-9() ]*[ \t]*[$]?[\n]?$`;
	match, err := regexp.MatchString(format, line);
	if err != nil {
		return p, err;
	}
	if match == false {
		return p, fmt.Errorf("Invalid production format: %q", line);
	}
	
	// Remove whitespace
	ws := " \t\n";
	for i := range ws {
		line = strings.Replace(line, string(ws[i]), "", -1);
	}

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
