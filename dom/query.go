package dom

import "github.com/eberle1080/go-textual/css"

// DOMQuery is the result of a CSS selector query on a DOM subtree.
type DOMQuery struct {
	root     Node
	selector string
	deep     bool
	filterFn func(Node) bool
	exclude  string
}

func newQuery(root Node, selector string, deep bool) *DOMQuery {
	return &DOMQuery{root: root, selector: selector, deep: deep}
}

// Results returns all nodes matching the selector.
func (q *DOMQuery) Results() ([]Node, error) {
	selectorSets, err := css.ParseSelectors(q.selector)
	if err != nil {
		return nil, &InvalidQueryFormatError{Selector: q.selector, Cause: err}
	}

	var excludeSets []css.SelectorSet
	if q.exclude != "" {
		excludeSets, err = css.ParseSelectors(q.exclude)
		if err != nil {
			return nil, &InvalidQueryFormatError{Selector: q.exclude, Cause: err}
		}
	}

	var results []Node
	var walk func(Node)
	walk = func(node Node) {
		for _, child := range node.Children().Slice() {
			if css.Match(selectorSets, child) {
				if len(excludeSets) > 0 && css.Match(excludeSets, child) {
					// excluded
				} else if q.filterFn != nil && !q.filterFn(child) {
					// filtered
				} else {
					results = append(results, child)
				}
			}
			if q.deep {
				walk(child)
			}
		}
	}
	walk(q.root)
	return results, nil
}

// First returns the first matching node.
func (q *DOMQuery) First() (Node, error) {
	results, err := q.Results()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &NoMatchesError{Selector: q.selector}
	}
	return results[0], nil
}

// Last returns the last matching node.
func (q *DOMQuery) Last() (Node, error) {
	results, err := q.Results()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &NoMatchesError{Selector: q.selector}
	}
	return results[len(results)-1], nil
}

// Filter returns a new DOMQuery with an additional filter.
func (q *DOMQuery) Filter(fn func(Node) bool) *DOMQuery {
	cp := *q
	if q.filterFn != nil {
		prev := q.filterFn
		cp.filterFn = func(n Node) bool { return prev(n) && fn(n) }
	} else {
		cp.filterFn = fn
	}
	return &cp
}

// Exclude returns a new DOMQuery that excludes nodes matching excludeSelector.
func (q *DOMQuery) Exclude(excludeSelector string) *DOMQuery {
	cp := *q
	cp.exclude = excludeSelector
	return &cp
}
