package interval

import (
	"fmt"
	"slices"
	"sort"
)

// CanonicalSet is a canonical representation of a set of Interval objects
type CanonicalSet struct {
	IntervalSet []Interval // sorted list of non-overlapping intervals
}

func NewCanonicalIntervalSet() *CanonicalSet {
	return &CanonicalSet{
		IntervalSet: []Interval{},
	}
}

// IsEmpty returns true if the  CanonicalSet is empty
func (c *CanonicalSet) IsEmpty() bool {
	return len(c.IntervalSet) == 0
}

// Equal returns true if the CanonicalSet equals the input CanonicalSet
func (c *CanonicalSet) Equal(other CanonicalSet) bool {
	if len(c.IntervalSet) != len(other.IntervalSet) {
		return false
	}
	for index := range c.IntervalSet {
		if !(c.IntervalSet[index].Equal(other.IntervalSet[index])) {
			return false
		}
	}
	return true
}

// AddInterval adds a new interval range to the set
func (c *CanonicalSet) AddInterval(v Interval) {
	set := c.IntervalSet
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
	c.IntervalSet = slices.Replace(c.IntervalSet, left, right, v)
}

// AddHole updates the current CanonicalSet object by removing the input Interval from the set
func (c *CanonicalSet) AddHole(hole Interval) {
	newIntervalSet := []Interval{}
	for _, interval := range c.IntervalSet {
		newIntervalSet = append(newIntervalSet, interval.subtract(hole)...)
	}
	c.IntervalSet = newIntervalSet
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
	for _, interval := range c.IntervalSet {
		res += getNumAsStr(interval.Start)
		if interval.Start != interval.End {
			res += "-" + getNumAsStr(interval.End)
		}
		res += ","
	}
	return res[:len(res)-1]
}

// Union updates the CanonicalSet object with the union result of the input CanonicalSet
func (c *CanonicalSet) Union(other CanonicalSet) {
	for _, interval := range other.IntervalSet {
		c.AddInterval(interval)
	}
}

// Copy returns a new copy of the CanonicalSet object
func (c *CanonicalSet) Copy() CanonicalSet {
	return CanonicalSet{IntervalSet: append([]Interval(nil), c.IntervalSet...)}
}

func (c *CanonicalSet) Contains(n int64) bool {
	i := FromInterval(n, n)
	return i.ContainedIn(*c)
}

// ContainedIn returns true of the current CanonicalSet is contained in the input CanonicalSet
func (c *CanonicalSet) ContainedIn(other CanonicalSet) bool {
	larger := other.IntervalSet
	for _, target := range c.IntervalSet {
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

// Intersect updates current CanonicalSet with intersection result of input CanonicalSet
func (c *CanonicalSet) Intersect(other CanonicalSet) {
	newIntervalSet := []Interval{}
	for _, interval := range c.IntervalSet {
		for _, otherInterval := range other.IntervalSet {
			newIntervalSet = append(newIntervalSet, interval.intersection(otherInterval)...)
		}
	}
	c.IntervalSet = newIntervalSet
}

// Overlaps returns true if current CanonicalSet overlaps with input CanonicalSet
func (c *CanonicalSet) Overlaps(other *CanonicalSet) bool {
	for _, selfInterval := range c.IntervalSet {
		for _, otherInterval := range other.IntervalSet {
			if selfInterval.overlaps(otherInterval) {
				return true
			}
		}
	}
	return false
}

// Subtract updates current CanonicalSet with subtraction result of input CanonicalSet
func (c *CanonicalSet) Subtract(other CanonicalSet) {
	for _, i := range other.IntervalSet {
		c.AddHole(i)
	}
}

func (c *CanonicalSet) IsSingleNumber() bool {
	if len(c.IntervalSet) == 1 && c.IntervalSet[0].Start == c.IntervalSet[0].End {
		return true
	}
	return false
}

func FromInterval(start, end int64) *CanonicalSet {
	return &CanonicalSet{IntervalSet: []Interval{{Start: start, End: end}}}
}
