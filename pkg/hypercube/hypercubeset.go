/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package hypercube

import (
	"errors"
	"log"
	"slices"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

// CanonicalSet is a canonical representation for set of n-dimensional cubes, from integer intervals
type CanonicalSet struct {
	layers     map[*interval.CanonicalSet]*CanonicalSet
	dimensions int
}

// NewCanonicalSet returns a new empty CanonicalSet with n dimensions
func NewCanonicalSet(n int) *CanonicalSet {
	return &CanonicalSet{
		layers:     map[*interval.CanonicalSet]*CanonicalSet{},
		dimensions: n,
	}
}

// Equal return true if c equals other (same canonical form)
func (c *CanonicalSet) Equal(other *CanonicalSet) bool {
	if c == other {
		return true
	}
	if c.dimensions != other.dimensions {
		return false
	}
	if len(c.layers) != len(other.layers) {
		return false
	}
	if len(c.layers) == 0 {
		return true
	}
	mapByString := map[string]*CanonicalSet{}
	for k, v := range c.layers {
		mapByString[k.String()] = v
	}
	for k, v := range other.layers {
		if w, ok := mapByString[k.String()]; !ok || !v.Equal(w) {
			return false
		}
	}
	return true
}

// Union returns a new CanonicalSet object that results from union of c with other
func (c *CanonicalSet) Union(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	remainingFromOther := map[*interval.CanonicalSet]*interval.CanonicalSet{}
	for otherKey := range other.layers {
		remainingFromOther[otherKey] = otherKey.Copy()
	}
	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		remainingFromSelf := k.Copy()
		for otherKey, otherVal := range other.layers {
			commonElem := k.Intersect(otherKey)
			if commonElem.IsEmpty() {
				continue
			}
			remainingFromOther[otherKey] = remainingFromOther[otherKey].Subtract(commonElem)
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			newSubElem := NewCanonicalSet(0)
			if c.dimensions != 1 {
				newSubElem = v.Union(otherVal)
			}
			layers[commonElem] = newSubElem
		}
		if !remainingFromSelf.IsEmpty() {
			layers[remainingFromSelf] = v.Copy()
		}
	}
	for k, v := range remainingFromOther {
		if !v.IsEmpty() {
			layers[v] = other.layers[k].Copy()
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

// IsEmpty returns true if c is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.layers) == 0
}

// Intersect returns a new CanonicalSet object that results from intersection of c with other
func (c *CanonicalSet) Intersect(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}

	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		for otherKey, otherVal := range other.layers {
			commonELem := k.Intersect(otherKey)
			if commonELem.IsEmpty() {
				continue
			}
			if c.dimensions == 1 {
				layers[commonELem] = NewCanonicalSet(0)
				continue
			}
			newSubElem := v.Intersect(otherVal)
			if !newSubElem.IsEmpty() {
				layers[commonELem] = newSubElem
			}
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

// Subtract returns a new CanonicalSet object that results from subtraction other from c
func (c *CanonicalSet) Subtract(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return NewCanonicalSet(c.dimensions)
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	layers := map[*interval.CanonicalSet]*CanonicalSet{}
	for k, v := range c.layers {
		remainingFromSelf := k.Copy()
		for otherKey, otherVal := range other.layers {
			commonElem := k.Intersect(otherKey)
			if commonElem.IsEmpty() {
				continue
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			if c.dimensions == 1 {
				continue
			}
			newSubElem := v.Subtract(otherVal)
			if !newSubElem.IsEmpty() {
				layers[commonElem] = newSubElem
			}
		}
		if !remainingFromSelf.IsEmpty() {
			layers[remainingFromSelf] = v.Copy()
		}
	}
	return &CanonicalSet{
		layers:     getElementsUnionPerLayer(layers),
		dimensions: c.dimensions,
	}
}

// ContainedIn returns true if c is subset of other
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) (bool, error) {
	if c == other {
		return true, nil
	}
	if c.dimensions != other.dimensions {
		return false, errors.New("ContainedIn mismatch between num of dimensions for input args")
	}
	if c.dimensions == 0 {
		if len(c.layers) != 0 || len(other.layers) != 0 {
			return false, errors.New("unexpected non-empty object of dimension size 0")
		}
		return true, nil
	}

	isSubsetCount := 0
	for currentLayer, v := range c.layers {
		for otherKey, otherVal := range other.layers {
			commonKey := currentLayer.Intersect(otherKey)
			remaining := currentLayer.Subtract(commonKey)
			if !commonKey.IsEmpty() {
				subContainment, err := v.ContainedIn(otherVal)
				if !subContainment || err != nil {
					return subContainment, err
				}
				if !remaining.IsEmpty() {
					currentLayer = remaining
				} else {
					isSubsetCount += 1
					break
				}
			}
		}
	}
	return isSubsetCount == len(c.layers), nil
}

// Copy returns a new CanonicalSet object, copied from c
func (c *CanonicalSet) Copy() *CanonicalSet {
	res := NewCanonicalSet(c.dimensions)
	for k, v := range c.layers {
		res.layers[k.Copy()] = v.Copy()
	}
	return res
}

func getCubeStr(cube []*interval.CanonicalSet) string {
	strList := []string{}
	for _, v := range cube {
		strList = append(strList, "("+v.String()+")")
	}
	return "[" + strings.Join(strList, ",") + "]"
}

// String returns a string representation of c
func (c *CanonicalSet) String() string {
	cubesList := c.GetCubesList()
	strList := []string{}
	for _, cube := range cubesList {
		strList = append(strList, getCubeStr(cube))
	}
	sort.Strings(strList)
	return strings.Join(strList, "; ")
}

// GetCubesList returns the list of cubes in c, each cube as a slice of CanonicalSet
func (c *CanonicalSet) GetCubesList() [][]*interval.CanonicalSet {
	res := [][]*interval.CanonicalSet{}
	if c.dimensions == 1 {
		for k := range c.layers {
			res = append(res, []*interval.CanonicalSet{k})
		}
		return res
	}
	for k, v := range c.layers {
		subRes := v.GetCubesList()
		for _, subList := range subRes {
			cube := []*interval.CanonicalSet{k}
			cube = append(cube, subList...)
			res = append(res, cube)
		}
	}
	return res
}

// SwapDimensions returns a new CanonicalSet object, built from the input CanonicalSet object,
// with dimensions dim1 and dim2 swapped
func (c *CanonicalSet) SwapDimensions(dim1, dim2 int) *CanonicalSet {
	if c.IsEmpty() || dim1 == dim2 {
		return c.Copy()
	}
	if min(dim1, dim2) < 0 || max(dim1, dim2) >= c.dimensions {
		log.Panicf("invalid dimensions: %d, %d", dim1, dim2)
	}
	res := NewCanonicalSet(c.dimensions)
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

func getElementsUnionPerLayer(layers map[*interval.CanonicalSet]*CanonicalSet) map[*interval.CanonicalSet]*CanonicalSet {
	type pair struct {
		hc *CanonicalSet            // hypercube set object
		is []*interval.CanonicalSet // interval-set list
	}
	equivClasses := map[string]*pair{}
	for k, v := range layers {
		if _, ok := equivClasses[v.String()]; ok {
			equivClasses[v.String()].is = append(equivClasses[v.String()].is, k)
		} else {
			equivClasses[v.String()] = &pair{hc: v, is: []*interval.CanonicalSet{k}}
		}
	}
	newLayers := map[*interval.CanonicalSet]*CanonicalSet{}
	for _, p := range equivClasses {
		newVal := p.hc
		newKey := p.is[0]
		for i := 1; i < len(p.is); i += 1 {
			newKey = newKey.Union(p.is[i])
		}
		newLayers[newKey] = newVal
	}
	return newLayers
}

// FromCube returns a new CanonicalSet created from a single input cube
// the input cube is a slice of CanonicalSet, treated as ordered list of dimension values
func FromCube(cube []*interval.CanonicalSet) *CanonicalSet {
	if len(cube) == 0 {
		return nil
	}
	if len(cube) == 1 {
		res := NewCanonicalSet(1)
		res.layers[cube[0].Copy()] = NewCanonicalSet(0)
		return res
	}
	res := NewCanonicalSet(len(cube))
	res.layers[cube[0].Copy()] = FromCube(cube[1:])
	return res
}

// Cube returns a new CanonicalSet created from a single input cube
// the input cube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func Cube(values ...int64) *CanonicalSet {
	cube := []*interval.CanonicalSet{}
	for i := 0; i < len(values); i += 2 {
		cube = append(cube, interval.New(values[i], values[i+1]).ToSet())
	}
	return FromCube(cube)
}
