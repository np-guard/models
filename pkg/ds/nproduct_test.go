// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
)

type Hypercube = ds.NProduct[*interval.CanonicalSet]

// ncube returns a new ds.NProduct created from a single input ncube
// the input ncube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func ncube(values ...int64) *Hypercube {
	partition := []*interval.CanonicalSet{}
	for i := 0; i < len(values); i += 2 {
		partition = append(partition, interval.NewSetFromInterval(interval.New(values[i], values[i+1])))
	}
	return ds.CartesianN(partition)
}

func union[S ds.Set[S]](set S, sets ...S) S {
	for _, c := range sets {
		set = set.Union(c)
	}
	return set.Copy()
}

func TestHCBasic(t *testing.T) {
	a := ncube(1, 100)
	b := ncube(1, 100)
	c := ncube(1, 200)
	d := ncube(1, 100, 1, 100)
	e := ncube(1, 100, 1, 100)
	f := ncube(1, 100, 1, 200)

	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))

	require.False(t, a.Equal(c))
	require.False(t, c.Equal(a))

	require.False(t, a.Equal(d))
	require.False(t, d.Equal(a))

	require.True(t, d.Equal(e))
	require.True(t, e.Equal(d))

	require.False(t, d.Equal(f))
	require.False(t, f.Equal(d))
}

