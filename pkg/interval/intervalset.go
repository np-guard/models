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
		log.Fatal("cannot take min from empty interval set")
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

// Union updates the CanonicalSet object with the union result of the input CanonicalSet
func (c *CanonicalSet) Union(other *CanonicalSet) {
	if c == other {
		return
	}
	for _, interval := range other.intervalSet {
		c.AddInterval(interval)
	}
}

// Copy returns a new copy of the CanonicalSet object
func (c *CanonicalSet) Copy() *CanonicalSet {
	return &CanonicalSet{intervalSet: slices.Clone(c.intervalSet)}
}

func (c *CanonicalSet) Contains(n int64) bool {
	i := FromInterval(n, n)
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

// Intersect updates current CanonicalSet with intersection result of input CanonicalSet
func (c *CanonicalSet) Intersect(other *CanonicalSet) {
	if c == other {
		return
	}
	newIntervalSet := []Interval{}
	for _, interval := range c.intervalSet {
		for _, otherInterval := range other.intervalSet {
			newIntervalSet = append(newIntervalSet, interval.intersection(otherInterval)...)
		}
	}
	c.intervalSet = newIntervalSet
}

// Overlaps returns true if current CanonicalSet overlaps with input CanonicalSet
func (c *CanonicalSet) Overlaps(other *CanonicalSet) bool {
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

// Subtract updates current CanonicalSet with subtraction result of input CanonicalSet
func (c *CanonicalSet) Subtract(other *CanonicalSet) {
	for _, i := range other.intervalSet {
		c.AddHole(i)
	}
}

func (c *CanonicalSet) IsSingleNumber() bool {
	if len(c.intervalSet) == 1 && c.intervalSet[0].Start == c.intervalSet[0].End {
		return true
	}
	return false
}

func FromInterval(start, end int64) *CanonicalSet {
	return &CanonicalSet{intervalSet: []Interval{{Start: start, End: end}}}
}
