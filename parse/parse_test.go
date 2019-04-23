package parse 

import (
	"fmt"
	"testing"
	"io"
	"strings"
)


/*
 *******************************************************************************
 *                        Tests for WhiteSpace Stripper                        *
 *******************************************************************************
*/


func TestDropWhiteSpaceOnWhiteSpace (t *testing.T) {
	s := "  \n\n\t\n\r \r \f\t\t\v ";
	rs := []rune(s);
	dropWhiteSpace(&rs);
	if len(rs) != 0 {
		t.Errorf("Not all whitespace dropped!");
	}
}

func TestDropWhiteSpaceNoLeadWhiteSpace (t *testing.T) {
	s := "H \n\t\t";
	rs := []rune(s);
	pre_len := len(rs);
	dropWhiteSpace(&rs);
	post_len := len(rs);
	if pre_len != post_len {
		t.Errorf("No whitespace should have been removed!");
	}
}

func TestDropWhiteSpaceNotTooMuch (t *testing.T) {
	s := "  \tx helo";
	rs := []rune(s);
	pre_len := len(rs);
	dropWhiteSpace(&rs);
	post_len := len(rs);
	if (pre_len - post_len) != 3 {
		t.Errorf("Should have removed 3 runes, but removed %d\n", pre_len - post_len);
	}
}


/*
 *******************************************************************************
 *                           Test for Non-Terminals                            *
 *******************************************************************************
*/


func TestParseEmptyNonTerminal (t *testing.T) {
	s := "";
	rs := []rune(s);
	
	_, err := parseNonTerminal(&rs);
	if err != io.EOF {
		t.Errorf("Parsing empty non-terminal must return EOF");
	}
}

func TestParseNonTerminal (t *testing.T) {
	s := "      \tXx' xX";
	rs := []rune(s);
	
	tok, err := parseNonTerminal(&rs);
	if err != nil {
		t.Errorf("NonTerminal should not raise error");
	}
	if strings.Compare(tok.Name, "Xx'") != 0 {
		t.Errorf("Expected Xx' as non-terminal!");
	}
	if tok.IsNT == false {
		t.Errorf("NonTerminal should have type IsNT true");
	}
}

func TestParseInvalidNonTerminal (t *testing.T) {
	s := "      \t1Xx' xX";
	rs := []rune(s);
	
	_, err := parseNonTerminal(&rs);
	if err == nil {
		t.Errorf("Invalid non-terminal should not parse without error");
	}
}

func TestParseNonTerminal2 (t *testing.T) {
	s := "A";
	rs := []rune(s);
	
	tok, err := parseNonTerminal(&rs);
	if err != nil {
		t.Errorf("NonTerminal should not raise error");
	}
	if strings.Compare(tok.Name, "A") != 0 {
		t.Errorf("Expected A as non-terminal!");
	}
	if tok.IsNT == false {
		t.Errorf("NonTerminal should have type IsNT true");
	}
}

func TestParseNonTerminalRange (t *testing.T) {
	AtoZ := []rune{ 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
					'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
					'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 
					'Y', 'Z'};
	atoz := []rune{ 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
					'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
					'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 
					'y', 'z'};

	all := append(append(AtoZ, atoz...), '\'');
	
	// All non-terminals must start with a capitol letter
	for _, a := range AtoZ {
		
		// All legal combinations of length 3 will be tested
		for _, b := range all {
			for _, c := range all {
				rs := []rune{a,b,c};
				_, err := parseNonTerminal(&rs);
				if err != nil {
					t.Errorf("%s is a non-terminal but got err!", string(rs));
				}
			}
		} 
	}
}


/*
 *******************************************************************************
 *                             Test for Operators                              *
 *******************************************************************************
*/


func TestParseDefineOperator (t *testing.T) {
	if len(DefineOperator) == 0 {
		t.Errorf("Def-op should not have zero length!");
	}
	s := "\t\t\r\n" + DefineOperator;
	rs := []rune(s);
	err := parseDefineOperator(&rs);

	if err != nil {
		t.Errorf("Failed to parse def-op: %s correctly!", DefineOperator);
	}
}

func TestNotParseDefineOperator (t *testing.T) {
	if len(DefineOperator) == 0 {
		t.Errorf("Def-op should not have zero length!");
	}
	rs := []rune(DefineOperator);
	rs[0] = rs[0] + 1;
	err := parseDefineOperator(&rs);
	if err == nil {
		t.Errorf("Mutated def-op should have raised error!");
	}
}


/*
 *******************************************************************************
 *                             Tests for Terminals                             *
 *******************************************************************************
*/


func TestParseEmptyTerminal (t *testing.T) {
	s := "";
	rs := []rune(s);
	
	_, err := parseTerminal(&rs);
	if err != io.EOF {
		t.Errorf("Parsing empty terminal must return EOF");
	}
}

func TestParseTerminal (t *testing.T) {
	s := "      \tpx' xX";
	rs := []rune(s);
	
	tok, err := parseTerminal(&rs);
	if err != nil {
		t.Errorf("Terminal should not raise error");
	}
	if strings.Compare(tok.Name, "px'") != 0 {
		t.Errorf("Expected px' as non-terminal!");
	}
	if tok.IsNT == true {
		t.Errorf("Terminal should have type IsNT false");
	}
}

