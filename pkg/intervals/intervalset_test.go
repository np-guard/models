package intervals_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/intervals"
)

func TestInterval(t *testing.T) {
	it1 := intervals.Interval{3, 7}
	it2 := intervals.Interval{3, 5}
	it3 := intervals.Interval{2, 7}

	require.Equal(t, "[3-7]", it1.String())

	require.True(t, it3.Lt(it1))
	require.True(t, it2.Lt(it1))
}

func TestIntervalSet(t *testing.T) {
	is1 := intervals.NewCanonicalIntervalSet()
	is1.AddInterval(intervals.Interval{5, 10})
	is1.AddInterval(intervals.Interval{0, 1})
	is1.AddInterval(intervals.Interval{3, 3})
	is1.AddInterval(intervals.Interval{70, 80})
	is1.AddHole(intervals.Interval{7, 9})
	require.True(t, is1.Contains(5))
	require.False(t, is1.Contains(8))

	is2 := intervals.NewCanonicalIntervalSet()
	require.Equal(t, "Empty", is2.String())
	is2.AddInterval(intervals.Interval{6, 8})
	require.Equal(t, "6-8", is2.String())
	require.False(t, is2.IsSingleNumber())
	require.False(t, is2.ContainedIn(*is1))
	require.False(t, is1.ContainedIn(*is2))
	require.False(t, is2.Equal(*is1))
	require.False(t, is1.Equal(*is2))
	require.True(t, is1.Overlaps(is2))
	require.True(t, is2.Overlaps(is1))

	is1.Subtraction(*is2)
	require.False(t, is2.ContainedIn(*is1))
	require.False(t, is1.ContainedIn(*is2))
	require.False(t, is1.Overlaps(is2))
	require.False(t, is2.Overlaps(is1))

	is1.Union(*is2)
	is1.Union(*intervals.CreateFromInterval(7, 9))
	require.True(t, is2.ContainedIn(*is1))
	require.False(t, is1.ContainedIn(*is2))
	require.True(t, is1.Overlaps(is2))
	require.True(t, is2.Overlaps(is1))

	is3 := is1.Copy()
	is3.Intersection(*is2)
	require.True(t, is3.Equal(*is2))
	require.True(t, is2.ContainedIn(is3))

	require.True(t, intervals.CreateFromInterval(1, 1).IsSingleNumber())
}
