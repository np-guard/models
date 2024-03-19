// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package interval

import (
	"fmt"
	"log"
	"slices"
	"sort"
)

// CanonicalSet is a set of int64 integers, implemented using an ordered slice of non-overlapping, non-touching interval
type CanonicalSet struct {
	intervalSet []Interval
}

func NewCanonicalSet() *CanonicalSet {
	return &CanonicalSet{
		intervalSet: []Interval{},
	}
}

func (c *CanonicalSet) Intervals() []Interval {
	return slices.Clone(c.intervalSet)
}

func (c *CanonicalSet) NumIntervals() int {
	return len(c.intervalSet)
}

func (c *CanonicalSet) Min() int64 {
	if len(c.intervalSet) == 0 {
		log.Panic("cannot take min from empty interval set")
	}
	return c.intervalSet[0].Start()
}

// IsEmpty returns true if the CanonicalSet is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.intervalSet) == 0
}

func (c *CanonicalSet) CalculateSize() int64 {
	var res int64 = 0
	for _, r := range c.intervalSet {
		res += r.Size()
	}
	return res
}

// Equal returns true if the CanonicalSet equals the input CanonicalSet
func (c *CanonicalSet) Equal(other *CanonicalSet) bool {
	if c == other {
		return true
	}
	if len(c.intervalSet) != len(other.intervalSet) {
		return false
	}
	for index := range c.intervalSet {
		if !(c.intervalSet[index].Equal(other.intervalSet[index])) {
			return false
		}
	}
	return true
}

// AddInterval adds a new interval range to the set
func (c *CanonicalSet) AddInterval(v Interval) {
	if v.IsEmpty() {
		return
	}
	set := c.intervalSet
	left := sort.Search(len(set), func(i int) bool {
		return set[i].End() >= v.Start()-1
	})
	if left < len(set) && set[left].Start() <= v.End() {
		v = New(min(v.Start(), set[left].Start()), v.End())
	}
	right := sort.Search(len(set), func(j int) bool {
		return set[j].Start() > v.End()+1
	})
	if right > 0 && set[right-1].End() >= v.Start() {
		v = New(v.Start(), max(v.End(), set[right-1].End()))
	}
	c.intervalSet = slices.Replace(c.intervalSet, left, right, v)
}

// AddHole updates the current CanonicalSet object by removing the input Interval from the set
func (c *CanonicalSet) AddHole(hole Interval) {
	newIntervalSet := []Interval{}
	for _, interval := range c.intervalSet {
		newIntervalSet = append(newIntervalSet, interval.subtract(hole)...)
	}
	c.intervalSet = newIntervalSet
}

// String returns a string representation of the current CanonicalSet object
func (c *CanonicalSet) String() string {
	if c.IsEmpty() {
		return "Empty"
	}
	res := ""
	for _, interval := range c.intervalSet {
		if interval.Start() != interval.End() {
			res += fmt.Sprintf("%v-%v", interval.Start(), interval.End())
		} else {
			res += fmt.Sprintf("%v", interval.Start())
		}
		res += ","
	}
	return res[:len(res)-1]
}

// Union returns the union of the two sets
func (c *CanonicalSet) Union(other *CanonicalSet) *CanonicalSet {
	res := c.Copy()
	if c == other {
		return res
	}
	for _, v := range other.intervalSet {
		res.AddInterval(v)
	}
	return res
}

// Copy returns a new copy of the CanonicalSet object
func (c *CanonicalSet) Copy() *CanonicalSet {
	return &CanonicalSet{intervalSet: slices.Clone(c.intervalSet)}
}

func (c *CanonicalSet) Contains(n int64) bool {
	return New(n, n).ToSet().ContainedIn(c)
}

// ContainedIn returns true of the current CanonicalSet is contained in the input CanonicalSet
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) bool {
	if c == other {
		return true
	}
	larger := other.intervalSet
	for _, target := range c.intervalSet {
		left := sort.Search(len(larger), func(i int) bool {
			return larger[i].End() >= target.End()
		})
		if left == len(larger) || larger[left].Start() > target.Start() {
			return false
		}
		// Optimization
		larger = larger[left:]
	}
	return true
}

// Intersect returns the intersection of the current set with the input set
func (c *CanonicalSet) Intersect(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return c.Copy()
	}
	res := NewCanonicalSet()
	for _, left := range c.intervalSet {
		for _, right := range other.intervalSet {
			res.AddInterval(left.intersect(right))
		}
	}
	return res
}

// Overlap returns true if current CanonicalSet overlaps with input CanonicalSet
func (c *CanonicalSet) Overlap(other *CanonicalSet) bool {
	if c == other {
		return !c.IsEmpty()
	}
	for _, selfInterval := range c.intervalSet {
		for _, otherInterval := range other.intervalSet {
			if selfInterval.overlap(otherInterval) {
				return true
			}
		}
	}
	return false
}

// Subtract returns the subtraction result of input CanonicalSet
func (c *CanonicalSet) Subtract(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return NewCanonicalSet()
	}
	res := c.Copy()
	for _, i := range other.intervalSet {
		res.AddHole(i)
	}
	return res
}

func (c *CanonicalSet) IsSingleNumber() bool {
	if len(c.intervalSet) == 1 && c.intervalSet[0].Size() == 1 {
		return true
	}
	return false
}

// Elements returns a slice with all the numbers contained in the set.
// USE WITH CARE. It can easily run out of memory for large sets.
func (c *CanonicalSet) Elements() []int64 {
	// allocate memory up front, to fail early
	res := make([]int64, c.CalculateSize())
	i := 0
	for _, interval := range c.intervalSet {
		for v := interval.Start(); v <= interval.End(); v++ {
			res[i] = v
			i++
		}
	}
	return res
}

func NewSetFromInterval(span Interval) *CanonicalSet {
	return &CanonicalSet{intervalSet: []Interval{span}}
}
