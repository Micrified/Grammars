package parse

import (
	"unicode"
	"errors"
	"io"
	"fmt"
	"strings"
	"bufio"
)

/*
 ********** Grammar Description *********
 * - NonTerminal:	[A-Z][a-zA-Z']*
 * - Terminal:		[^A-Z][^ \t\n]*
 * - Production: 	NT<DEF-OP>(NT|T)*
 * - Grammar:		NT+
 * (note: only printable symbols allowed)
 *
 ****************************************
*/


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Simple token. Has identifier and type (terminal or nonterminal)
type Token struct {
	Name 	string
	IsNT	bool
}


// Production. Has a RHS token and a LHS token set
type Production struct {
	Lhs		Token
	Rhs		[]Token
}


// Grammar. A set of productions
type Grammar []Production


// Symbol used for define operator
const DefineOperator	=	"->"


// Maximum line count for a grammar
const MaxLineCount		=	16


/*
 *******************************************************************************
 *                            Formatting Functions                             *
 *******************************************************************************
*/


// Converts Token to string descriptor
func (t *Token) String () string {
	return t.Name;
}


// Converts Production to string descriptor
func (p *Production) String () string {
	lhs := p.Lhs.String();
	s := fmt.Sprintf("%s %s ", lhs, DefineOperator);

	if len(p.Rhs) == 0 {
		return s + `ε`;
	}

	for _, tok := range p.Rhs {
		s = s + tok.String() + " ";
	}
	return strings.TrimSuffix(s, " ");
}


// Converts a Grammar to string descriptor
func (g *Grammar) String () string {
	s := "";
	for _, p := range *g {
		s = s + p.String() + "\n";
	}
	return s;
}


/*
 *******************************************************************************
 *                              Parsing Functions                              *
 *******************************************************************************
*/


// Removes whitespace from head of rune slice
func dropWhiteSpace (rs *[]rune) {
	i := 0;
	for _, r := range (*rs) {
		if !unicode.IsSpace(r) {
			break;
		}
		i++;
	}
	*rs = (*rs)[i:];
}


// Returns next non-terminal in slice. Otherwise returns error
func parseNonTerminal (rs *[]rune) (Token, error) {
	var tok Token = Token{"", true};
	dropWhiteSpace(rs);
	if len(*rs) == 0 {
		return tok, io.EOF;
	}
	if r := (*rs)[0]; !(unicode.IsLetter(r) && unicode.IsUpper(r)) {
		return tok, errors.New("Non-terminal must begin with uppercase letter!");
	}
	i := 1;
	for _, r := range (*rs)[1:] {
		if !(unicode.IsLetter(r) || r == '\'') {
			break;
		}
		i++;
	}
	tok.Name = string((*rs)[:i]);
	*rs = (*rs)[i:];
	return tok, nil;
}


// Returns nil if define operator in slice. Otherwise returns error
func parseDefineOperator (rs *[]rune) error {
	op := []rune(DefineOperator);
	dropWhiteSpace(rs);
	if len(*rs) == 0 {
		return io.EOF;
	}
	if len(*rs) < len(op) {
		return fmt.Errorf("Missing define operator: %s", DefineOperator);
	}
	for i, r := range op {
		if r != (*rs)[i] {
			return fmt.Errorf("Missing define operator: %s", DefineOperator);
		}
	}
	*rs = (*rs)[len(op):];
	return nil;
}


// Returns next terminal in slice. Otherwise returns error
func parseTerminal (rs *[]rune) (Token, error) {
	tok := Token{"", false};
	dropWhiteSpace(rs);
	if len(*rs) == 0 {
		return tok, io.EOF;
	}
	if r := (*rs)[0]; unicode.IsLetter(r) && unicode.IsUpper(r) {
		return tok, errors.New("Terminal cannot begin with uppercase letter!");
	}
	i := 1;
	for _, r := range (*rs)[1:] {
		if unicode.IsSpace(r) {
			break;
		}
		if !unicode.IsPrint(r) {
			return tok, fmt.Errorf("Invalid rune: %q", r);
		}
		i++;
	}
	tok.Name = string((*rs)[:i]);
	*rs = (*rs)[i:];
	return tok, nil;
}


// Parses a Production from the slice. Otherwise returns error
func parseProduction (rs *[]rune) (Production, error) {
	var p Production;
	var tok Token;
	var err error;

	// Parse RHS non-terminal
	tok, err = parseNonTerminal(rs);
	if err != nil {
		return p, err;
	}
	p = Production{Lhs: tok, Rhs: []Token{}};

	// Parse define operator
	if err = parseDefineOperator(rs); err != nil {
		return p, errors.New("Missing define operator!");
	}

	// Parse zero or more terminals or non-terminals
	for {
		tok, err = parseNonTerminal(rs);
		if err == io.EOF {
			break;
		}
		if err == nil {
			goto next;
		}
		tok, err = parseTerminal(rs);
		if err != nil {
			return p, err;
		}
		next:
		p.Rhs = append(p.Rhs, tok);
	}
	return p, nil;
}

// Parses string into slice of Productions
func ParseGrammarFromString (s string) (Grammar, error) {
	lines := strings.Split(s, "\n");
	grammar := []Production{};

	if len(lines) == 0 || len(lines) > MaxLineCount {
		return grammar, fmt.Errorf("Invalid production range: (0..%d]", MaxLineCount);
	}
	
	for n, line := range lines {
		rs := []rune(line);
		production, err := parseProduction(&rs);
		if err == io.EOF {
			continue;
		}
		if err != nil {
			return grammar, fmt.Errorf("Line %d: %s", n + 1, err);
		}
		grammar = append(grammar, production);
	}
	
	return grammar, nil;
}


// Parses buffered input into slice of Productions
func ParseGrammarFromReader (r *bufio.Reader) (Grammar, error) {
	grammar := []Production{};
	var line string; var prod Production; var err error; var n int;

	for n = 0; n < MaxLineCount; n++ {
		if line, err = r.ReadString('\n'); err != nil {
			break;
		}
		rs := []rune(line);

		// Here, we ignore io.EOF since the buffered reader will find it.
		if prod, err = parseProduction(&rs); err == io.EOF {
			continue;
		}
		if err != nil {
			break;
		}
		grammar = append(grammar, prod);
	}

	if (n == 0 || n > MaxLineCount) {
		return grammar, 
			fmt.Errorf("Invalid production range: (0..%d]", MaxLineCount);
	}

	if err != io.EOF {
		return grammar, fmt.Errorf("Line %d: %s", n + 1, err);
	}
	
	return grammar, nil;
}
		
