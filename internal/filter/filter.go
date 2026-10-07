package filter

import "github.com/SanjayM-2002/chat-filter/internal/trie"

type AbuseFilter struct {
	trie *trie.Trie
}

func New() *AbuseFilter {
	return &AbuseFilter{trie: trie.New()}
}

// AddWord normalizes a raw word ("SH1T", "son of a bitch") and inserts it.
func (f *AbuseFilter) AddWord(raw string) {
	w := make([]byte, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		if n := normalize(raw[i]); n != 0 {
			w = append(w, n)
		}
	}
	if len(w) > 0 {
		f.trie.Insert(string(w))
	}
}

// AddWords inserts every word in the list.
func (f *AbuseFilter) AddWords(words []string) {
	for _, w := range words {
		f.AddWord(w)
	}
}

// matchAt walks the trie starting at text[start].
// Returns end index (exclusive) of the LONGEST bad word found, or -1.
func (f *AbuseFilter) matchAt(text string, start int) int {
	node := f.trie.Root()
	best := -1
	var last byte

	for j := start; j < len(text); j++ {
		c := normalize(text[j])

		if c == 0 {
			// not a letter: skip it if it's a separator (s.h.i.t), else stop
			if j > start && isSeparator(text[j]) {
				continue
			}
			break
		}

		if node.ContainsKey(c) {
			node = node.GetNode(c) // normal step down the trie
			last = c
		} else if c == last {
			// repeated char: shiiiit -> stay on the same node
		} else {
			break // dead end, no bad word continues this way
		}

		// full bad word here AND the word really ends here (not "assistant")
		if node.CheckEnd() && atBoundary(text, j+1) {
			best = j + 1
		}
	}
	return best
}

// Mask replaces every bad word in text with '*' (spaces inside a match are kept).
func (f *AbuseFilter) Mask(text string) string {
	out := []byte(text)
	i := 0
	for i < len(text) {
		// only try to match where a word starts (so "class" never matches "ass")
		if isWordChar(text[i]) && atBoundary(text, i-1) {
			if end := f.matchAt(text, i); end != -1 {
				for k := i; k < end; k++ {
					if text[k] != ' ' {
						out[k] = '*'
					}
				}
				i = end
				continue
			}
		}
		i++
	}
	return string(out)
}
