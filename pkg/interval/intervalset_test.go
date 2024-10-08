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

func TestIntervalSet(t *testing.T) {
	is1 := interval.NewCanonicalSet()
	is1.AddInterval(interval.New(5, 10))
	is1.AddInterval(interval.New(0, 1))
	is1.AddInterval(interval.New(3, 3))
	is1.AddInterval(interval.New(70, 80))
	is1.AddInterval(interval.New(0, -1))
	is1 = is1.Subtract(interval.New(7, 9).ToSet())
	require.True(t, is1.Contains(5))
	require.False(t, is1.Contains(8))

	is2 := interval.NewCanonicalSet()
	require.True(t, is2.IsEmpty())
	is2.AddInterval(interval.New(6, 8))
	require.Equal(t, []int64{6, 7, 8}, is2.Elements())
	require.False(t, is2.IsSingleNumber())
	require.False(t, is2.IsSubset(is1))
	require.False(t, is1.IsSubset(is2))
	require.False(t, is2.Equal(is1))
	require.False(t, is1.Equal(is2))
	require.True(t, is1.Overlap(is2))
	require.True(t, is2.Overlap(is1))

	is1 = is1.Subtract(is2)
	require.False(t, is2.IsSubset(is1))
	require.False(t, is1.IsSubset(is2))
	require.False(t, is1.Overlap(is2))
	require.False(t, is2.Overlap(is1))

	is1 = is1.Union(is2).Union(interval.New(7, 9).ToSet())
	require.True(t, is2.IsSubset(is1))
	require.False(t, is1.IsSubset(is2))
	require.True(t, is1.Overlap(is2))
	require.True(t, is2.Overlap(is1))

	is3 := is1.Intersect(is2)
	require.True(t, is3.Equal(is2))
	require.True(t, is2.IsSubset(is3))

	require.True(t, interval.New(1, 1).ToSet().IsSingleNumber())
}

func TestIntervalSetSubtract(t *testing.T) {
	s := interval.New(1, 100).ToSet()
	s.AddInterval(interval.New(400, 700))
	d := *interval.New(50, 100).ToSet()
	d.AddInterval(interval.New(400, 700))
	actual := s.Subtract(&d)
	expected := interval.New(1, 49).ToSet()
	require.Equal(t, expected.Elements(), actual.Elements())
}
