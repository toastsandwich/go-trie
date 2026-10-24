package trie_test

import (
	"slices"
	"testing"

	"github.com/toastsandwich/go-trie"
)

func insertAll(words ...string) *trie.Trie {
	tr := trie.New()
	for _, w := range words {
		tr.Insert(w)
	}
	return tr
}

func nodeAt(t *testing.T, tr *trie.Trie, path string) *trie.Node {
	t.Helper()
	curr := tr.Root
	for _, r := range path {
		idx := slices.IndexFunc(curr.Children, func(n *trie.Node) bool {
			return n.Value == r
		})
		if idx < 0 {
			t.Fatalf("no node for %q at rune %q", path, r)
		}
		curr = curr.Children[idx]
	}
	return curr
}

func childValues(n *trie.Node) string {
	out := make([]rune, len(n.Children))
	for i, c := range n.Children {
		out[i] = c.Value
	}
	return string(out)
}

func TestNew(t *testing.T) {
	tr := trie.New()
	if tr == nil {
		t.Fatal("New() returned nil")
	}
	if tr.Root == nil {
		t.Fatal("root is nil")
	}
	if tr.Root.Value != -1 {
		t.Errorf("root value = %d, want -1", tr.Root.Value)
	}
	if tr.Root.IsEnd {
		t.Error("root IsEnd = true, want false")
	}
	if len(tr.Root.Children) != 0 {
		t.Errorf("root children = %d, want 0", len(tr.Root.Children))
	}
	if !tr.Search("") {
		t.Error(`Search("") = false, want true`)
	}
	if tr.Search("a") {
		t.Error(`Search("a") = true, want false`)
	}
}

func TestTriesAreIndependent(t *testing.T) {
	a := trie.New()
	b := trie.New()
	a.Insert("only-a")
	b.Insert("only-b")

	if !a.Search("only-a") || a.Search("only-b") {
		t.Errorf("a: Search(only-a)=%v Search(only-b)=%v", a.Search("only-a"), a.Search("only-b"))
	}
	if !b.Search("only-b") || b.Search("only-a") {
		t.Errorf("b: Search(only-b)=%v Search(only-a)=%v", b.Search("only-b"), b.Search("only-a"))
	}
}

func TestReset(t *testing.T) {
	tr := trie.New()
	tr.Insert("alpha")
	tr.Insert("beta")
	old := tr.Root

	tr.Reset()

	if tr.Root == nil {
		t.Fatal("root is nil after Reset")
	}
	if tr.Root == old {
		t.Error("Reset kept the previous root")
	}
	if tr.Root.Value != -1 {
		t.Errorf("root value = %d, want -1", tr.Root.Value)
	}
	if len(tr.Root.Children) != 0 {
		t.Errorf("root children = %d, want 0", len(tr.Root.Children))
	}
	if tr.Search("alpha") || tr.Search("beta") {
		t.Error("Search found a word inserted before Reset")
	}
	if !tr.Search("") {
		t.Error(`Search("") = false, want true`)
	}

	tr.Insert("alpha")
	if !tr.Search("alpha") {
		t.Error(`Search("alpha") = false after reinsert, want true`)
	}
	if tr.Search("beta") {
		t.Error(`Search("beta") = true, want false`)
	}

	tr.Reset()
	tr.Reset()
	if tr.Search("alpha") {
		t.Error(`Search("alpha") = true after a second Reset, want false`)
	}
}

