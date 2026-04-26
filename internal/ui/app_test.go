package ui

import "testing"

func TestDefaultSpecUsesResearchTitleAndModules(t *testing.T) {
	spec := DefaultSpec()

	if spec.Title != "Moscovium Statera Go" {
		t.Fatalf("Title = %q", spec.Title)
	}
	if spec.Width < 900 || spec.Height < 600 {
		t.Fatalf("unexpected window size: %dx%d", spec.Width, spec.Height)
	}

	wantModules := []string{"Education", "Research", "Design"}
	if len(spec.Modules) != len(wantModules) {
		t.Fatalf("module count = %d, want %d", len(spec.Modules), len(wantModules))
	}
	for i, want := range wantModules {
		if spec.Modules[i].Name != want {
			t.Fatalf("module[%d] name = %q, want %q", i, spec.Modules[i].Name, want)
		}
		if spec.Modules[i].Description == "" {
			t.Fatalf("module[%d] missing description", i)
		}
	}
}
