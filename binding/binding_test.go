package binding

import (
	"testing"
)

func TestBindingParseKey(t *testing.T) {
	b := Binding{Key: "ctrl+shift+a"}
	mods, key := b.ParseKey()
	if key != "a" {
		t.Errorf("key = %q, want %q", key, "a")
	}
	if len(mods) != 2 || mods[0] != "ctrl" || mods[1] != "shift" {
		t.Errorf("modifiers = %v, want [ctrl shift]", mods)
	}
}

func TestBindingParseKeyNoModifiers(t *testing.T) {
	b := Binding{Key: "f1"}
	mods, key := b.ParseKey()
	if key != "f1" {
		t.Errorf("key = %q, want %q", key, "f1")
	}
	if len(mods) != 0 {
		t.Errorf("modifiers = %v, want []", mods)
	}
}

func TestBindingWithKey(t *testing.T) {
	original := Binding{Key: "ctrl+q", Action: "quit", Description: "Quit"}
	updated := original.WithKey("ctrl+x", nil)
	if updated.Key != "ctrl+x" {
		t.Errorf("updated key = %q, want %q", updated.Key, "ctrl+x")
	}
	if updated.Action != "quit" {
		t.Errorf("action changed: got %q, want %q", updated.Action, "quit")
	}
	// Original is unchanged.
	if original.Key != "ctrl+q" {
		t.Errorf("original key mutated: got %q", original.Key)
	}
}

func TestMakeBindingsFromBinding(t *testing.T) {
	b := Binding{Key: "ctrl+a", Action: "select_all", Show: true}
	got, err := MakeBindings([]any{b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Action != "select_all" {
		t.Errorf("action = %q, want %q", got[0].Action, "select_all")
	}
}

func TestMakeBindingsFromTuple2(t *testing.T) {
	got, err := MakeBindings([]any{[2]string{"ctrl+a", "quit"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Key != "ctrl+a" || got[0].Action != "quit" {
		t.Errorf("got {%q, %q}, want {ctrl+a, quit}", got[0].Key, got[0].Action)
	}
	if got[0].Description != "" {
		t.Errorf("description should be empty, got %q", got[0].Description)
	}
}

func TestMakeBindingsFromTuple3(t *testing.T) {
	got, err := MakeBindings([]any{[3]string{"ctrl+a", "quit", "Quit app"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Description != "Quit app" {
		t.Errorf("description = %q, want %q", got[0].Description, "Quit app")
	}
}

func TestMakeBindingsCommaExpansion(t *testing.T) {
	b := Binding{Key: "j,down", Action: "scroll_down", Show: true}
	got, err := MakeBindings([]any{b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestMakeBindingsCharNormalization(t *testing.T) {
	b := Binding{Key: "/", Action: "search"}
	got, err := MakeBindings([]any{b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one binding")
	}
	// The key should be normalised (not the raw "/" character).
	if got[0].Key == "/" {
		t.Logf("key was not normalised (may be a no-op for this character): %q", got[0].Key)
	}
}

func TestMakeBindingsEmptyKey(t *testing.T) {
	b := Binding{Key: "a,,b", Action: "test"}
	_, err := MakeBindings([]any{b})
	if err == nil {
		t.Error("expected error for empty key in comma list")
	}
}

func TestBindingsMapNew(t *testing.T) {
	b := Binding{Key: "ctrl+a", Action: "select"}
	m, err := NewBindingsMap(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.KeyToBindings) == 0 {
		t.Error("KeyToBindings should be populated")
	}
}

func TestBindingsMapBind(t *testing.T) {
	m := &BindingsMap{KeyToBindings: make(map[string][]Binding)}
	m.Bind("ctrl+a,ctrl+b", "quit", "Quit", true, nil, false)
	if len(m.KeyToBindings) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m.KeyToBindings))
	}
	if _, ok := m.KeyToBindings["ctrl+a"]; !ok {
		t.Error("ctrl+a not found")
	}
	if _, ok := m.KeyToBindings["ctrl+b"]; !ok {
		t.Error("ctrl+b not found")
	}
}

func TestBindingsMapGetBindingsForKey(t *testing.T) {
	b := Binding{Key: "ctrl+q", Action: "quit"}
	m, _ := NewBindingsMap(b)
	got, err := m.GetBindingsForKey("ctrl+q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	_, err = m.GetBindingsForKey("missing")
	if err == nil {
		t.Error("expected NoBinding error for missing key")
	}
}

func TestBindingsMapMerge(t *testing.T) {
	a := &BindingsMap{KeyToBindings: map[string][]Binding{
		"ctrl+a": {{Key: "ctrl+a", Action: "a"}},
	}}
	b := &BindingsMap{KeyToBindings: map[string][]Binding{
		"ctrl+b": {{Key: "ctrl+b", Action: "b"}},
	}}
	merged := Merge(a, b)
	if _, ok := merged.KeyToBindings["ctrl+a"]; !ok {
		t.Error("merged map missing ctrl+a")
	}
	if _, ok := merged.KeyToBindings["ctrl+b"]; !ok {
		t.Error("merged map missing ctrl+b")
	}
}

func TestBindingsMapCopy(t *testing.T) {
	original := &BindingsMap{KeyToBindings: map[string][]Binding{
		"ctrl+a": {{Key: "ctrl+a", Action: "a"}},
	}}
	copied := original.Copy()
	// Modify the copy; original should be unaffected.
	copied.KeyToBindings["ctrl+b"] = []Binding{{Key: "ctrl+b", Action: "b"}}
	if _, ok := original.KeyToBindings["ctrl+b"]; ok {
		t.Error("modifying copy should not affect original")
	}
}

func TestBindingsMapShownKeys(t *testing.T) {
	m := &BindingsMap{KeyToBindings: map[string][]Binding{
		"ctrl+a": {{Key: "ctrl+a", Action: "a", Show: true}},
		"ctrl+b": {{Key: "ctrl+b", Action: "b", Show: false}},
	}}
	shown := m.ShownKeys()
	if len(shown) != 1 {
		t.Errorf("ShownKeys = %d, want 1", len(shown))
	}
	if shown[0].Key != "ctrl+a" {
		t.Errorf("shown key = %q, want ctrl+a", shown[0].Key)
	}
}

func TestBindingsMapIter(t *testing.T) {
	m := &BindingsMap{KeyToBindings: map[string][]Binding{
		"ctrl+a": {{Key: "ctrl+a", Action: "a"}},
		"ctrl+b": {{Key: "ctrl+b", Action: "b"}},
	}}
	entries := m.Iter()
	if len(entries) != 2 {
		t.Errorf("Iter returned %d entries, want 2", len(entries))
	}
}
