package trie_test

import (
	"testing"

	"github.com/toastsandwich/go-trie"
)

func TestTrie_Search(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   string
		s    string
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test#1",
			in:   "hello",
			s:    "hi",
			want: false,
		},
		{
			name: "test#2",
			in:   "stumped",
			s:    "stumped",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trie.New()
			tr.Insert(tt.in)
			got := tr.Search(tt.s)
			if got != tt.want {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
