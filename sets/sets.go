package sets

import (
	"fmt"
	"strings"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Defines a Set of anything
type Set 	 	[]interface{}


// Defines a comparison function
type Compare	func(interface{}, interface{}) bool


// Defines a stringify function
type Stringify	func(interface{}) string


/*
 *******************************************************************************
 *                                Set Functions                                *
 *******************************************************************************
*/


// Returns set length
func (s Set) Len () int {
	return len(s);
}


// Returns set capacity
func (s Set) Cap () int {
	return cap(s);
}


// Returns true if set contains element
func (s Set) Contains (e interface{}, f Compare) bool {
	for _, x := range s {
		if f(e,x) {
			return true;
		}
	}
	return false;
}


// Inserts element into given set
func (s *Set) Insert (e interface{}, f Compare) {
	for _, x := range *s {
		if f(e, x) {
			return;
		}
	}
	*s = append(*s, e);
}


// Returns copy of given set (caution - applies to pointer sets)
func (s Set) Copy () Set {
	return s;
}


// Removes element from given set
func (s *Set) Remove (e interface{}, f Compare) {
	s_new := Set{};
	for _, x := range *s {
		if f(e, x) {
			continue;
		}
		s_new = append(s_new, x);
	}
	*s = s_new;
}


// Combines elements of two sets. Returns a new set
func Union (a, b *Set, f Compare) Set {
	s := *a;
	for _, e := range *b {
		s.Insert(e, f);
	}
	return s;
}


// Returns string form of set
func (s *Set) String (f Stringify) string {
	d := "{";
	l := s.Len();
	i := 0;
	if l == i {
		goto end;
	}
	for {
		tok := f((*s)[i]);
		d += tok;
		i++;
		if (i >= l) {
			break;
		}
		d += ",";
	}
	end:
	d += "}";
	return d;
}


/*
 *******************************************************************************
 *                          Default Sorting Functions                          *
 *******************************************************************************
*/


// String: Defines a comparison (equality) function
func CompareString (a, b interface{}) bool {
	return (strings.Compare(a.(string), b.(string)) == 0);
}


// String: Defines a stringify function
func ShowString(a interface{}) string {
	return a.(string);
}

// Int: Defines a comparison (equality) function
func CompareInt (a, b interface{}) bool {
	return (a.(int) == b.(int));
}


// Int: Defines a stringify function
func ShowInt (a interface{}) string {
	return fmt.Sprintf("%d", a.(int));
}


