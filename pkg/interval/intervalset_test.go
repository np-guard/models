// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package interval_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/interval"
)

func TestInterval(t *testing.T) {
	it1 := interval.Interval{3, 7}

	require.Equal(t, "[3-7]", it1.String())
}

func TestIntervalSet(t *testing.T) {
	is1 := interval.NewCanonicalSet()
	is1.AddInterval(interval.Interval{5, 10})
	is1.AddInterval(interval.Interval{0, 1})
	is1.AddInterval(interval.Interval{3, 3})
	is1.AddInterval(interval.Interval{70, 80})
	is1 = is1.Subtract(interval.New(7, 9).ToSet())
	require.True(t, is1.Contains(5))
	require.False(t, is1.Contains(8))

	is2 := interval.NewCanonicalSet()
	require.Equal(t, "Empty", is2.String())
	is2.AddInterval(interval.Interval{6, 8})
	require.Equal(t, "6-8", is2.String())
	require.False(t, is2.IsSingleNumber())
	require.False(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.False(t, is2.Equal(is1))
	require.False(t, is1.Equal(is2))
	require.True(t, is1.Overlap(is2))
	require.True(t, is2.Overlap(is1))

	is1 = is1.Subtract(is2)
	require.False(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.False(t, is1.Overlap(is2))
	require.False(t, is2.Overlap(is1))

	is1 = is1.Union(is2).Union(interval.New(7, 9).ToSet())
	require.True(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.True(t, is1.Overlap(is2))
	require.True(t, is2.Overlap(is1))

	is3 := is1.Intersect(is2)
	require.True(t, is3.Equal(is2))
	require.True(t, is2.ContainedIn(is3))

	require.True(t, interval.New(1, 1).ToSet().IsSingleNumber())
}

func TestIntervalSetSubtract(t *testing.T) {
	s := interval.New(1, 100).ToSet()
	s.AddInterval(interval.Interval{Start: 400, End: 700})
	d := *interval.New(50, 100).ToSet()
	d.AddInterval(interval.Interval{Start: 400, End: 700})
	actual := s.Subtract(&d)
	expected := interval.New(1, 49).ToSet()
	require.Equal(t, expected.String(), actual.String())
}
