package main

import "fmt"

type Node struct {
	val   rune
	leafs map[rune]*Node
}

var leets map[rune][]rune

func main() {
	words := []string{"hello", "wassup", "cookie"}
	leets = map[rune][]rune{'e': {'3'}, 'o': {'0'}, 'l': {'7', '1'}}

	trie := Node{val: '_', leafs: make(map[rune]*Node)}
	for _, word := range words {
		addWord(&trie, word)
	}
	fmt.Println(isInTrie(&trie, "hello"), isInTrie(&trie, "noooo"), isInTrie(&trie, "h3llo"), isInTrie(&trie, "h3ll0"), isInTrie(&trie, "h3770"))
}

// addWord adds a new word to the trie
func addWord(trie *Node, word string) {
	current := trie
	for i, v := range word {
		if leets[v] != nil {
			for _, c := range leets[v] {
				// fmt.Println(string(c) + word[i+1:])
				// newNode := Node{val: c, leafs: make(map[rune]*Node)}
				// current.leafs[c] = &newNode
				addWord(current, string(c)+word[i+1:])
			}
		}
		if current.leafs[v] == nil {
			newNode := Node{val: v, leafs: make(map[rune]*Node)}
			current.leafs[v] = &newNode
			current = &newNode
		} else {
			current = current.leafs[v]
		}
	}
}

func isInTrie(trie *Node, word string) bool {
	current := trie
	for _, v := range word {
		if current.leafs[v] == nil {
			return false
		} else {
			current = current.leafs[v]
		}
	}
	return true
}
