// Package partition provides a generic partition utility.
// It is a direct port of textual/_partition.py.
package partition

// Partition divides items into two slices based on predicate.
// Items for which predicate returns false go into the first return value;
// items for which it returns true go into the second.
//
// The relative order of items within each group is preserved.
func Partition[T any](predicate func(T) bool, items []T) (falseItems, trueItems []T) {
	for _, item := range items {
		if predicate(item) {
			trueItems = append(trueItems, item)
		} else {
			falseItems = append(falseItems, item)
		}
	}
	return falseItems, trueItems
}
