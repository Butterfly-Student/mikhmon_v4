package util

import (
	"strings"
	"testing"
)

func TestRandN(t *testing.T) {
	s := RandN(6)
	if len(s) != 6 {
		t.Errorf("RandN(6) length = %d, want 6", len(s))
	}
	for _, ch := range s {
		if !strings.ContainsRune(digits, ch) {
			t.Errorf("RandN() contains non-digit: %c", ch)
		}
	}
}

func TestRandLC(t *testing.T) {
	s := RandLC(8)
	if len(s) != 8 {
		t.Errorf("RandLC(8) length = %d, want 8", len(s))
	}
	for _, ch := range s {
		if !strings.ContainsRune(lower, ch) {
			t.Errorf("RandLC() contains non-lowercase: %c", ch)
		}
	}
}

func TestRandUC(t *testing.T) {
	s := RandUC(5)
	if len(s) != 5 {
		t.Errorf("RandUC(5) length = %d, want 5", len(s))
	}
	for _, ch := range s {
		if !strings.ContainsRune(upper, ch) {
			t.Errorf("RandUC() contains non-uppercase: %c", ch)
		}
	}
}

func TestRandNLC(t *testing.T) {
	s := RandNLC(10)
	if len(s) != 10 {
		t.Errorf("RandNLC(10) length = %d, want 10", len(s))
	}
	for _, ch := range s {
		if !strings.ContainsRune(digits+lower, ch) {
			t.Errorf("RandNLC() contains invalid char: %c", ch)
		}
	}
}

func TestRandNULC(t *testing.T) {
	for i := 0; i < 20; i++ {
		s := RandNULC(6)
		if len(s) != 6 {
			t.Errorf("RandNULC(6) length = %d, want 6", len(s))
		}
	}
}

func TestUniqueness(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		s := RandNULC(8)
		seen[s] = true
	}
	// Should have many unique values (not all identical)
	if len(seen) < 50 {
		t.Errorf("expected high uniqueness, got %d unique values out of 100", len(seen))
	}
}
