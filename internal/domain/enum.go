package domain

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Enum represents an enumeration like:
// "3"        -> Parts=[3], Total=3
// "3,5"      -> Parts=[3,5], Total=8
// "4-4"      -> Parts=[4,4], Total=8
// "3,4-5,2"  -> Parts=[3,4,5,2], Total=14
type Enum struct {
	Raw   string
	Parts []int
	Total int
}

// ParseEnum parses common crossword enumerations:
// - digits separated by ',' and/or '-' (commas and hyphens both mean "split into parts")
// - ignores whitespace
//
// It does NOT attempt to encode word boundaries semantics (comma vs hyphen).
// For validation we only need the total length and parts.
func ParseEnum(s string) (Enum, error) {
	raw := strings.TrimSpace(s)
	if raw == "" {
		return Enum{}, fmt.Errorf("enum is empty")
	}

	// normalize: remove spaces
	compact := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, raw)

	// Strict delimiter validation:
	// - cannot start/end with delimiter
	// - cannot contain consecutive delimiters
	isDelim := func(r byte) bool { return r == ',' || r == '-' }

	if len(compact) == 0 {
		return Enum{}, fmt.Errorf("enum is empty")
	}
	if isDelim(compact[0]) || isDelim(compact[len(compact)-1]) {
		return Enum{}, fmt.Errorf("enum %q cannot start or end with a delimiter", raw)
	}
	for i := 1; i < len(compact); i++ {
		if isDelim(compact[i]) && isDelim(compact[i-1]) {
			return Enum{}, fmt.Errorf("enum %q contains consecutive delimiters", raw)
		}
	}

	// split on comma or hyphen
	fields := strings.FieldsFunc(compact, func(r rune) bool {
		return r == ',' || r == '-'
	})

	if len(fields) == 0 {
		return Enum{}, fmt.Errorf("enum %q has no parts", raw)
	}

	parts := make([]int, 0, len(fields))
	total := 0

	for _, f := range fields {
		if f == "" {
			return Enum{}, fmt.Errorf("enum %q contains empty part", raw)
		}
		n, err := strconv.Atoi(f)
		if err != nil || n <= 0 {
			return Enum{}, fmt.Errorf("enum %q has invalid part %q", raw, f)
		}
		parts = append(parts, n)
		total += n
	}

	return Enum{Raw: raw, Parts: parts, Total: total}, nil
}

// NormalizedAnswerLen returns the letter-count of an answer ignoring:
// spaces, hyphens, apostrophes, and underscores.
//
// This matches typical crossword expectations where enumeration counts letters.
func NormalizedAnswerLen(answer string) int {
	n := 0
	for _, r := range answer {
		switch {
		case unicode.IsLetter(r):
			n++
		case unicode.IsDigit(r):
			// some puzzles allow digits; treat as countable
			n++
		default:
			// ignore punctuation/spacing
		}
	}
	return n
}
