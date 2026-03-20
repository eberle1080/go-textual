// Package spatial provides a grid-based spatial map for efficient region
// queries. It is used by the layout engine to find visible widget placements
// within a viewport region.
package spatial

import (
	"github.com/eberle1080/go-textual/geometry"
)

// SpatialMap stores values associated with rectangular regions and supports
// efficient retrieval by overlapping query region.
//
// Values are indexed into a grid of cells; a single value can span multiple
// cells. "Fixed" values (added without a region) are always returned.
type SpatialMap[T any] struct {
	gridWidth  int
	gridHeight int
	total      geometry.Region
	grid       map[[2]int][]T
	fixed      []T
}

// Entry associates a value with the region it occupies.
type Entry[T any] struct {
	Region geometry.Region
	Value  T
	Fixed  bool // if true, Value is returned for every query
}

// New creates a SpatialMap that covers totalRegion and divides it into a
// grid of cells each approximately cellWidth×cellHeight in size.
func New[T any](totalRegion geometry.Region, cellWidth, cellHeight int) *SpatialMap[T] {
	if cellWidth < 1 {
		cellWidth = 1
	}
	if cellHeight < 1 {
		cellHeight = 1
	}
	gw := (totalRegion.Width + cellWidth - 1) / cellWidth
	gh := (totalRegion.Height + cellHeight - 1) / cellHeight
	if gw < 1 {
		gw = 1
	}
	if gh < 1 {
		gh = 1
	}
	return &SpatialMap[T]{
		gridWidth:  cellWidth,
		gridHeight: cellHeight,
		total:      totalRegion,
		grid:       make(map[[2]int][]T),
	}
}

// Insert adds a slice of entries into the map.
func (s *SpatialMap[T]) Insert(entries []Entry[T]) {
	for _, e := range entries {
		if e.Fixed {
			s.fixed = append(s.fixed, e.Value)
			continue
		}
		r := e.Region
		colStart := r.X / s.gridWidth
		colEnd := (r.X + r.Width - 1) / s.gridWidth
		rowStart := r.Y / s.gridHeight
		rowEnd := (r.Y + r.Height - 1) / s.gridHeight
		for col := colStart; col <= colEnd; col++ {
			for row := rowStart; row <= rowEnd; row++ {
				key := [2]int{col, row}
				s.grid[key] = append(s.grid[key], e.Value)
			}
		}
	}
}

// GetValuesInRegion returns all values whose region overlaps with query.
// Fixed values are always included. Duplicate entries are deduplicated by
// pointer identity using a simple seen-map.
func (s *SpatialMap[T]) GetValuesInRegion(query geometry.Region) []T {
	colStart := query.X / s.gridWidth
	colEnd := (query.X + query.Width - 1) / s.gridWidth
	rowStart := query.Y / s.gridHeight
	rowEnd := (query.Y + query.Height - 1) / s.gridHeight

	// Collect, dedup via a flat slice (small grids are common).
	var result []T
	result = append(result, s.fixed...)

	seen := make(map[any]bool)
	for _, v := range s.fixed {
		seen[any(v)] = true
	}

	for col := colStart; col <= colEnd; col++ {
		for row := rowStart; row <= rowEnd; row++ {
			for _, v := range s.grid[[2]int{col, row}] {
				key := any(v)
				if !seen[key] {
					seen[key] = true
					result = append(result, v)
				}
			}
		}
	}
	return result
}

// TotalRegion returns the bounding region this spatial map was created for.
func (s *SpatialMap[T]) TotalRegion() geometry.Region { return s.total }
