// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package interval

import (
	"errors"
	"fmt"
	"slices"
	"sort"
)

// CanonicalSet is a set of int64 integers, implemented using an ordered slice of non-overlapping, non-touching interval
type CanonicalSet struct {
	intervalSet []Interval
}

func NewCanonicalIntervalSet() *CanonicalSet {
	return &CanonicalSet{
		intervalSet: []Interval{},
	}
}

// IsEmpty returns true if the  CanonicalSet is empty
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
	set := c.intervalSet
	left := sort.Search(len(set), func(i int) bool {
		return set[i].End >= v.Start-1
	})
	if left < len(set) && set[left].Start <= v.End {
		v.Start = min(v.Start, set[left].Start)
	}
	right := sort.Search(len(set), func(j int) bool {
		return set[j].Start > v.End+1
	})
	if right > 0 && set[right-1].End >= v.Start {
		v.End = max(v.End, set[right-1].End)
	}
	c.intervalSet = slices.Replace(c.intervalSet, left, right, v)
}

// String returns a string representation of the current CanonicalSet object
func (c *CanonicalSet) String() string {
	if c.IsEmpty() {
		return "Empty"
	}
	res := ""
	for _, interval := range c.intervalSet {
		if interval.Start != interval.End {
			res += fmt.Sprintf("%v-%v", interval.Start, interval.End)
		} else {
			res += fmt.Sprintf("%v", interval.Start)
		}
		res += ","
	}
	return res[:len(res)-1]
}

// Union returns the union of the two sets
func (c *CanonicalSet) Union(other *CanonicalSet) *CanonicalSet {
	res := c.Copy()
	for _, interval := range other.intervalSet {
		res.AddInterval(interval)
	}
	return res
}

// Copy returns a new copy of the CanonicalSet object
func (c *CanonicalSet) Copy() *CanonicalSet {
	return &CanonicalSet{intervalSet: slices.Clone(c.intervalSet)}
}

func (c *CanonicalSet) Contains(n int64) bool {
	i := CreateSetFromInterval(n, n)
	return i.ContainedIn(c)
}

// ContainedIn returns true of the current interval.CanonicalSet is contained in the input interval.CanonicalSet
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) bool {
	larger := other.intervalSet
	for _, target := range c.intervalSet {
		left := sort.Search(len(larger), func(i int) bool {
			return larger[i].End >= target.End
		})
		if left == len(larger) || larger[left].Start > target.Start {
			return false
		}
		// Optimization
		larger = larger[left:]
	}
	return true
}

// Intersect returns the intersection of the current set with the input set
func (c *CanonicalSet) Intersect(other *CanonicalSet) *CanonicalSet {
	res := NewCanonicalIntervalSet()
	for _, interval := range c.intervalSet {
		for _, otherInterval := range other.intervalSet {
			res.intervalSet = append(res.intervalSet, interval.intersection(otherInterval)...)
		}
	}
	return res
}

// Overlaps returns true if current CanonicalSet overlaps with input CanonicalSet
func (c *CanonicalSet) Overlaps(other *CanonicalSet) bool {
	for _, selfInterval := range c.intervalSet {
		for _, otherInterval := range other.intervalSet {
			if selfInterval.overlaps(otherInterval) {
				return true
			}
		}
	}
	return false
}

// Subtract returns the subtraction result of input CanonicalSet
func (c *CanonicalSet) Subtract(other *CanonicalSet) *CanonicalSet {
	res := slices.Clone(c.intervalSet)
	for _, hole := range other.intervalSet {
		newIntervalSet := []Interval{}
		for _, interval := range res {
			newIntervalSet = append(newIntervalSet, interval.subtract(hole)...)
		}
		res = newIntervalSet
	}
	return &CanonicalSet{
		intervalSet: res,
	}
}

func (c *CanonicalSet) IsSingleNumber() bool {
	if len(c.intervalSet) == 1 && c.intervalSet[0].Size() == 1 {
		return true
	}
	return false
}

func (c *CanonicalSet) Min() (int64, error) {
	if len(c.intervalSet) > 0 {
		return c.intervalSet[0].Start, nil
	}
	return 0, errors.New("empty interval set")
}

// Split returns a set of canonical set objects, each with a single interval
func (c *CanonicalSet) Split() []*CanonicalSet {
	res := make([]*CanonicalSet, len(c.intervalSet))
	for index, ipr := range c.intervalSet {
		res[index] = CreateSetFromInterval(ipr.Start, ipr.End)
	}
	return res
}

func (c *CanonicalSet) Intervals() []Interval {
	return slices.Clone(c.intervalSet)
}

func (c *CanonicalSet) NumIntervals() int {
	return len(c.intervalSet)
}

func (c *CanonicalSet) Elements() []int64 {
	res := []int64{}
	for _, interval := range c.intervalSet {
		for i := interval.Start; i <= interval.End; i++ {
			res = append(res, i)
		}
	}
	return res
}

func CreateSetFromInterval(start, end int64) *CanonicalSet {
	return &CanonicalSet{intervalSet: []Interval{{Start: start, End: end}}}
}
