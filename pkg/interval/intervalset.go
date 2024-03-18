package interval

import (
	"fmt"
	"log"
	"slices"
	"sort"
)

// CanonicalSet is a canonical representation of a set of Interval objects
type CanonicalSet struct {
	intervalSet []Interval // sorted list of non-overlapping intervals
}

func NewCanonicalIntervalSet() *CanonicalSet {
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
	return c.intervalSet[0].Start
}

// IsEmpty returns true if the  CanonicalSet is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.intervalSet) == 0
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

// AddHole updates the current CanonicalSet object by removing the input Interval from the set
func (c *CanonicalSet) AddHole(hole Interval) {
	newIntervalSet := []Interval{}
	for _, interval := range c.intervalSet {
		newIntervalSet = append(newIntervalSet, interval.subtract(hole)...)
	}
	c.intervalSet = newIntervalSet
}

func getNumAsStr(num int64) string {
	return fmt.Sprintf("%v", num)
}

// String returns a string representation of the current CanonicalSet object
func (c *CanonicalSet) String() string {
	if c.IsEmpty() {
		return "Empty"
	}
	res := ""
	for _, interval := range c.intervalSet {
		res += getNumAsStr(interval.Start)
		if interval.Start != interval.End {
			res += "-" + getNumAsStr(interval.End)
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
	i := NewSetFromInterval(New(n, n))
	return i.ContainedIn(c)
}

// ContainedIn returns true of the current CanonicalSet is contained in the input CanonicalSet
func (c *CanonicalSet) ContainedIn(other *CanonicalSet) bool {
	if c == other {
		return true
	}
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
	if c == other {
		return c.Copy()
	}
	res := NewCanonicalIntervalSet()
	for _, interval := range c.intervalSet {
		for _, otherInterval := range other.intervalSet {
			for _, span := range interval.intersection(otherInterval) {
				res.AddInterval(span)
			}
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
			if selfInterval.overlaps(otherInterval) {
				return true
			}
		}
	}
	return false
}

// Subtract returns the subtraction result of input CanonicalSet
func (c *CanonicalSet) Subtract(other *CanonicalSet) *CanonicalSet {
	if c == other {
		return NewCanonicalIntervalSet()
	}
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
	if len(c.intervalSet) == 1 && c.intervalSet[0].Start == c.intervalSet[0].End {
		return true
	}
	return false
}

func NewSetFromInterval(span Interval) *CanonicalSet {
	return &CanonicalSet{intervalSet: []Interval{span}}
}
