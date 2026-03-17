package activities

import (
	"testing"
)

func TestAllActivities(t *testing.T) {
	all := All()
	if len(all) < 20 {
		t.Errorf("expected at least 20 activities, got %d", len(all))
	}
}

func TestRandom(t *testing.T) {
	a := Random()
	if a.Name == "" {
		t.Error("expected non-empty activity name")
	}
	if a.Category == "" {
		t.Error("expected non-empty category")
	}
	if a.DurationSec <= 0 {
		t.Error("expected positive duration")
	}
}

func TestRandomFromCategory(t *testing.T) {
	categories := []Category{Stretch, Eyes, Movement, Breathing}
	for _, cat := range categories {
		a := RandomFromCategory(cat)
		if a.Category != cat {
			t.Errorf("expected category %s, got %s", cat, a.Category)
		}
	}
}

func TestCategories(t *testing.T) {
	counts := map[Category]int{}
	for _, a := range All() {
		counts[a.Category]++
	}

	for _, cat := range []Category{Stretch, Eyes, Movement, Breathing} {
		if counts[cat] == 0 {
			t.Errorf("no activities in category %s", cat)
		}
	}
}