// TestSearch checks path membership. A query matches when every rune exists
// as a child chain from the root. The empty string and proper prefixes match.
// Search does not consult Node.IsEnd.
func TestSearch(t *testing.T) {
	words := []string{
		"a",
		"apple",
		"bee",
		"book",
		"bookkeeper",
		"car",
		"card",
		"cat",
		"dog",
		"door",
	}

	tests := []struct {
		s    string
		want bool
	}{
		{s: "", want: true},
		{s: "a", want: true},
		{s: "ap", want: true},
		{s: "apple", want: true},
		{s: "apples", want: false},
		{s: "b", want: true},
		{s: "be", want: true},
		{s: "bee", want: true},
		{s: "beet", want: false},
		{s: "boo", want: true},
		{s: "book", want: true},
		{s: "bookkeeper", want: true},
		{s: "bookkeepers", want: false},
		{s: "box", want: false},
		{s: "c", want: true},
		{s: "ca", want: true},
		{s: "car", want: true},
		{s: "card", want: true},
		{s: "care", want: false},
		{s: "cat", want: true},
		{s: "catdog", want: false},
		{s: "d", want: true},
		{s: "do", want: true},
		{s: "dog", want: true},
		{s: "door", want: true},
		{s: "doors", want: false},
		{s: "dot", want: false},
		{s: "Hello", want: false},
		{s: "z", want: false},
	}

	tr := insertAll(words...)
	// A miss before a hit must leave the stored paths intact.
	if tr.Search("zzz") {
		t.Fatal(`Search("zzz") = true, want false`)
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := tr.Search(tt.s)
			if got != tt.want {
				t.Errorf("Search(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestSearch_SingleInsert(t *testing.T) {
	tests := []struct {
		name string
		in   string
		s    string
		want bool
	}{
		{name: "different word", in: "hello", s: "hi", want: false},
		{name: "exact word", in: "stumped", s: "stumped", want: true},
		{name: "proper prefix", in: "stumped", s: "stump", want: true},
		{name: "first rune", in: "stumped", s: "s", want: true},
		{name: "extra suffix", in: "stumped", s: "stumpedx", want: false},
		{name: "empty word", in: "", s: "", want: true},
		{name: "empty word has no rune path", in: "", s: "a", want: false},
		{name: "case differs", in: "Hello", s: "hello", want: false},
		{name: "case exact", in: "Hello", s: "Hello", want: true},
		{name: "repeated runes", in: "aaa", s: "aaa", want: true},
		{name: "repeated runes short prefix", in: "aaa", s: "aa", want: true},
		{name: "repeated runes too long", in: "aaa", s: "aaaa", want: false},
		{name: "interior space", in: "hello world", s: "hello world", want: true},
		{name: "space is part of the path", in: "hello world", s: "hello", want: true},
		{name: "trailing space missing", in: "hello world", s: "hello world ", want: false},
		{name: "null byte", in: "a\x00b", s: "a\x00b", want: true},
		{name: "null byte is not skipped", in: "a\x00b", s: "ab", want: false},
		{name: "punctuation", in: "rock'n'roll", s: "rock'n'roll", want: true},
		{name: "punctuation prefix", in: "rock'n'roll", s: "rock'", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trie.New()
			tr.Insert(tt.in)
			got := tr.Search(tt.s)
			if got != tt.want {
				t.Errorf("Insert(%q) Search(%q) = %v, want %v", tt.in, tt.s, got, tt.want)
			}
		})
	}
}

func TestSearch_InsertOrder(t *testing.T) {
	words := []string{"card", "car", "cat", "dog", "door", "apple", "book", "bookkeeper", "bee", "a"}
	queries := []string{
		"", "a", "ap", "apple", "apples",
		"c", "ca", "car", "card", "care", "cat",
		"d", "do", "dog", "door", "doors",
		"boo", "book", "bookkeeper", "bee", "beet",
		"Hello", "z",
	}

	forward := insertAll(words...)
	backward := trie.New()
	for i := len(words) - 1; i >= 0; i-- {
		backward.Insert(words[i])
	}

	for _, q := range queries {
		got := backward.Search(q)
		want := forward.Search(q)
		if got != want {
			t.Errorf("Search(%q) = %v with reverse insert, want %v", q, got, want)
		}
	}
}

func TestSearch_ShorterWordAfterLongerWord(t *testing.T) {
	tr := trie.New()
	tr.Insert("carpet")
	tr.Insert("car")

	for _, s := range []string{"c", "ca", "car", "carp", "carpe", "carpet"} {
		if !tr.Search(s) {
			t.Errorf("Search(%q) = false, want true", s)
		}
	}
	if tr.Search("cars") {
		t.Error(`Search("cars") = true, want false`)
	}
}

func TestSearch_ManySiblings(t *testing.T) {
	t.Run("full alphabet", func(t *testing.T) {
		tr := trie.New()
		for r := 'z'; r >= 'a'; r-- {
			tr.Insert(string(r))
			tr.Insert("p" + string(r))
		}

		if got, want := len(tr.Root.Children), 26; got != want {
			t.Fatalf("root children = %d, want %d", got, want)
		}
		p := nodeAt(t, tr, "p")
		if got, want := len(p.Children), 26; got != want {
			t.Fatalf("p children = %d, want %d", got, want)
		}

		for r := 'a'; r <= 'z'; r++ {
			letter := string(r)
			if !tr.Search(letter) {
				t.Errorf("Search(%q) = false, want true", letter)
			}
			withP := "p" + letter
			if !tr.Search(withP) {
				t.Errorf("Search(%q) = false, want true", withP)
			}
			// Single-letter words have no further children. "pp" is the
			// inserted word p + 'p', so it is not used as the miss probe.
			if tr.Search(letter + "!") {
				t.Errorf("Search(%q) = true, want false", letter+"!")
			}
		}
		if !tr.Search("p") {
			t.Error(`Search("p") = false, want true`)
		}
		for _, missing := range []string{"A", "`", "{", "P", "pA", "p{"} {
			if tr.Search(missing) {
				t.Errorf("Search(%q) = true, want false", missing)
			}
		}
	})

	t.Run("gaps", func(t *testing.T) {
		tr := trie.New()
		for r := 'a'; r <= 'z'; r += 2 {
			tr.Insert(string(r))
		}
		for r := 'a'; r <= 'z'; r++ {
			got := tr.Search(string(r))
			want := (r-'a')%2 == 0
			if got != want {
				t.Errorf("Search(%q) = %v, want %v", string(r), got, want)
			}
		}
		if tr.Search("ac") {
			t.Error(`Search("ac") = true, want false`)
		}
	})
}

func TestInsert_DuplicateDoesNotBranch(t *testing.T) {
	tr := trie.New()
	tr.Insert("go")
	tr.Insert("go")
	tr.Insert("go")

	if !tr.Search("go") {
		t.Fatal(`Search("go") = false, want true`)
	}
	if got, want := childValues(tr.Root), "g"; got != want {
		t.Fatalf("root children = %q, want %q", got, want)
	}
	g := nodeAt(t, tr, "g")
	if got, want := childValues(g), "o"; got != want {
		t.Fatalf("g children = %q, want %q", got, want)
	}
	if len(nodeAt(t, tr, "go").Children) != 0 {
		t.Fatal("duplicate insert added children under the terminal")
	}
}

func TestInsert_EmptyStringAddsNoChildren(t *testing.T) {
	tr := trie.New()
	tr.Insert("keep")
	tr.Insert("")
	tr.Insert("")

	if got, want := childValues(tr.Root), "k"; got != want {
		t.Fatalf("root children = %q, want %q", got, want)
	}
	if !tr.Search("keep") {
		t.Error(`Search("keep") = false, want true`)
	}
}

func TestChildrenAreSortedByRune(t *testing.T) {
	tr := insertAll("ba", "bc", "bb", "ab", "aa", "az", "a ")

	if got, want := childValues(tr.Root), "ab"; got != want {
		t.Fatalf("root children = %q, want %q", got, want)
	}
	if got, want := childValues(nodeAt(t, tr, "a")), " abz"; got != want {
		t.Errorf("a children = %q, want %q", got, want)
	}
	if got, want := childValues(nodeAt(t, tr, "b")), "abc"; got != want {
		t.Errorf("b children = %q, want %q", got, want)
	}
}

func TestIsEnd(t *testing.T) {
	t.Run("single letter", func(t *testing.T) {
		tr := insertAll("a")
		if !nodeAt(t, tr, "a").IsEnd {
			t.Error("IsEnd = false, want true")
		}
		if tr.Root.IsEnd {
			t.Error("root IsEnd = true, want false")
		}
	})

	t.Run("only the new terminal", func(t *testing.T) {
		tr := insertAll("trie")
		for _, prefix := range []string{"t", "tr", "tri"} {
			if nodeAt(t, tr, prefix).IsEnd {
				t.Errorf("node %q IsEnd = true, want false", prefix)
			}
		}
		if !nodeAt(t, tr, "trie").IsEnd {
			t.Error(`node "trie" IsEnd = false, want true`)
		}
	})

	t.Run("shorter word stays terminal after a longer word", func(t *testing.T) {
		tr := insertAll("car", "card", "book", "bookkeeper")
		for _, path := range []string{"car", "card", "book", "bookkeeper"} {
			if !nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = false, want true", path)
			}
		}
		for _, path := range []string{"ca", "boo", "bookk", "bookkee"} {
			if nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = true, want false", path)
			}
		}
	})

	t.Run("shorter word after a longer word", func(t *testing.T) {
		tr := trie.New()
		tr.Insert("carpet")
		tr.Insert("card")
		tr.Insert("car")
		tr.Insert("bookkeeper")
		tr.Insert("book")

		for _, path := range []string{"car", "card", "carpet", "book", "bookkeeper"} {
			if !nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = false, want true", path)
			}
		}
		for _, path := range []string{"c", "ca", "carp", "boo", "bookk", "bookkee"} {
			if nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = true, want false", path)
			}
		}
		if got, want := childValues(nodeAt(t, tr, "car")), "dp"; got != want {
			t.Errorf("car children = %q, want %q", got, want)
		}
	})

	t.Run("multibyte terminal", func(t *testing.T) {
		tr := insertAll("café", "éx", "🙂", "a🙂b", "こんにちは")
		ends := []string{"café", "éx", "🙂", "a🙂b", "こんにちは"}
		for _, path := range ends {
			if !nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = false, want true", path)
			}
		}
		open := []string{"c", "ca", "caf", "é", "a", "a🙂", "こん", "こんにち"}
		for _, path := range open {
			if nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = true, want false", path)
			}
		}
	})

	t.Run("shorter multibyte word after a longer word", func(t *testing.T) {
		tr := insertAll("caféteria", "café")
		if !nodeAt(t, tr, "café").IsEnd {
			t.Error(`node "café" IsEnd = false, want true`)
		}
		if !nodeAt(t, tr, "caféteria").IsEnd {
			t.Error(`node "caféteria" IsEnd = false, want true`)
		}
		for _, path := range []string{"c", "ca", "caf", "cafét"} {
			if nodeAt(t, tr, path).IsEnd {
				t.Errorf("node %q IsEnd = true, want false", path)
			}
		}
	})
}

func TestUnicodePaths(t *testing.T) {
	const zwj = "👨‍💻" // man, zero-width joiner, personal computer

	tr := insertAll("こんにちは", "🙂", "café", "a🙂b", zwj)

	tests := []struct {
		s    string
		want bool
	}{
		{s: "こんにちは", want: true},
		{s: "こん", want: true},
		{s: "こんにちは!", want: false},
		{s: "🙂", want: true},
		{s: "☺", want: false},
		{s: "café", want: true},
		{s: "cafe", want: false},
		{s: "caf", want: true},
		{s: "a🙂b", want: true},
		{s: "a🙂", want: true},
		{s: "ab", want: false},
		{s: zwj, want: true},
		{s: "👨", want: true},
		{s: "👨💻", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := tr.Search(tt.s)
			if got != tt.want {
				t.Errorf("Search(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}

	n := nodeAt(t, tr, "café")
	if n.Value != 'é' {
		t.Errorf("terminal value = %q, want é", n.Value)
	}
	if got, want := childValues(nodeAt(t, tr, "caf")), "é"; got != want {
		t.Errorf("children under caf = %q, want %q", got, want)
	}
}
