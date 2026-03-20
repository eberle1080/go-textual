package partition

import (
	"reflect"
	"testing"
)

func TestPartition_Integers(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	odds, evens := Partition(isEven, []int{1, 2, 3, 4, 5, 6})

	if !reflect.DeepEqual(odds, []int{1, 3, 5}) {
		t.Fatalf("expected odds [1 3 5], got %v", odds)
	}
	if !reflect.DeepEqual(evens, []int{2, 4, 6}) {
		t.Fatalf("expected evens [2 4 6], got %v", evens)
	}
}

func TestPartition_AllTrue(t *testing.T) {
	falseItems, trueItems := Partition(func(s string) bool { return true }, []string{"a", "b"})
	if len(falseItems) != 0 {
		t.Fatalf("expected no false items, got %v", falseItems)
	}
	if !reflect.DeepEqual(trueItems, []string{"a", "b"}) {
		t.Fatalf("expected [a b], got %v", trueItems)
	}
}

func TestPartition_AllFalse(t *testing.T) {
	falseItems, trueItems := Partition(func(s string) bool { return false }, []string{"a", "b"})
	if !reflect.DeepEqual(falseItems, []string{"a", "b"}) {
		t.Fatalf("expected [a b], got %v", falseItems)
	}
	if len(trueItems) != 0 {
		t.Fatalf("expected no true items, got %v", trueItems)
	}
}

func TestPartition_Empty(t *testing.T) {
	falseItems, trueItems := Partition(func(int) bool { return true }, nil)
	if falseItems != nil || trueItems != nil {
		t.Fatal("expected nil slices for empty input")
	}
}
