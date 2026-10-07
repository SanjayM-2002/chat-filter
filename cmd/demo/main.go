package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/SanjayM-2002/chat-filter/internal/filter"
	"github.com/SanjayM-2002/chat-filter/internal/wordlist"
)

func main() {
	wordsPath := flag.String("words", "data/words.json", "path to word list JSON")
	flag.Parse()

	words, err := wordlist.LoadFile(*wordsPath)
	if err != nil {
		log.Fatal(err)
	}

	f := filter.New()
	f.AddWords(words)
	fmt.Printf("loaded %d words from %s\n\n", len(words), *wordsPath)

	tests := []string{
		"mondays are shit",
		"mondays are SHIT!",
		"this is sh1t",
		"$hit happens",
		"shiiiiit man",
		"s.h.i.t",
		"s h i t",
		"f-u-c-k you",
		"you ass",
		"I passed the class",
		"assistant manager",
		"this hit hard",
		"you son of a bitch",
		"fuckfuck",
	}

	for _, t := range tests {
		fmt.Printf("%-24s -> %s\n", t, f.Mask(t))
	}
}