func TestCopy(t *testing.T) {
	a := ncube(1, 100)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

// func TestString(t *testing.T) {
// 	require.Equal(t, "[(1-3)]", ncube(1, 3).String())
// 	require.Equal(t, "[(1-3),(2-4)]", ncube(1, 3, 2, 4).String())
// }

// func TestOr(t *testing.T) {
// 	a := ncube(1, 100, 1, 100)
// 	b := ncube(1, 90, 1, 200)
// 	c := a.Union(b)
// 	require.Equal(t, "[(1-90),(1-200)]; [(91-100),(1-100)]", c.String())
// }

func TestBasic1(t *testing.T) {
	a := union(
		ncube(1, 2),
		ncube(5, 6),
		ncube(3, 4),
	)
	b := ncube(1, 6)
	require.True(t, a.Equal(b))
}

func TestBasic2(t *testing.T) {
	a := union(
		ncube(1, 2, 1, 5),
		ncube(1, 2, 7, 9),
		ncube(1, 2, 6, 7),
	)
	b := ncube(1, 2, 1, 9)
	require.True(t, a.Equal(b))
}

// func TestNew(t *testing.T) {
// 	a := union(
// 		ncube(10, 20, 10, 20, 1, 65535),
// 		ncube(1, 65535, 15, 40, 1, 65535),
// 		ncube(1, 65535, 100, 200, 30, 80),
// 	)
// 	expectedStr := "[(1-9,21-65535),(100-200),(30-80)]; " +
// 		"[(1-9,21-65535),(15-40),(1-65535)]; " +
// 		"[(10-20),(10-40),(1-65535)]; " +
// 		"[(10-20),(100-200),(30-80)]"
// 	require.Equal(t, expectedStr, a.String())
// }

func checkContained[S ds.Set[S]](t *testing.T, a, b S, expected bool) {
	t.Helper()
	contained := a.IsSubset(b)
	require.Equal(t, expected, contained)
}

func TestIsSubset(t *testing.T) {
	a := ncube(1, 100, 200, 300)
	b := ncube(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = b.Union(ncube(10, 200, 210, 280))
	checkContained(t, b, a, false)
}

func TestIsSubset1(t *testing.T) {
	checkContained(t, ncube(1, 3), ncube(2, 4), false)
	checkContained(t, ncube(2, 4), ncube(1, 3), false)
	checkContained(t, ncube(1, 3), ncube(1, 4), true)
	checkContained(t, ncube(1, 4), ncube(1, 3), false)
}

func TestIsSubset2(t *testing.T) {
	c := union(
		ncube(1, 100, 200, 300),
		ncube(150, 180, 20, 300),
		ncube(200, 240, 200, 300),
		ncube(241, 300, 200, 350),
	)

	a := union(
		ncube(1, 100, 200, 300),
		ncube(150, 180, 20, 300),
		ncube(200, 240, 200, 300),
		ncube(242, 300, 200, 350),
	)
	d := ncube(210, 220, 210, 280)
	e := ncube(210, 310, 210, 280)
	f := ncube(210, 250, 210, 280)
	f1 := ncube(210, 240, 210, 280)
	f2 := ncube(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestIsSubset3(t *testing.T) {
	a := ncube(105, 105, 54, 54)
	b := union(
		ncube(0, 204, 0, 255),
		ncube(205, 205, 0, 53),
		ncube(205, 205, 55, 255),
		ncube(206, 254, 0, 255),
	)
	checkContained(t, a, b, true)
}

func TestIsSubset4(t *testing.T) {
	a := ncube(105, 105, 54, 54)
	b := ncube(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestIsSubset5(t *testing.T) {
	a := ncube(100, 200, 54, 65, 60, 300)
	b := ncube(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestEqual1(t *testing.T) {
	a := ncube(1, 2)
	b := ncube(1, 2)
	require.True(t, a.Equal(b))

	c := ncube(1, 2, 1, 5)
	d := ncube(1, 2, 1, 5)
	require.True(t, c.Equal(d))
}

func TestEqual2(t *testing.T) {
	c := union(
		ncube(1, 2, 1, 5),
		ncube(1, 2, 7, 9),
		ncube(1, 2, 6, 7),
		ncube(4, 8, 1, 9),
	)
	res := union(
		ncube(4, 8, 1, 9),
		ncube(1, 2, 1, 9),
	)
	require.True(t, res.Equal(c))

	d := union(
		ncube(1, 2, 1, 5),
		ncube(5, 6, 1, 5),
		ncube(3, 4, 1, 5),
	)
	res2 := ncube(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestBasicAddCube(t *testing.T) {
	a := union(
		ncube(1, 2),
		ncube(8, 10),
	)
	b := union(
		a,
		ncube(1, 2),
		ncube(6, 10),
		ncube(1, 10),
	)
	res := ncube(1, 10)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

// func TestFourHoles(t *testing.T) {
// 	a := ncube(1, 2, 1, 2)
// 	require.Equal(t, "[(1-2),(1-2)]", a.String())
// 	require.Equal(t, "[(1),(2)]; [(2),(1-2)]", a.Subtract(ncube(1, 1, 1, 1)).String())
// 	require.Equal(t, "[(1),(1)]; [(2),(1-2)]", a.Subtract(ncube(1, 1, 2, 2)).String())
// 	require.Equal(t, "[(1),(1-2)]; [(2),(2)]", a.Subtract(ncube(2, 2, 1, 1)).String())
// 	require.Equal(t, "[(1),(1-2)]; [(2),(1)]", a.Subtract(ncube(2, 2, 2, 2)).String())
// }

func TestBasicSubtract1(t *testing.T) {
	a := ncube(1, 10)
	require.True(t, a.Subtract(ncube(3, 7)).Equal(union(ncube(1, 2), ncube(8, 10))))
	require.True(t, a.Subtract(ncube(3, 20)).Equal(ncube(1, 2)))
	require.True(t, a.Subtract(ncube(0, 20)).IsEmpty())
	require.True(t, a.Subtract(ncube(0, 5)).Equal(ncube(6, 10)))
	require.True(t, a.Subtract(ncube(12, 14)).Equal(ncube(1, 10)))
}

func TestBasicSubtract2(t *testing.T) {
	a := ncube(1, 100, 200, 300).Subtract(ncube(50, 60, 220, 300))
	resA := union(
		ncube(61, 100, 200, 300),
		ncube(50, 60, 200, 219),
		ncube(1, 49, 200, 300),
	)
	require.True(t, a.Equal(resA))

	b := ncube(1, 100, 200, 300).Subtract(ncube(50, 1000, 0, 250))
	resB := union(
		ncube(50, 100, 251, 300),
		ncube(1, 49, 200, 300),
	)
	require.True(t, b.Equal(resB))

	c := union(
		ncube(1, 100, 200, 300),
		ncube(400, 700, 200, 300),
	).Subtract(ncube(50, 1000, 0, 250))
	resC := union(
		ncube(50, 100, 251, 300),
		ncube(1, 49, 200, 300),
		ncube(400, 700, 251, 300),
	)
	require.True(t, c.Equal(resC))

	d := ncube(1, 100, 200, 300).Subtract(ncube(50, 60, 220, 300))
	dRes := union(
		ncube(1, 49, 200, 300),
		ncube(50, 60, 200, 219),
		ncube(61, 100, 200, 300),
	)
	require.True(t, d.Equal(dRes))
}

func TestAddHole2(t *testing.T) {
	c := union(
		ncube(80, 100, 20, 300),
		ncube(250, 400, 20, 300),
	).Subtract(ncube(30, 300, 100, 102))
	d := union(
		ncube(80, 100, 20, 99),
		ncube(80, 100, 103, 300),
		ncube(250, 300, 20, 99),
		ncube(250, 300, 103, 300),
		ncube(301, 400, 20, 300),
	)
	require.True(t, c.Equal(d))
}

func TestSubtractToEmpty(t *testing.T) {
	c := ncube(1, 100, 200, 300).Subtract(ncube(1, 100, 200, 300))
	require.True(t, c.IsEmpty())
}

func TestUnion1(t *testing.T) {
	c := union(
		ncube(1, 100, 200, 300),
		ncube(101, 200, 200, 300),
	)
	cExpected := ncube(1, 200, 200, 300)
	require.True(t, cExpected.Equal(c))
}

func TestUnion2(t *testing.T) {
	c := union(
		ncube(1, 100, 200, 300),
		ncube(101, 200, 200, 300),
		ncube(201, 300, 200, 300),
		ncube(301, 400, 200, 300),
		ncube(402, 500, 200, 300),
		ncube(500, 600, 200, 700),
		ncube(601, 700, 200, 700),
	)
	cExpected := union(
		ncube(1, 400, 200, 300),
		ncube(402, 500, 200, 300),
		ncube(500, 700, 200, 700),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(ncube(702, 800, 200, 700))
	dExpected := cExpected.Union(ncube(702, 800, 200, 700))
	require.True(t, d.Equal(dExpected))
}

func TestIntersect(t *testing.T) {
	c := ncube(5, 15, 3, 10).Intersect(ncube(8, 30, 7, 20))
	d := ncube(8, 15, 7, 10)
	require.True(t, c.Equal(d))
}

func TestUnionMerge(t *testing.T) {
	a := union(
		ncube(5, 15, 3, 6),
		ncube(5, 30, 7, 10),
		ncube(8, 30, 11, 20),
	)
	excepted := union(
		ncube(5, 15, 3, 10),
		ncube(8, 30, 7, 20),
	)
	require.True(t, excepted.Equal(a))
}

func TestSubtract(t *testing.T) {
	g := ncube(5, 15, 3, 10).Subtract(ncube(8, 30, 7, 20))
	h := union(
		ncube(5, 7, 3, 10),
		ncube(8, 15, 3, 6),
	)
	require.True(t, g.Equal(h))
}

func TestIntersectEmpty(t *testing.T) {
	a := ncube(5, 15, 3, 10)
	b := union(
		ncube(1, 3, 7, 20),
		ncube(20, 23, 7, 20),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestOr2(t *testing.T) {
	a := union(
		ncube(1, 79, 10054, 10054),
		ncube(80, 100, 10053, 10054),
		ncube(101, 65535, 10054, 10054),
	)
	expected := union(
		ncube(80, 100, 10053, 10053),
		ncube(1, 65535, 10054, 10054),
	)
	require.True(t, expected.Equal(a))
}

func TestSwapDimensions(t *testing.T) {
	s := ds.NewNProduct[*interval.CanonicalSet](2)
	require.True(t, s.Swap(0, 1).Equal(s))

	require.True(t, ncube(1, 2).Swap(0, 0).Equal(ncube(1, 2)))

	require.True(t, ncube(1, 2, 3, 4).Swap(0, 1).Equal(ncube(3, 4, 1, 2)))
	require.True(t, ncube(1, 2, 1, 2).Swap(0, 1).Equal(ncube(1, 2, 1, 2)))

	require.True(t, ncube(1, 2, 3, 4, 5, 6).Swap(0, 1).Equal(ncube(3, 4, 1, 2, 5, 6)))
	require.True(t, ncube(1, 2, 3, 4, 5, 6).Swap(1, 2).Equal(ncube(1, 2, 5, 6, 3, 4)))
	require.True(t, ncube(1, 2, 3, 4, 5, 6).Swap(0, 2).Equal(ncube(5, 6, 3, 4, 1, 2)))

	require.True(t, union(
		ncube(1, 3, 7, 20),
		ncube(20, 23, 7, 20),
	).Swap(0, 1).Equal(union(
		ncube(7, 20, 1, 3),
		ncube(7, 20, 20, 23),
	)))

	require.Panics(t, func() { ncube(1, 2).Swap(0, 1) })
	require.Panics(t, func() { ncube(1, 2).Swap(-1, 0) })
}
