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
	is1 := interval.NewCanonicalIntervalSet()
	is1.AddInterval(interval.Interval{5, 10})
	is1.AddInterval(interval.Interval{0, 1})
	is1.AddInterval(interval.Interval{3, 3})
	is1.AddInterval(interval.Interval{70, 80})
	is1.AddHole(interval.Interval{7, 9})
	require.True(t, is1.Contains(5))
	require.False(t, is1.Contains(8))

	is2 := interval.NewCanonicalIntervalSet()
	require.Equal(t, "Empty", is2.String())
	is2.AddInterval(interval.Interval{6, 8})
	require.Equal(t, "6-8", is2.String())
	require.False(t, is2.IsSingleNumber())
	require.False(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.False(t, is2.Equal(is1))
	require.False(t, is1.Equal(is2))
	require.True(t, is1.Overlaps(is2))
	require.True(t, is2.Overlaps(is1))

	is1.Subtract(is2)
	require.False(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.False(t, is1.Overlaps(is2))
	require.False(t, is2.Overlaps(is1))

	is1.Union(is2)
	is1.Union(interval.FromInterval(7, 9))
	require.True(t, is2.ContainedIn(is1))
	require.False(t, is1.ContainedIn(is2))
	require.True(t, is1.Overlaps(is2))
	require.True(t, is2.Overlaps(is1))

	is3 := is1.Copy()
	is3.Intersect(is2)
	require.True(t, is3.Equal(is2))
	require.True(t, is2.ContainedIn(is3))

	require.True(t, interval.FromInterval(1, 1).IsSingleNumber())
}
