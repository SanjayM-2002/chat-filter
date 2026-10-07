package trie

// TrieNode holds 26 links, one per letter 'a'..'z'.
type TrieNode struct {
	links [26]*TrieNode
	isEnd bool
}

func (n *TrieNode) ContainsKey(ch byte) bool  { return n.links[ch-'a'] != nil }
func (n *TrieNode) GetNode(ch byte) *TrieNode { return n.links[ch-'a'] }
func (n *TrieNode) CheckEnd() bool            { return n.isEnd }
func (n *TrieNode) SetEnd()                   { n.isEnd = true }
func (n *TrieNode) Put(ch byte, node *TrieNode) {
	n.links[ch-'a'] = node
}

type Trie struct {
	root *TrieNode
}

func New() *Trie {
	return &Trie{root: &TrieNode{}}
}

func (t *Trie) Root() *TrieNode { return t.root }

// Insert adds a word. The word must already be normalized (only 'a'..'z').
func (t *Trie) Insert(word string) {
	node := t.root
	for i := 0; i < len(word); i++ {
		ch := word[i]
		if !node.ContainsKey(ch) {
			node.Put(ch, &TrieNode{})
		}
		node = node.GetNode(ch)
	}
	node.SetEnd()
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for i := 0; i < len(word); i++ {
		if !node.ContainsKey(word[i]) {
			return false
		}
		node = node.GetNode(word[i])
	}
	return node.CheckEnd()
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for i := 0; i < len(prefix); i++ {
		if !node.ContainsKey(prefix[i]) {
			return false
		}
		node = node.GetNode(prefix[i])
	}
	return true
}