func TestParseInvalidTerminal (t *testing.T) {
	s := "      \tXx' xX";
	rs := []rune(s);
	
	_, err := parseTerminal(&rs);
	if err == nil {
		t.Errorf("Invalid terminal should not parse without error");
	}
}

func TestParseTerminal2 (t *testing.T) {
	s := "*()";
	rs := []rune(s);
	
	tok, err := parseTerminal(&rs);
	if err != nil {
		t.Errorf("Terminal should not raise error");
	}
	if strings.Compare(tok.Name, "*()") != 0 {
		t.Errorf("Expected *() as terminal!");
	}
	if tok.IsNT == true {
		t.Errorf("Terminal should have type IsNT false");
	}
}

func TestParseTerminalRange (t *testing.T) {
	AtoZ := []rune{ 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
					'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
					'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 
					'Y', 'Z'};
	atoz := []rune{ 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
					'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
					'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 
					'y', 'z'};
	misc := []rune{ '|', '!', '@', '#', '$', '%', '^', '&',
					'^', '&', '*', '(', ')', '_', '-', '+',
					'=', '{', '[', ']', '}', ':', ';', '\'',
					'"', '|', '\\', '~', '`', ',', '<', '.',
					'>', '/', '?', '0', '1', '2', '3', '4', 
					'5', '6', '7', '8', '9'};

	norm := append(atoz, misc...);
	all := append(norm, AtoZ...);
	
	// All terminals must not start with a capitol letter
	for _, a := range norm {
		
		// All legal combinations of length 3 will be tested
		for _, b := range all {
			for _, c := range all {
				rs := []rune{a,b,c};
				_, err := parseTerminal(&rs);
				if err != nil {
					t.Errorf("%s is a terminal but got err!", string(rs));
				}
			}
		} 
	}
}


/*
 *******************************************************************************
 *                            Test for Productions                             *
 *******************************************************************************
*/


func TestParseNonEmptyProduction (t *testing.T) {
	g := fmt.Sprintf("E %s T T'", DefineOperator);
	rs := []rune(g);

	p, err := parseProduction(&rs);

	if err != nil {
		t.Errorf("Production %s is valid but raised error!", g);
	}
	if strings.Compare(p.Lhs.Name, "E") != 0 {
		t.Errorf("Lhs token name should be E but is %s!", p.Lhs.Name);
	}
	if p.Lhs.IsNT != true {
		t.Errorf("Lhs should always be a non-terminal!");
	}
	if len(p.Rhs) != 2 {
		t.Errorf("Expected Rhs to have two tokens but got %d!", len(p.Rhs));
	}
}

func TestParseEmptyProduction (t *testing.T) {
	g := fmt.Sprintf("E %s", DefineOperator);
	rs := []rune(g);

	p, err := parseProduction(&rs);

	if err != nil {
		t.Errorf("Production %s is valid but raised error!", g);
	}
	if strings.Compare(p.Lhs.Name, "E") != 0 {
		t.Errorf("Lhs token name should be E but is %s!", p.Lhs.Name);
	}
	if p.Lhs.IsNT != true {
		t.Errorf("Lhs should always be a non-terminal!");
	}
	if len(p.Rhs) != 0 {
		t.Errorf("Expected Rhs to have zero tokens but got %d!", len(p.Rhs));
	}
}

func TestParseInvalidProduction (t *testing.T) {
	g := fmt.Sprintf("e %s e + e", DefineOperator);
	rs := []rune(g);

	_, err := parseProduction(&rs);

	if err == nil {
		t.Errorf("Production %s is invalid but raised no error!", g);
	}
}

func TestParseEmptyProduction2 (t *testing.T) {
	g := "";
	rs := []rune(g);

	_, err := parseProduction(&rs);

	if err != io.EOF {
		t.Errorf("Production %s should have returned error io.EOF", g);
	}
}


/*
 *******************************************************************************
 *                             Tests for Grammars                              *
 *******************************************************************************
*/


func TestParseWellFormedGrammarFromString (t *testing.T) {
	g := "E  ->   T E'\n"    +
		 "E' -> + T E'\n"    + 
		 "E' -> - T E'\n"    +
		 "E' ->\n"           +
		 "T  -> F T'\n"      +
		 "T' -> * T'\n"      +
		 "T' -> / T'\n"      +
		 "T' ->\n"           +
		 "F  -> id";
	
	_, err := ParseGrammarFromString(g);

	if err != nil {
		t.Errorf("Grammar:\n%s\n Is valid but returned error!", g);
	}
}

func TestParseOddGrammarFromString (t *testing.T) {
	g := "E  ->   T E'\n"    +
		 "E' -> + T E'\n"    + 
		 "E' -> - T E'\n"    +
		 "\n\n"				 +
		 "E' ->\n"           +
		 "T  -> F T'\n"      +
		 "T' -> * T'\n"      +
		 "T' ->\t / T'\n"      +
		 "\n\t\t\t\r"        +
		 "T' ->\n"           +
		 "F  -> id";
	
	_, err := ParseGrammarFromString(g);

	if err != nil {
		t.Errorf("Grammar:\n%s\n Is valid but returned error: %s", g, err);
	}
}