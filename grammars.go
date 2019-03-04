package main

import (
	"fmt"
	"grammars/form"
	"grammars/prop"
	"bufio"
	"os"
	"log"
)

func main () {
	var item form.Item;
	var err error;
	var reader *bufio.Reader;
	var shown []bool = make([]bool, 26);
	var firsts [][]rune = make([][]rune, 26);

	reader = bufio.NewReader(os.Stdin);

	item, err = form.ParseItem(reader);

	if err != nil {
		log.Fatal(err);
	}

	fmt.Println("Ok!");
	s := item.String(false);
	fmt.Printf("%s", s);

	fmt.Println("First-Sets:");
	for i, p := range item.Ps {
		lhs := p.Lhs;
		if (shown[lhs % 26] == true) {
			continue;
		}
		first, err := prop.FirstSet(p.Lhs, []rune{}, item);
		if err != nil {
			log.Fatal(err);
		}
		fmt.Printf("%d. First(%c) := %s\n", i, lhs, form.SetToString(first));
		firsts[lhs % 26] = first;
		shown[lhs % 26] = true;
	}

	// Reset.
	shown = make([]bool, 26);

	fmt.Println("Follow-Sets:");
	for i, p := range item.Ps {
		lhs := p.Lhs;
		if (shown[lhs % 26] == true) {
			continue;
		}

		follow, err := prop.FollowSet(lhs, []rune{}, firsts, item);
		if err != nil {
			log.Fatal(err);
		}
		fmt.Printf("%d. Follow(%c) := %s\n", i, lhs, form.SetToString(follow));
		shown[lhs % 26] = true;
	}

}