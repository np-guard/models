// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube

import (
	"log"
	"slices"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
)

// CanonicalSet is a canonical representation for set of n-dimensional cubes
type CanonicalSet[S ds.Set[S]] struct {
	layers     *ds.IMap[S, *CanonicalSet[S]]
	dimensions int
}

// NewCanonicalSet returns a new empty CanonicalSet with n dimensions
func NewCanonicalSet[S ds.Set[S]](n int) *CanonicalSet[S] {
	return &CanonicalSet[S]{
		layers:     ds.NewIMap[S, *CanonicalSet[S]](),
		dimensions: n,
	}
}

// Equal return true if c equals other (same canonical form)
func (c *CanonicalSet[S]) Equal(other *CanonicalSet[S]) bool {
	if c == other {
		return true
	}
	if c.dimensions != other.dimensions {
		return false
	}
	return c.layers.Equal(other.layers)
}

const (
	hashX = 3
	hashY = 5
)

func (c *CanonicalSet[S]) Hash() int {
	if c.dimensions == 0 {
		return 1
	}
	res := hashX
	for _, p := range c.layers.Pairs() {
		res ^= hashY*p.Value.dimensions + (p.Key.Hash() << 1) ^ p.Value.Hash()
	}
	return res
}

// Union returns a new CanonicalSet object that results from union of c with other
func (c *CanonicalSet[S]) Union(other *CanonicalSet[S]) *CanonicalSet[S] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &CanonicalSet[S]{dimensions: other.dimensions, layers: c.layers.Union(other.layers)}
}

// IsEmpty returns true if c is empty
func (c *CanonicalSet[S]) IsEmpty() bool {
	return c.layers.IsEmpty()
}

// Intersect returns a new CanonicalSet object that results from intersection of c with other
func (c *CanonicalSet[S]) Intersect(other *CanonicalSet[S]) *CanonicalSet[S] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &CanonicalSet[S]{dimensions: other.dimensions, layers: c.layers.Intersect(other.layers)}
}

// Subtract returns a new CanonicalSet object that results from subtraction other from c
func (c *CanonicalSet[S]) Subtract(other *CanonicalSet[S]) *CanonicalSet[S] {
	if c == other {
		return NewCanonicalSet[S](c.dimensions)
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &CanonicalSet[S]{dimensions: other.dimensions, layers: c.layers.Subtract(other.layers)}
}

// ContainedIn returns true if c contained in other
func (c *CanonicalSet[S]) ContainedIn(other *CanonicalSet[S]) bool {
	if c.dimensions != other.dimensions {
		log.Panic("dimensionality mismatch")
	}
	return c.layers.ContainedIn(other.layers)
}

// Copy returns a new CanonicalSet object, copied from c
func (c *CanonicalSet[S]) Copy() *CanonicalSet[S] {
	res := NewCanonicalSet[S](c.dimensions)
	for _, p := range c.layers.Pairs() {
		res.layers.Insert(p.Key, p.Value.Copy())
	}
	return res
}

// GetCubesList returns the list of cubes in c, each cube as a slice of CanonicalSet
func (c *CanonicalSet[S]) Paths() [][]S {
	res := [][]S{}
	if c.dimensions == 1 {
		for _, k := range c.layers.Keys() {
			res = append(res, []S{k})
		}
		return res
	}
	for _, pair := range c.layers.Pairs() {
		subRes := pair.Value.Paths()
		for _, subList := range subRes {
			path := []S{pair.Key}
			path = append(path, subList...)
			res = append(res, path)
		}
	}
	return res
}

// SwapDimensions returns a new CanonicalSet object, built from the input CanonicalSet object,
// with dimensions dim1 and dim2 swapped
func (c *CanonicalSet[S]) SwapDimensions(dim1, dim2 int) *CanonicalSet[S] {
	if c.IsEmpty() || dim1 == dim2 {
		return c.Copy()
	}
	if min(dim1, dim2) < 0 || max(dim1, dim2) >= c.dimensions {
		log.Panicf("invalid dimensions: %d, %d", dim1, dim2)
	}
	res := NewCanonicalSet[S](c.dimensions)
	for _, path := range c.Paths() {
		if !path[dim1].Equal(path[dim2]) {
			// Shallow clone should be enough, since we do shallow swap:
			path = slices.Clone(path)
			path[dim1], path[dim2] = path[dim2], path[dim1]
		}
		res = res.Union(FromPath[S](path))
	}
	return res
}

// FromPath returns a new CanonicalSet created from a single input path
// the input cube is a slice of CanonicalSet, treated as ordered list of dimension values
func FromPath[S ds.Set[S]](path []S) *CanonicalSet[S] {
	if len(path) == 0 {
		return nil
	}
	if len(path) == 1 {
		res := NewCanonicalSet[S](1)
		res.layers.Insert(path[0], NewCanonicalSet[S](0))
		return res
	}
	res := NewCanonicalSet[S](len(path))
	res.layers.Insert(path[0], FromPath(path[1:]))
	return res
}

// Cube returns a new hypercube.CanonicalSet created from a single input cube
// the input cube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func Cube(values ...int64) *CanonicalSet[*interval.CanonicalSet] {
	path := []*interval.CanonicalSet{}
	for i := 0; i < len(values); i += 2 {
		path = append(path, interval.NewSetFromInterval(interval.New(values[i], values[i+1])))
	}
	return FromPath[*interval.CanonicalSet](path)
}
