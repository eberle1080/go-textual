package cache

import "testing"

func TestFIFOCache_BasicHitMiss(t *testing.T) {
	c := NewFIFO[string, int](3)

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss on empty cache")
	}

	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected hit a=1, got %d ok=%v", v, ok)
	}
}

func TestFIFOCache_Eviction(t *testing.T) {
	c := NewFIFO[string, int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4) // evicts "a"

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected a to be evicted")
	}
	for _, key := range []string{"b", "c", "d"} {
		if _, ok := c.Get(key); !ok {
			t.Fatalf("expected %s to be present", key)
		}
	}
}

func TestFIFOCache_UpdateExisting(t *testing.T) {
	c := NewFIFO[string, int](3)
	c.Set("a", 1)
	c.Set("a", 99) // update without adding new key
	if v, ok := c.Get("a"); !ok || v != 99 {
		t.Fatalf("expected updated value 99, got %d", v)
	}
	if c.Len() != 1 {
		t.Fatalf("expected Len=1, got %d", c.Len())
	}
}

func TestFIFOCache_ZeroCapacity(t *testing.T) {
	c := NewFIFO[string, int](0)
	c.Set("a", 1)
	if _, ok := c.Get("a"); ok {
		t.Fatal("zero-capacity cache should never store")
	}
}

func TestFIFOCache_Clear(t *testing.T) {
	c := NewFIFO[string, int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Clear()
	if c.Len() != 0 {
		t.Fatal("expected Len=0 after Clear")
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss after Clear")
	}
}

func TestLRUCache_BasicHitMiss(t *testing.T) {
	c := NewLRU[string, int](3)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss on empty cache")
	}
	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected hit a=1, got %d ok=%v", v, ok)
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Get("a")    // make "a" most recently used
	c.Set("d", 4) // should evict "b" (LRU)

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected b to be evicted as LRU")
	}
	for _, key := range []string{"a", "c", "d"} {
		if _, ok := c.Get(key); !ok {
			t.Fatalf("expected %s to be present", key)
		}
	}
}

func TestLRUCache_UpdateExisting(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("a", 99) // update a, moves to front
	c.Set("d", 4)  // evicts LRU
	// "b" should be evicted since "a" was just updated (moved to front), "c" is next

	if v, ok := c.Get("a"); !ok || v != 99 {
		t.Fatalf("expected updated value a=99, got %d ok=%v", v, ok)
	}
}

func TestLRUCache_ZeroCapacity(t *testing.T) {
	c := NewLRU[string, int](0)
	c.Set("a", 1)
	if _, ok := c.Get("a"); ok {
		t.Fatal("zero-capacity cache should never store")
	}
}

func TestLRUCache_Clear(t *testing.T) {
	c := NewLRU[string, int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Clear()
	if c.Len() != 0 {
		t.Fatal("expected Len=0 after Clear")
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss after Clear")
	}
}
