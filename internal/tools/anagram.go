package tools

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type AnagramResult struct {
	Word  string
	Score int // placeholder (frequency scoring can come later)
}

func (wl *Wordlist) Anagrams(letters string, length int) ([]AnagramResult, error) {
	in := normalizeLetters(letters)
	if in == "" {
		return nil, fmt.Errorf("letters empty")
	}

	if length > 0 && length != len([]rune(in)) {
		// MVP is exact anagram: length must match
		return []AnagramResult{}, nil
	}

	sig := signature(in)
	words := wl.BySig[sig]

	results := make([]AnagramResult, 0, len(words))
	for _, w := range words {
		results = append(results, AnagramResult{
			Word:  w,
			Score: 0,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		// deterministic: alpha, later score desc then alpha
		return results[i].Word < results[j].Word
	})

	return results, nil
}

func normalizeLetters(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
		default:
			// ignore spaces/punct
		}
	}
	return b.String()
}
