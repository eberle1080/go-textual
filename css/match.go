package css

// SelectorMatchNode extends SelectorNode with the path of nodes from root to
// this node (used for descendant/child combinator matching).
type SelectorMatchNode interface {
	SelectorNode
	// CSSPathNodes returns the path of nodes from the DOM root to this node,
	// with this node last.
	CSSPathNodes() []SelectorNode
}

// Match reports whether a node matches any of the given selector sets.
func Match(selectorSets []SelectorSet, node SelectorMatchNode) bool {
	pathNodes := node.CSSPathNodes()
	for _, ss := range selectorSets {
		if CheckSelectors(ss.Selectors, pathNodes) {
			return true
		}
	}
	return false
}

// CheckSelectors reports whether a list of selectors matches the given DOM path.
// pathNodes is the list of ancestor nodes from the root down to (and including)
// the target node; the target node is pathNodes[len(pathNodes)-1].
func CheckSelectors(selectors []Selector, pathNodes []SelectorNode) bool {
	if len(selectors) == 0 || len(pathNodes) == 0 {
		return false
	}

	node := pathNodes[len(pathNodes)-1]
	pathCount := len(pathNodes)
	selectorCount := len(selectors)

	// Stack entries: (selectorIndex, nodeIndex)
	type frame struct{ si, ni int }
	stack := []frame{{0, 0}}

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		si, ni := top.si, top.ni

		if si == selectorCount || ni == pathCount {
			stack = stack[:len(stack)-1]
			continue
		}

		pathNode := pathNodes[ni]
		sel := selectors[si]

		if sel.Combinator == CombinatorDescendant {
			if sel.Check(pathNode) {
				// Is this a complete match?
				if pathNode == node && si == selectorCount-1 {
					return true
				}
				// Advance selector and possibly node
				top.si = si + 1
				top.ni = ni + sel.Advance
				// Also push a frame to continue scanning without consuming this selector
				stack = append(stack, frame{si, ni + 1})
			} else {
				top.ni = ni + 1
			}
		} else {
			// Child combinator: must match exactly
			if sel.Check(pathNode) {
				if pathNode == node && si == selectorCount-1 {
					return true
				}
				top.si = si + 1
				top.ni = ni + sel.Advance
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}
	return false
}
