package interval

import (
	"fmt"
	"slices"
	"sort"
)

// CanonicalIntervalSet is a canonical representation of a set of Interval objects
type CanonicalIntervalSet struct {
	IntervalSet []Interval // sorted list of non-overlapping intervals
}

func NewCanonicalIntervalSet() *CanonicalIntervalSet {
	return &CanonicalIntervalSet{
		IntervalSet: []Interval{},
	}
}

// IsEmpty returns true if the  CanonicalIntervalSet is empty
func (c *CanonicalIntervalSet) IsEmpty() bool {
	return len(c.IntervalSet) == 0
}

// Equal returns true if the CanonicalIntervalSet equals the input CanonicalIntervalSet
func (c *CanonicalIntervalSet) Equal(other CanonicalIntervalSet) bool {
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
func (c *CanonicalIntervalSet) AddInterval(v Interval) {
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

// AddHole updates the current CanonicalIntervalSet object by removing the input Interval from the set
func (c *CanonicalIntervalSet) AddHole(hole Interval) {
	newIntervalSet := []Interval{}
	for _, interval := range c.IntervalSet {
		newIntervalSet = append(newIntervalSet, interval.subtract(hole)...)
	}
	c.IntervalSet = newIntervalSet
}

func getNumAsStr(num int64) string {
	return fmt.Sprintf("%v", num)
}

// String returns a string representation of the current CanonicalIntervalSet object
func (c *CanonicalIntervalSet) String() string {
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

// Union updates the CanonicalIntervalSet object with the union result of the input CanonicalIntervalSet
func (c *CanonicalIntervalSet) Union(other CanonicalIntervalSet) {
	for _, interval := range other.IntervalSet {
		c.AddInterval(interval)
	}
}

// Copy returns a new copy of the CanonicalIntervalSet object
func (c *CanonicalIntervalSet) Copy() CanonicalIntervalSet {
	return CanonicalIntervalSet{IntervalSet: append([]Interval(nil), c.IntervalSet...)}
}

func (c *CanonicalIntervalSet) Contains(n int64) bool {
	i := CreateFromInterval(n, n)
	return i.ContainedIn(*c)
}

// ContainedIn returns true of the current CanonicalIntervalSet is contained in the input CanonicalIntervalSet
func (c *CanonicalIntervalSet) ContainedIn(other CanonicalIntervalSet) bool {
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

// Intersection updates current CanonicalIntervalSet with intersection result of input CanonicalIntervalSet
func (c *CanonicalIntervalSet) Intersection(other CanonicalIntervalSet) {
	newIntervalSet := []Interval{}
	for _, interval := range c.IntervalSet {
		for _, otherInterval := range other.IntervalSet {
			newIntervalSet = append(newIntervalSet, interval.intersection(otherInterval)...)
		}
	}
	c.IntervalSet = newIntervalSet
}

// Overlaps returns true if current CanonicalIntervalSet overlaps with input CanonicalIntervalSet
func (c *CanonicalIntervalSet) Overlaps(other *CanonicalIntervalSet) bool {
	for _, selfInterval := range c.IntervalSet {
		for _, otherInterval := range other.IntervalSet {
			if selfInterval.overlaps(otherInterval) {
				return true
			}
		}
	}
	return false
}

// Subtraction updates current CanonicalIntervalSet with subtraction result of input CanonicalIntervalSet
func (c *CanonicalIntervalSet) Subtraction(other CanonicalIntervalSet) {
	for _, i := range other.IntervalSet {
		c.AddHole(i)
	}
}

func (c *CanonicalIntervalSet) IsSingleNumber() bool {
	if len(c.IntervalSet) == 1 && c.IntervalSet[0].Start == c.IntervalSet[0].End {
		return true
	}
	return false
}

func CreateFromInterval(start, end int64) *CanonicalIntervalSet {
	return &CanonicalIntervalSet{IntervalSet: []Interval{{Start: start, End: end}}}
}
