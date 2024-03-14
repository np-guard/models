package interval

import "fmt"

// Interval is an integer interval from Start to End
type Interval struct {
	Start int64
	End   int64
}

// String returns a String representation of Interval object
func (i Interval) String() string {
	return fmt.Sprintf("[%v-%v]", i.Start, i.End)
}

// Equal returns true if current Interval obj is equal to the input Interval
func (i Interval) Equal(x Interval) bool {
	return i.Start == x.Start && i.End == x.End
}

func (i Interval) overlaps(other Interval) bool {
	return other.End >= i.Start && other.Start <= i.End
}

func (i Interval) isSubset(other Interval) bool {
	return other.Start <= i.Start && other.End >= i.End
}

// returns a list with up to 2 intervals
func (i Interval) subtract(other Interval) []Interval {
	if !i.overlaps(other) {
		return []Interval{i}
	}
	if i.isSubset(other) {
		return []Interval{}
	}
	if i.Start < other.Start && i.End > other.End {
		// self is split into two ranges by other
		return []Interval{{Start: i.Start, End: other.Start - 1}, {Start: other.End + 1, End: i.End}}
	}
	if i.Start < other.Start {
		return []Interval{{Start: i.Start, End: min(i.End, other.Start-1)}}
	}
	return []Interval{{Start: max(i.Start, other.End+1), End: i.End}}
}

func (i Interval) intersection(other Interval) []Interval {
	maxStart := max(i.Start, other.Start)
	minEnd := min(i.End, other.End)
	if minEnd < maxStart {
		return []Interval{}
	}
	return []Interval{{Start: maxStart, End: minEnd}}
}
