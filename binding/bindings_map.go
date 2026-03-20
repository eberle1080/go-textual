package binding

import "strings"

// BindingEntry is a (key, binding) pair returned by BindingsMap.Iter.
type BindingEntry struct {
	Key     string
	Binding Binding
}

// BindingsMap manages a map from key strings to slices of Binding values.
type BindingsMap struct {
	KeyToBindings map[string][]Binding
}

// NewBindingsMap constructs a BindingsMap from a slice of binding inputs.
// See MakeBindings for the accepted input types.
func NewBindingsMap(inputs ...any) (*BindingsMap, error) {
	m := &BindingsMap{KeyToBindings: make(map[string][]Binding)}
	bindings, err := MakeBindings(inputs)
	if err != nil {
		return nil, err
	}
	for _, b := range bindings {
		m.KeyToBindings[b.Key] = append(m.KeyToBindings[b.Key], b)
	}
	return m, nil
}

// FromKeys constructs a BindingsMap from a pre-built key→bindings map.
func FromKeys(keyMap map[string][]Binding) *BindingsMap {
	return &BindingsMap{KeyToBindings: keyMap}
}

// AddBinding adds a single binding to the map.
func (m *BindingsMap) AddBinding(b Binding) {
	m.KeyToBindings[b.Key] = append(m.KeyToBindings[b.Key], b)
}

// Iter returns all (key, binding) pairs in the map.
func (m *BindingsMap) Iter() []BindingEntry {
	var entries []BindingEntry
	for key, bindings := range m.KeyToBindings {
		for _, b := range bindings {
			entries = append(entries, BindingEntry{Key: key, Binding: b})
		}
	}
	return entries
}

// Copy returns a shallow copy of the BindingsMap.
func (m *BindingsMap) Copy() *BindingsMap {
	c := &BindingsMap{KeyToBindings: make(map[string][]Binding, len(m.KeyToBindings))}
	for k, v := range m.KeyToBindings {
		dst := make([]Binding, len(v))
		copy(dst, v)
		c.KeyToBindings[k] = dst
	}
	return c
}

// Merge combines multiple BindingsMaps into a new map. Keys from all maps are
// accumulated (duplicates are preserved per-key).
func Merge(maps ...*BindingsMap) *BindingsMap {
	result := &BindingsMap{KeyToBindings: make(map[string][]Binding)}
	for _, m := range maps {
		if m == nil {
			continue
		}
		for key, bindings := range m.KeyToBindings {
			result.KeyToBindings[key] = append(result.KeyToBindings[key], bindings...)
		}
	}
	return result
}

// Bind adds new bindings for a comma-separated list of keys.
func (m *BindingsMap) Bind(keyStr, action, description string, show bool, keyDisplay *string, priority bool) {
	for _, k := range strings.Split(keyStr, ",") {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		b := Binding{
			Key:         k,
			Action:      action,
			Description: description,
			Show:        show,
			KeyDisplay:  keyDisplay,
			Priority:    priority,
		}
		m.KeyToBindings[k] = append(m.KeyToBindings[k], b)
	}
}

// GetBindingsForKey returns the bindings registered for the given key.
// Returns a *NoBinding error if no bindings are found.
func (m *BindingsMap) GetBindingsForKey(key string) ([]Binding, error) {
	bindings, ok := m.KeyToBindings[key]
	if !ok || len(bindings) == 0 {
		return nil, &NoBinding{BindingError{Msg: "no binding for key: " + key}}
	}
	return bindings, nil
}

// ShownKeys returns all bindings where Show is true.
func (m *BindingsMap) ShownKeys() []Binding {
	var result []Binding
	for _, bindings := range m.KeyToBindings {
		for _, b := range bindings {
			if b.Show {
				result = append(result, b)
			}
		}
	}
	return result
}

// ApplyKeymap substitutes binding keys based on ID lookup in keymap.
// keymap maps binding ID strings to new key strings (comma-separated for multiple).
// Returns a KeymapApplyResult describing any clashed bindings.
func (m *BindingsMap) ApplyKeymap(keymap map[string]string) KeymapApplyResult {
	result := KeymapApplyResult{ClashedBindings: make(map[Binding]struct{})}
	if len(keymap) == 0 {
		return result
	}

	// Collect all bindings being remapped and which IDs are involved.
	type remapInfo struct {
		oldKey  string
		binding Binding
		newKeys []string
	}
	var remaps []remapInfo
	remappedIDs := make(map[string]bool)

	for key, bindings := range m.KeyToBindings {
		for _, b := range bindings {
			if b.ID == nil {
				continue
			}
			if newKeys, ok := keymap[*b.ID]; ok {
				remaps = append(remaps, remapInfo{
					oldKey:  key,
					binding: b,
					newKeys: splitKeys(newKeys),
				})
				remappedIDs[*b.ID] = true
			}
		}
	}

	// Remove old key entries for remapped bindings.
	for _, r := range remaps {
		existing := m.KeyToBindings[r.oldKey]
		filtered := existing[:0]
		for _, b := range existing {
			if b.ID != nil && remappedIDs[*b.ID] {
				continue
			}
			filtered = append(filtered, b)
		}
		if len(filtered) == 0 {
			delete(m.KeyToBindings, r.oldKey)
		} else {
			m.KeyToBindings[r.oldKey] = filtered
		}
	}

	// Accumulate remapped bindings in a separate overlay so that multiple IDs
	// remapped to the same target key can coexist (mirroring Textual's new_bindings
	// overlay). Clash detection and default-removal operate only on the original map.
	newBindings := make(map[string][]Binding)
	for _, r := range remaps {
		for _, newKey := range r.newKeys {
			if existing, ok := m.KeyToBindings[newKey]; ok {
				// Record displaced non-remapped bindings as clashes and remove them.
				for _, eb := range existing {
					if eb.ID == nil || !remappedIDs[*eb.ID] {
						result.ClashedBindings[eb] = struct{}{}
					}
				}
				delete(m.KeyToBindings, newKey)
			}
			// Clear KeyDisplay so footer rendering reflects the overridden key string.
			newBinding := r.binding.WithKey(newKey, nil)
			newBindings[newKey] = append(newBindings[newKey], newBinding)
		}
	}

	// Merge overlay into the live map.
	for key, bindings := range newBindings {
		m.KeyToBindings[key] = append(m.KeyToBindings[key], bindings...)
	}

	return result
}

// splitKeys splits a comma-separated key string and trims whitespace.
func splitKeys(s string) []string {
	var out []string
	for _, k := range strings.Split(s, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			out = append(out, k)
		}
	}
	return out
}
