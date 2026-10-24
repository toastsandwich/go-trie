# go-trie

A trie that stores strings one rune at a time. Children of each node stay sorted by rune, and lookup uses binary search.

```
go get github.com/toastsandwich/go-trie@v0.0.2
```

The module requires Go 1.27.2.

## Usage

```go
package main

import (
	"fmt"

	"github.com/toastsandwich/go-trie"
)

func main() {
	t := trie.New()
	t.Insert("car")
	t.Insert("card")
	t.Insert("café")

	fmt.Println(t.Search("car"))  // true
	fmt.Println(t.Search("ca"))   // true: the path exists
	fmt.Println(t.Search("cat"))  // false
	fmt.Println(hasWord(t, "ca")) // false: "ca" was never inserted
	fmt.Println(hasWord(t, "café"))

	t.Reset()
	fmt.Println(t.Search("car")) // false
}

// hasWord reports whether word was inserted. Search only checks that the
// rune path exists, so a proper prefix also returns true.
func hasWord(t *trie.Trie, word string) bool {
	n := t.Root
	for _, r := range word {
		var next *trie.Node
		for _, child := range n.Children {
			if child.Value == r {
				next = child
				break
			}
		}
		if next == nil {
			return false
		}
		n = next
	}
	return n.IsEnd
}
```

## API

`New` returns an empty trie. The root value is `-1` and it is not a stored rune.

`Insert` adds one node per rune. The last rune of the word gets `IsEnd` set, including when that node already exists because a longer word was inserted first. Inserting the same word again does not add duplicate children. `Insert("")` leaves the trie unchanged.

`Search` walks the query one rune at a time and returns whether every rune is present. An empty query returns true. A proper prefix of an inserted word also returns true. `Search` does not read `IsEnd`.

`Reset` replaces the root and drops every inserted word.

`Node` exposes `Value`, `IsEnd`, and `Children`. `IsEnd` is true only for the final rune of a word passed to `Insert`.

Matching is case-sensitive and follows Go runes. `café` and `cafe` are different words. A sequence such as `👨‍💻` is several runes, so a prefix of that sequence still matches `Search`.
