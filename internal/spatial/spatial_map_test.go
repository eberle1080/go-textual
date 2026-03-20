package spatial

import (
	"testing"

	"github.com/eberle1080/go-textual/geometry"
)

func TestSpatialMap_BasicQuery(t *testing.T) {
	total := geometry.Region{X: 0, Y: 0, Width: 100, Height: 100}
	sm := New[int](total, 10, 10)

	sm.Insert([]Entry[int]{
		{Region: geometry.Region{X: 0, Y: 0, Width: 20, Height: 20}, Value: 1},
		{Region: geometry.Region{X: 50, Y: 50, Width: 20, Height: 20}, Value: 2},
	})

	// Query overlapping value 1 only.
	results := sm.GetValuesInRegion(geometry.Region{X: 0, Y: 0, Width: 10, Height: 10})
	if !containsInt(results, 1) {
		t.Fatalf("expected value 1 in results, got %v", results)
	}
	if containsInt(results, 2) {
		t.Fatalf("unexpected value 2 in results for non-overlapping query")
	}
}

func TestSpatialMap_FixedAlwaysReturned(t *testing.T) {
	total := geometry.Region{X: 0, Y: 0, Width: 100, Height: 100}
	sm := New[int](total, 10, 10)

	sm.Insert([]Entry[int]{
		{Fixed: true, Value: 99},
		{Region: geometry.Region{X: 80, Y: 80, Width: 10, Height: 10}, Value: 1},
	})

	// Query nowhere near value 1, but fixed should always appear.
	results := sm.GetValuesInRegion(geometry.Region{X: 0, Y: 0, Width: 5, Height: 5})
	if !containsInt(results, 99) {
		t.Fatal("expected fixed value 99 to always be returned")
	}
}

func TestSpatialMap_NoDuplicates(t *testing.T) {
	total := geometry.Region{X: 0, Y: 0, Width: 100, Height: 100}
	sm := New[int](total, 10, 10)

	// Large region spans many cells.
	sm.Insert([]Entry[int]{
		{Region: geometry.Region{X: 0, Y: 0, Width: 90, Height: 90}, Value: 42},
	})

	results := sm.GetValuesInRegion(geometry.Region{X: 0, Y: 0, Width: 90, Height: 90})
	count := 0
	for _, v := range results {
		if v == 42 {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly one occurrence of value 42, got %d", count)
	}
}

func containsInt(slice []int, v int) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}
