/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package interval_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/interval"
)

func span(start, end int64) interval.Interval {
	return interval.New(start, end)
}

func empty() interval.Interval {
	return interval.New(0, -1)
}

// To avoid recursion, only use helper functions that are declared earlier in the file.

func requireEqual(t *testing.T, actual, expected interval.Interval) {
	t.Helper()
	require.Equal(t, actual.Start(), expected.Start())
	require.Equal(t, actual.End(), expected.End())
}

func requireIntersection(t *testing.T, i1, i2, expected interval.Interval) {
	t.Helper()
	require.Equal(t, !expected.IsEmpty(), i1.Overlap(i2))
	require.Equal(t, !expected.IsEmpty(), i2.Overlap(i1))
	requireEqual(t, i1.Intersect(i2), expected)
	requireEqual(t, i2.Intersect(i1), expected)
}

func requireUnrelated(t *testing.T, i1, i2 interval.Interval) {
	t.Helper()
	require.False(t, i1.IsSubset(i2))
	require.False(t, i2.IsSubset(i1))
	require.True(t, i1.Intersect(i2).Size() < min(i1.Size(), i2.Size()))
}

func requireSubset(t *testing.T, small, large interval.Interval) {
	t.Helper()
	require.True(t, small.IsSubset(large))
	requireIntersection(t, small, large, small)
	require.Equal(t, small.Equal(large), large.IsSubset(small))
}

func TestInterval_Elements(t *testing.T) {
	it1 := span(3, 7)

	require.Equal(t, int64(5), it1.Size())
	require.Equal(t, []int64{3, 4, 5, 6, 7}, it1.Elements())
}

func TestInterval_Empty(t *testing.T) {
	requireEqual(t, empty(), empty())
	requireEqual(t, span(5, 4), empty())

	require.Equal(t, int64(0), span(-1, -2).Size())
	require.Equal(t, []int64{}, span(-1, -2).Elements())
}

func TestInterval_Intersect(t *testing.T) {
	s := span(4, 6)
	requireIntersection(t, s, s, span(4, 6))

	// requireIntersection checks both directions; no need to add tests for all combinations
	requireIntersection(t, empty(), empty(), empty())
	requireIntersection(t, span(4, 6), span(3, 7), span(4, 6))
	requireIntersection(t, span(4, 6), span(3, 5), span(4, 5))
	requireIntersection(t, span(4, 6), empty(), empty())
	requireIntersection(t, span(4, 6), span(7, 8), empty())
}

func TestInterval_IsSubset(t *testing.T) {
	requireSubset(t, empty(), empty())
	requireSubset(t, empty(), span(1, 2))
	requireSubset(t, span(1, 2), span(1, 2))
	requireSubset(t, span(1, 2), span(0, 3))

	requireUnrelated(t, span(1, 2), span(2, 3))
	requireUnrelated(t, span(1, 2), span(3, 4))
	requireUnrelated(t, span(1, 2), span(5, 6))
}

func TestInterval_SubtractSplit(t *testing.T) {
	require.Equal(t, []interval.Interval{}, empty().SubtractSplit(empty()))
	require.Equal(t, []interval.Interval{}, empty().SubtractSplit(span(1, 3)))
	require.Equal(t, []interval.Interval{}, span(1, 3).SubtractSplit(span(1, 3)))
	require.Equal(t, []interval.Interval{}, span(1, 3).SubtractSplit(span(0, 3)))
	require.Equal(t, []interval.Interval{}, span(1, 3).SubtractSplit(span(1, 4)))
	require.Equal(t, []interval.Interval{}, span(1, 3).SubtractSplit(span(0, 4)))

	require.Equal(t, []interval.Interval{span(1, 1), span(3, 3)}, span(1, 3).SubtractSplit(span(2, 2)))
	require.Equal(t, []interval.Interval{span(1, 2), span(5, 6)}, span(1, 6).SubtractSplit(span(3, 4)))

	require.Equal(t, []interval.Interval{span(3, 4)}, span(3, 4).SubtractSplit(span(1, 2)))
	require.Equal(t, []interval.Interval{span(3, 4)}, span(3, 4).SubtractSplit(span(5, 6)))
}
