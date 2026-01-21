package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWordlist_AnagramsAndPattern(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(fp, []byte("react\ntrace\ncrate\ncater\ncat\n"), 0644); err != nil {
		t.Fatalf("write temp wordlist: %v", err)
	}

	wl, err := LoadWordlist(fp)
	if err != nil {
		t.Fatalf("LoadWordlist: %v", err)
	}

	ana, err := wl.Anagrams("REACT", 5)
	if err != nil {
		t.Fatalf("Anagrams: %v", err)
	}
	if len(ana) == 0 {
		t.Fatalf("expected anagrams, got none")
	}

	pm, err := wl.PatternMatch("re?c?", 5)
	if err != nil {
		t.Fatalf("PatternMatch: %v", err)
	}
	if len(pm) == 0 {
		t.Fatalf("expected pattern matches for %q, got none", "tr?c?")
	}
}
