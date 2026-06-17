// Package moderation provides sensitive word filtering using AC automaton.
package moderation

import (
	"strings"
	"sync"
)

// Filter provides sensitive word detection.
type Filter struct {
	mu       sync.RWMutex
	root     *node
	wordList []string
}

type node struct {
	children map[rune]*node
	fail     *node
	isEnd    bool
	word     string
}

// NewFilter creates a new sensitive word filter with a built-in word list.
func NewFilter(words []string) *Filter {
	f := &Filter{
		root:     &node{children: make(map[rune]*node)},
		wordList: words,
	}
	f.build()
	return f
}

// build constructs the AC automaton.
func (f *Filter) build() {
	// Build trie
	for _, word := range f.wordList {
		f.insert(word)
	}
	// Build failure links (BFS)
	queue := make([]*node, 0)
	for _, child := range f.root.children {
		child.fail = f.root
		queue = append(queue, child)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for r, child := range current.children {
			queue = append(queue, child)
			failNode := current.fail
			for failNode != nil {
				if next, ok := failNode.children[r]; ok {
					child.fail = next
					break
				}
				failNode = failNode.fail
			}
			if child.fail == nil {
				child.fail = f.root
			}
		}
	}
}

func (f *Filter) insert(word string) {
	current := f.root
	for _, r := range word {
		if _, ok := current.children[r]; !ok {
			current.children[r] = &node{children: make(map[rune]*node)}
		}
		current = current.children[r]
	}
	current.isEnd = true
	current.word = word
}

// Contains returns true if the text contains any sensitive word.
func (f *Filter) Contains(text string) bool {
	return len(f.FindAll(text)) > 0
}

// FindAll returns all sensitive words found in the text.
func (f *Filter) FindAll(text string) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	text = strings.ToLower(text)
	found := make(map[string]bool)
	current := f.root

	for _, r := range text {
		for current != f.root {
			if _, ok := current.children[r]; ok {
				break
			}
			current = current.fail
		}
		if next, ok := current.children[r]; ok {
			current = next
		} else {
			current = f.root
		}
		// Check all pattern-ending nodes along the failure chain
		for temp := current; temp != f.root; temp = temp.fail {
			if temp.isEnd {
				found[temp.word] = true
			}
		}
	}

	result := make([]string, 0, len(found))
	for word := range found {
		result = append(result, word)
	}
	return result
}

// DefaultChineseWords returns a basic Chinese sensitive word list.
func DefaultChineseWords() []string {
	return []string{
		"赌博", "色情", "暴力", "毒品",
		"诈骗", "欺诈", "自杀", "自残",
		"违禁", "管制",
	}
}
