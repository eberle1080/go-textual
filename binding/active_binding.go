package binding

// BindingNode is a forward-reference interface for the DOM node that owns a
// binding. It avoids importing the future dom package.
type BindingNode interface {
	NodeID() string
}

// ActiveBinding pairs a Binding with the DOM node where it is defined and
// whether it is currently enabled.
type ActiveBinding struct {
	Node    BindingNode // owning DOM node (provides namespace for grouping/metadata)
	Binding Binding     // The binding configuration
	Enabled bool        // Whether the binding is currently enabled
	Tooltip string      // Optional tooltip
}

// NodeID returns the string ID of the owning node, or "" if Node is nil.
func (a ActiveBinding) NodeID() string {
	if a.Node == nil {
		return ""
	}
	return a.Node.NodeID()
}

// KeymapApplyResult is returned by BindingsMap.ApplyKeymap and reports bindings
// displaced by the keymap substitution.
type KeymapApplyResult struct {
	ClashedBindings map[Binding]struct{} // set of bindings displaced by remapping
}

// HasClashes reports whether any bindings were displaced by the keymap.
func (r KeymapApplyResult) HasClashes() bool {
	return len(r.ClashedBindings) > 0
}
