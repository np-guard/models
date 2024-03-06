package intervals

import (
	"fmt"
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

func (c *CanonicalIntervalSet) findIntervalLeft(interval Interval) int {
	if c.IsEmpty() {
		return -1
	}
	low := 0
	high := len(c.IntervalSet)
	for {
		if low == high {
			break
		}
		mid := (low + high) / 2
		if c.IntervalSet[mid].End < interval.Start-1 {
			if mid == len(c.IntervalSet)-1 || c.IntervalSet[mid+1].End >= interval.Start-1 {
				return mid
			}
			low = mid + 1
		} else {
			high = mid
		}
	}
	if low == len(c.IntervalSet) {
		low -= 1
	}
	if c.IntervalSet[low].End >= interval.Start-1 {
		return -1
	}
	return low
}

func (c *CanonicalIntervalSet) findIntervalRight(interval Interval) int {
	if c.IsEmpty() {
		return -1
	}
	low := 0
	high := len(c.IntervalSet)
	for {
		if low == high {
			break
		}
		mid := (low + high) / 2
		if c.IntervalSet[mid].Start > interval.End+1 {
			if mid == 0 || c.IntervalSet[mid-1].Start <= interval.End+1 {
				return mid
			}
			high = mid
		} else {
			low = mid + 1
		}
	}
	if low == len(c.IntervalSet) {
		low -= 1
	}
	if c.IntervalSet[low].Start <= interval.End+1 {
		return -1
	}
	return low
}

func insert(array []Interval, element Interval, i int) []Interval {
	return append(array[:i], append([]Interval{element}, array[i:]...)...)
}

// AddInterval updates the current CanonicalIntervalSet with a new Interval to add
//
//gocyclo:ignore
func (c *CanonicalIntervalSet) AddInterval(intervalToAdd Interval) {
	if c.IsEmpty() {
		c.IntervalSet = append(c.IntervalSet, intervalToAdd)
		return
	}
	left := c.findIntervalLeft(intervalToAdd)
	right := c.findIntervalRight(intervalToAdd)

	// interval_to_add has no overlapping/touching intervals between left to right
	if left >= 0 && right >= 0 && right-left == 1 {
		c.IntervalSet = insert(c.IntervalSet, intervalToAdd, left+1)
		return
	}

	// interval_to_add has no overlapping/touching intervals and is smaller than first interval
	if left == -1 && right == 0 {
		c.IntervalSet = insert(c.IntervalSet, intervalToAdd, 0)
		return
	}

	// interval_to_add has no overlapping/touching intervals and is greater than last interval
	if right == -1 && left == len(c.IntervalSet)-1 {
		c.IntervalSet = append(c.IntervalSet, intervalToAdd)
		return
	}

	// update left/right indexes to be the first potential overlapping/touching intervals from left/right
	left += 1
	if right >= 0 {
		right -= 1
	} else {
		right = len(c.IntervalSet) - 1
	}
	// check which of left/right is overlapping/touching interval_to_add
	leftOverlaps := c.IntervalSet[left].overlaps(intervalToAdd) || c.IntervalSet[left].touches(intervalToAdd)
	rightOverlaps := c.IntervalSet[right].overlaps(intervalToAdd) || c.IntervalSet[right].touches(intervalToAdd)
	newIntervalStart := intervalToAdd.Start
	if leftOverlaps && c.IntervalSet[left].Start < newIntervalStart {
		newIntervalStart = c.IntervalSet[left].Start
	}
	newIntervalEnd := intervalToAdd.End
	if rightOverlaps && c.IntervalSet[right].End > newIntervalEnd {
		newIntervalEnd = c.IntervalSet[right].End
	}
	newInterval := Interval{Start: newIntervalStart, End: newIntervalEnd}
	tmp := c.IntervalSet[right+1:]
	c.IntervalSet = append(c.IntervalSet[:left], newInterval)
	c.IntervalSet = append(c.IntervalSet, tmp...)
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

/*func Union(a, b CanonicalIntervalSet) CanonicalIntervalSet {
	res := a.Copy()
	res.Union(b)
	return res
}*/

func (c *CanonicalIntervalSet) Contains(n int64) bool {
	i := CreateFromInterval(n, n)
	return i.ContainedIn(*c)
}

// ContainedIn returns true of the current CanonicalIntervalSet is contained in the input CanonicalIntervalSet
func (c *CanonicalIntervalSet) ContainedIn(other CanonicalIntervalSet) bool {
	if len(c.IntervalSet) == 1 && len(other.IntervalSet) == 1 {
		return c.IntervalSet[0].isSubset(other.IntervalSet[0])
	}
	for _, interval := range c.IntervalSet {
		left := other.findIntervalLeft(interval)
		if left == len(other.IntervalSet)-1 {
			return false
		}
		if !interval.isSubset(other.IntervalSet[left+1]) {
			return false
		}
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

func (c *CanonicalIntervalSet) Elements() []int {
	res := []int{}
	for _, interval := range c.IntervalSet {
		for i := interval.Start; i <= interval.End; i++ {
			res = append(res, int(i))
		}
	}
	return res
}

func CreateFromInterval(start, end int64) *CanonicalIntervalSet {
	return &CanonicalIntervalSet{IntervalSet: []Interval{{Start: start, End: end}}}
}
