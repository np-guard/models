// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

type SetOps[S any] interface {
	IsEmpty() bool
	ContainedIn(S) bool
	Intersect(S) S
	Union(S) S
	Subtract(S) S
	fmt.Stringer
}

type HashableSet[S any] interface {
	Hashable[S]
	SetOps[S]
}

// CanonicalSet is a canonical representation for set of n-dimensional cubes
type CanonicalSet[T HashableSet[T]] struct {
	layers     Map[T, *CanonicalSet[T]]
	dimensions int
}

// NewCanonicalSet returns a new empty CanonicalSet with n dimensions
func NewCanonicalSet[T HashableSet[T]](n int) *CanonicalSet[T] {
	return &CanonicalSet[T]{
		layers:     NewMap[T, *CanonicalSet[T]](),
		dimensions: n,
	}
}

// Equal return true if c equals other (same canonical form)
func (c *CanonicalSet[T]) Equal(other *CanonicalSet[T]) bool {
	if c == other {
		return true
	}
	if c.dimensions != other.dimensions {
		return false
	}
	return c.layers.Equal(&other.layers)
}

const (
	hashX = 3
	hashY = 5
)

func (c *CanonicalSet[T]) Hash() int {
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
func (c *CanonicalSet[T]) Union(other *CanonicalSet[T]) *CanonicalSet[T] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	res := NewCanonicalSet[T](c.dimensions)
	remainingFromOther := NewMap[T, T]()
	for _, k := range other.layers.Keys() {
		remainingFromOther.Insert(k, k.Copy())
	}
	for _, pair := range c.layers.Pairs() {
		remainingFromSelf := pair.Key.Copy()
		for _, otherPair := range other.layers.Pairs() {
			commonElem := pair.Key.Intersect(otherPair.Key)
			if commonElem.IsEmpty() {
				continue
			}
			if v, ok := remainingFromOther.At(otherPair.Key); ok {
				remainingFromOther.Insert(otherPair.Key, v.Subtract(commonElem))
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			if c.dimensions == 1 {
				res.layers.Insert(commonElem, NewCanonicalSet[T](0))
				continue
			}
			newSubElem := pair.Value.Union(otherPair.Value)
			res.layers.Insert(commonElem, newSubElem)
		}
		if !remainingFromSelf.IsEmpty() {
			res.layers.Insert(remainingFromSelf, pair.Value.Copy())
		}
	}
	for _, pair := range remainingFromOther.Pairs() {
		if !pair.Value.IsEmpty() {
			if otherValue, ok := other.layers.At(pair.Key); ok {
				res.layers.Insert(pair.Value, otherValue.Copy())
			}
		}
	}
	res.canonicalize()
	return res
}

// IsEmpty returns true if c is empty
func (c *CanonicalSet[T]) IsEmpty() bool {
	return c.layers.IsEmpty()
}

// Intersect returns a new CanonicalSet object that results from intersection of c with other
func (c *CanonicalSet[T]) Intersect(other *CanonicalSet[T]) *CanonicalSet[T] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	res := NewCanonicalSet[T](c.dimensions)
	for _, pair := range c.layers.Pairs() {
		for _, otherPair := range other.layers.Pairs() {
			commonELem := pair.Key.Intersect(otherPair.Key)
			if commonELem.IsEmpty() {
				continue
			}
			if c.dimensions == 1 {
				res.layers.Insert(commonELem, NewCanonicalSet[T](0))
				continue
			}
			newSubElem := pair.Value.Intersect(otherPair.Value)
			if !newSubElem.IsEmpty() {
				res.layers.Insert(commonELem, newSubElem)
			}
		}
	}
	res.canonicalize()
	return res
}

