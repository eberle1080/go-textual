package layoutresolve

import (
	"reflect"
	"testing"
)

type testEdge struct {
	size     *int
	fraction int
	minSize  int
}

func (e testEdge) Size() *int    { return e.size }
func (e testEdge) Fraction() int { return e.fraction }
func (e testEdge) MinSize() int  { return e.minSize }

func intPtr(n int) *int { return &n }

func TestResolve_ExplicitSizes(t *testing.T) {
	edges := []Edge{
		testEdge{size: intPtr(10)},
		testEdge{size: intPtr(20)},
		testEdge{size: intPtr(30)},
	}
	result := Resolve(100, edges)
	if !reflect.DeepEqual(result, []int{10, 20, 30}) {
		t.Fatalf("expected [10 20 30], got %v", result)
	}
}

func TestResolve_FractionalEqual(t *testing.T) {
	edges := []Edge{
		testEdge{fraction: 1},
		testEdge{fraction: 1},
		testEdge{fraction: 1},
	}
	result := Resolve(90, edges)
	sum := 0
	for _, v := range result {
		sum += v
	}
	if sum != 90 {
		t.Fatalf("expected sum 90, got %d (result=%v)", sum, result)
	}
}

func TestResolve_Mixed(t *testing.T) {
	edges := []Edge{
		testEdge{size: intPtr(20)}, // fixed
		testEdge{fraction: 1},      // gets remaining / 2
		testEdge{fraction: 1},      // gets remaining / 2
	}
	result := Resolve(100, edges)
	if result[0] != 20 {
		t.Fatalf("expected result[0]=20, got %d", result[0])
	}
	if result[1]+result[2] != 80 {
		t.Fatalf("expected fractional sum=80, got %d+%d=%d", result[1], result[2], result[1]+result[2])
	}
}

func TestResolve_MinSize(t *testing.T) {
	edges := []Edge{
		testEdge{fraction: 1, minSize: 5},
		testEdge{fraction: 1, minSize: 5},
	}
	// total is tiny; min-size should be respected
	result := Resolve(4, edges)
	for i, v := range result {
		if v < 5 {
			t.Fatalf("result[%d]=%d violates minSize=5", i, v)
		}
	}
}

func TestResolve_Empty(t *testing.T) {
	result := Resolve(100, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}
