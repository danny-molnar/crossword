package tools

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Wordlist struct {
	Words   []string
	ByLen   map[int][]string    // length -> words
	BySig   map[string][]string // sorted letters -> words (exact anagrams)
	ByLenLC map[int][]string    // length -> lowercase words (handy for matching)
}

func LoadWordlist(path string) (*Wordlist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open wordlist: %w", err)
	}
	defer f.Close()

	wl := &Wordlist{
		ByLen:   make(map[int][]string),
		BySig:   make(map[string][]string),
		ByLenLC: make(map[int][]string),
	}

	seen := map[string]bool{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		w := normalizeWord(line)
		if w == "" {
			continue
		}
		if seen[w] {
			continue
		}
		seen[w] = true

		wl.Words = append(wl.Words, w)

		l := letterCount(w)
		if l <= 0 {
			continue
		}
		wl.ByLen[l] = append(wl.ByLen[l], w)
		wl.ByLenLC[l] = append(wl.ByLenLC[l], strings.ToLower(w))

		sig := signature(w)
		wl.BySig[sig] = append(wl.BySig[sig], w)
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan wordlist: %w", err)
	}

	// Stable ordering for predictable output
	sort.Strings(wl.Words)
	for k := range wl.ByLen {
		sort.Strings(wl.ByLen[k])
	}
	for k := range wl.BySig {
		sort.Strings(wl.BySig[k])
	}
	return wl, nil
}

func normalizeWord(s string) string {
	// Keep letters/digits; allow internal apostrophes/hyphens if you want later.
	// For MVP: strip everything except letters/digits.
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
		default:
			// drop punctuation/whitespace
		}
	}
	return b.String()
}

func letterCount(s string) int {
	n := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			n++
		}
	}
	return n
}

func signature(s string) string {
	rs := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			rs = append(rs, unicode.ToLower(r))
		}
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i] < rs[j] })
	return string(rs)
}