// Subtract returns a new CanonicalSet object that results from subtraction other from c
func (c *CanonicalSet[T]) Subtract(other *CanonicalSet[T]) *CanonicalSet[T] {
	if c == other {
		return NewCanonicalSet[T](c.dimensions)
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	res := NewCanonicalSet[T](c.dimensions)
	for _, pair := range c.layers.Pairs() {
		remainingFromSelf := pair.Key.Copy()
		for _, otherPair := range other.layers.Pairs() {
			commonELem := pair.Key.Intersect(otherPair.Key)
			if commonELem.IsEmpty() {
				continue
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonELem)
			if c.dimensions == 1 {
				continue
			}
			newSubElem := pair.Value.Subtract(otherPair.Value)
			if !newSubElem.IsEmpty() {
				res.layers.Insert(commonELem, newSubElem)
			}
		}
		if !remainingFromSelf.IsEmpty() {
			res.layers.Insert(remainingFromSelf, pair.Value.Copy())
		}
	}
	res.canonicalize()
	return res
}

// ContainedIn returns true if c contained in other
func (c *CanonicalSet[T]) ContainedIn(other *CanonicalSet[T]) (bool, error) {
	if c == other {
		return true, nil
	}
	if c.dimensions != other.dimensions {
		return false, errors.New("ContainedIn mismatch between num of dimensions for input args")
	}
	if c.dimensions == 0 {
		return true, nil
	}

	subsetCount := 0
	for _, pair := range c.layers.Pairs() {
		LeftoverKey := pair.Key.Copy()
		for _, otherPair := range other.layers.Pairs() {
			commonKey := otherPair.Key.Intersect(LeftoverKey)
			if commonKey.IsEmpty() {
				continue
			}
			subContainment, err := pair.Value.ContainedIn(otherPair.Value)
			if err != nil {
				return false, err
			}
			if !subContainment {
				return false, nil
			}
			LeftoverKey = LeftoverKey.Subtract(commonKey)
			if LeftoverKey.IsEmpty() {
				subsetCount += 1
				break
			}
		}
	}
	return subsetCount == c.layers.Size(), nil
}

// Copy returns a new CanonicalSet object, copied from c
func (c *CanonicalSet[T]) Copy() *CanonicalSet[T] {
	res := NewCanonicalSet[T](c.dimensions)
	for _, p := range c.layers.Pairs() {
		res.layers.Insert(p.Key, p.Value.Copy())
	}
	return res
}

func getCubeStr[T HashableSet[T]](cube []T) string {
	strList := []string{}
	for _, v := range cube {
		strList = append(strList, "("+v.String()+")")
	}
	return "[" + strings.Join(strList, ",") + "]"
}

// String returns a string representation of c
func (c *CanonicalSet[T]) String() string {
	cubesList := c.GetCubesList()
	strList := []string{}
	for _, cube := range cubesList {
		strList = append(strList, getCubeStr(cube))
	}
	sort.Strings(strList)
	return strings.Join(strList, "; ")
}

// GetCubesList returns the list of cubes in c, each cube as a slice of CanonicalSet
func (c *CanonicalSet[T]) GetCubesList() [][]T {
	res := [][]T{}
	if c.dimensions == 1 {
		for _, k := range c.layers.Keys() {
			res = append(res, []T{k})
		}
		return res
	}
	for _, pair := range c.layers.Pairs() {
		subRes := pair.Value.GetCubesList()
		for _, subList := range subRes {
			cube := []T{pair.Key}
			cube = append(cube, subList...)
			res = append(res, cube)
		}
	}
	return res
}

// SwapDimensions returns a new CanonicalSet object, built from the input CanonicalSet object,
// with dimensions dim1 and dim2 swapped
func (c *CanonicalSet[T]) SwapDimensions(dim1, dim2 int) *CanonicalSet[T] {
	if c.IsEmpty() || dim1 == dim2 {
		return c.Copy()
	}
	if min(dim1, dim2) < 0 || max(dim1, dim2) >= c.dimensions {
		log.Panicf("invalid dimensions: %d, %d", dim1, dim2)
	}
	res := NewCanonicalSet[T](c.dimensions)
	for _, cube := range c.GetCubesList() {
		if !cube[dim1].Equal(cube[dim2]) {
			// Shallow clone should be enough, since we do shallow swap:
			cube = slices.Clone(cube)
			cube[dim1], cube[dim2] = cube[dim2], cube[dim1]
		}
		res = res.Union(FromCube(cube))
	}
	return res
}

func (c *CanonicalSet[T]) canonicalize() {
	newLayers := NewMap[T, *CanonicalSet[T]]()
	for _, p := range InverseMap(&c.layers).Pairs() {
		items := p.Value.Items()
		if len(items) == 0 {
			continue
		}
		newKey := items[0]
		for _, v := range items[1:] {
			newKey = newKey.Union(v)
		}
		newLayers.Insert(newKey, p.Key)
	}
	c.layers = newLayers
}

// FromCube returns a new CanonicalSet created from a single input cube
// the input cube is a slice of CanonicalSet, treated as ordered list of dimension values
func FromCube[T HashableSet[T]](cube []T) *CanonicalSet[T] {
	if len(cube) == 0 {
		return nil
	}
	if len(cube) == 1 {
		res := NewCanonicalSet[T](1)
		res.layers.Insert(cube[0], NewCanonicalSet[T](0))
		return res
	}
	res := NewCanonicalSet[T](len(cube))
	res.layers.Insert(cube[0], FromCube[T](cube[1:]))
	return res
}

// Cube returns a new hypercube.CanonicalSet created from a single input cube
// the input cube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func Cube(values ...int64) *CanonicalSet[*interval.CanonicalSet] {
	cube := []*interval.CanonicalSet{}
	for i := 0; i < len(values); i += 2 {
		cube = append(cube, interval.NewSetFromInterval(interval.New(values[i], values[i+1])))
	}
	return FromCube(cube)
}
