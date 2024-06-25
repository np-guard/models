/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package interval

import "fmt"

// Interval is an integer interval from start to end inclusive
type Interval struct {
	start int64
	end   int64
}

func New(start, end int64) Interval {
	if end < start {
		return Interval{start: 0, end: -1}
	}
	return Interval{start: start, end: end}
}

func (i Interval) Start() int64 {
	return i.start
}

func (i Interval) End() int64 {
	return i.end
}

// String returns a String representation of Interval object
func (i Interval) String() string {
	if i.IsEmpty() {
		return "[]"
	}
	return fmt.Sprintf("[%v-%v]", i.start, i.end)
}

// ShortString returns a compacted String representation of Interval object:
// "v" instead of "v-v", without braces
func (i Interval) ShortString() string {
	if i.IsEmpty() {
		return ""
	}
	if i.start == i.end {
		return fmt.Sprintf("%v", i.start)
	}
	return fmt.Sprintf("%v-%v", i.start, i.end)
}

// IsEmpty returns true if the interval is empty, false otherwise.
// An interval is considered empty if its start is greater than its end.
func (i Interval) IsEmpty() bool {
	return i.end < i.start
}

// Equal returns true if current Interval obj is equal to the input Interval
func (i Interval) Equal(x Interval) bool {
	return i.start == x.start && i.end == x.end
}

func (i Interval) Size() int64 {
	return i.end - i.start + 1
}

func (i Interval) overlap(other Interval) bool {
	if i.IsEmpty() {
		return false
	}
	return other.end >= i.start && other.start <= i.end
}

func (i Interval) isSubset(other Interval) bool {
	if i.IsEmpty() {
		return true
	}
	return other.start <= i.start && other.end >= i.end
}

// returns a list with up to 2 intervals
func (i Interval) subtract(other Interval) []Interval {
	if !i.overlap(other) {
		return []Interval{i}
	}
	if i.isSubset(other) {
		return []Interval{}
	}
	if i.start < other.start && i.end > other.end {
		// self is split into two ranges by other
		return []Interval{{start: i.start, end: other.start - 1}, {start: other.end + 1, end: i.end}}
	}
	if i.start < other.start {
		return []Interval{{start: i.start, end: min(i.end, other.start-1)}}
	}
	return []Interval{{start: max(i.start, other.end+1), end: i.end}}
}

func (i Interval) intersect(other Interval) Interval {
	return New(
		max(i.start, other.start),
		min(i.end, other.end),
	)
}

func (i Interval) Elements() []int64 {
	var res []int64
	for v := i.start; v <= i.end; v++ {
		res = append(res, v)
	}
	return res
}

func (i Interval) ToSet() *CanonicalSet {
	return NewSetFromInterval(i)
}
