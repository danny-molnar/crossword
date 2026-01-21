package tools

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type PatternResult struct {
	Word string
}

func (wl *Wordlist) PatternMatch(pattern string, length int) ([]PatternResult, error) {
	p := strings.TrimSpace(pattern)
	if p == "" {
		return nil, fmt.Errorf("pattern empty")
	}
	if length <= 0 {
		// Infer from pattern if it uses ? as single-char wildcard.
		length = len([]rune(p))
	}

	// Convert simple crossword pattern to regex:
	// A?I?E => ^A.I.E$
	// Accept '.' already if user gives regex-like pattern; escape everything else.
	reg := "^" + regexp.QuoteMeta(p) + "$"
	reg = strings.ReplaceAll(reg, `\?`, ".")
	re, err := regexp.Compile(reg)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	candidates := wl.ByLen[length]
	out := make([]PatternResult, 0, 32)
	for _, w := range candidates {
		if re.MatchString(strings.ToLower(w)) {
			out = append(out, PatternResult{Word: w})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Word < out[j].Word })
	return out, nil
}
