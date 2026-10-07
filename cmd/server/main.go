package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/SanjayM-2002/chat-filter/internal/filter"
	"github.com/SanjayM-2002/chat-filter/internal/wordlist"
	"github.com/SanjayM-2002/chat-filter/internal/ws"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	wordsPath := flag.String("words", "data/words.json", "path to word list JSON")
	flag.Parse()

	words, err := wordlist.LoadFile(*wordsPath)
	if err != nil {
		log.Fatal(err)
	}
	f := filter.New()
	f.AddWords(words)
	log.Printf("loaded %d words from %s", len(words), *wordsPath)

	hub := ws.NewHub()

	http.Handle("/ws", ws.Handler(hub, f))
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("listening on http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
