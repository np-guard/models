/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
)

// cubioid returns a new ds.NProduct created from a single input cubioid
// the input cubioid is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func cubioid(s1, e1, s2, e2, s3, e3 int64) ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianRightTriple(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet(), interval.New(s3, e3).ToSet())
}

func TestCubioidEqual(t *testing.T) {
	d := cubioid(1, 100, 1, 100, 2, 100)
	e := cubioid(1, 100, 1, 100, 2, 100)
	f := cubioid(1, 100, 1, 200, 2, 100)

	require.True(t, d.Equal(e))
	require.True(t, e.Equal(d))

	require.False(t, d.Equal(f))
	require.False(t, f.Equal(d))
}

func TestCubioidCopy(t *testing.T) {
	a := cubioid(1, 100, 3, 4, 5, 6)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestCubioidBasic1(t *testing.T) {
	a := union(
		cubioid(1, 2, 3, 4, 5, 6),
		cubioid(5, 6, 3, 4, 5, 6),
		cubioid(3, 4, 3, 4, 5, 6),
	)
	b := cubioid(1, 6, 3, 4, 5, 6)
	require.True(t, a.Equal(b))
}

func TestCubioidBasic2(t *testing.T) {
	a := union(
		cubioid(1, 2, 1, 5, 0, 3),
		cubioid(1, 2, 7, 9, 0, 3),
		cubioid(1, 2, 6, 7, 0, 3),
	)
	b := cubioid(1, 2, 1, 9, 0, 3)
	require.True(t, a.Equal(b))
}

func TestCubioidIsSubset(t *testing.T) {
	a := cubioid(1, 100, 200, 300, 0, 3)
	b := cubioid(10, 80, 210, 280, 0, 3)
	checkContained(t, b, a, true)
	b = b.Union(cubioid(10, 200, 210, 280, 0, 3))
	checkContained(t, b, a, false)
}

func TestCubioidIsSubset1(t *testing.T) {
	checkContained(t, cubioid(1, 3, 0, 3, 3, 5), cubioid(2, 4, 0, 3, 3, 5), false)
	checkContained(t, cubioid(2, 4, 0, 3, 3, 5), cubioid(1, 3, 0, 3, 3, 5), false)
	checkContained(t, cubioid(1, 3, 0, 3, 3, 5), cubioid(1, 4, 0, 3, 3, 5), true)
	checkContained(t, cubioid(1, 4, 0, 3, 3, 5), cubioid(1, 3, 0, 3, 3, 5), false)
}

func TestCubioidIsSubset2(t *testing.T) {
	c := union(
		cubioid(1, 100, 200, 300, 3, 5),
		cubioid(150, 180, 20, 300, 3, 5),
		cubioid(200, 240, 200, 300, 3, 5),
		cubioid(241, 300, 200, 350, 3, 5),
	)

	a := union(
		cubioid(1, 100, 200, 300, 3, 5),
		cubioid(150, 180, 20, 300, 3, 5),
		cubioid(200, 240, 200, 300, 3, 5),
		cubioid(242, 300, 200, 350, 3, 5),
	)
	d := cubioid(210, 220, 210, 280, 3, 5)
	e := cubioid(210, 310, 210, 280, 3, 5)
	f := cubioid(210, 250, 210, 280, 3, 5)
	f1 := cubioid(210, 240, 210, 280, 3, 5)
	f2 := cubioid(241, 250, 210, 280, 3, 5)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestCubioidIsSubset3(t *testing.T) {
	a := cubioid(105, 105, 54, 54, 3, 5)
	b := union(
		cubioid(0, 204, 0, 255, 3, 5),
		cubioid(205, 205, 0, 53, 3, 5),
		cubioid(205, 205, 55, 255, 3, 5),
		cubioid(206, 254, 0, 255, 3, 5),
	)
	checkContained(t, a, b, true)
}

func TestCubioidIsSubset4(t *testing.T) {
	a := cubioid(105, 105, 54, 54, 3, 5)
	b := cubioid(200, 204, 0, 255, 3, 5)
	checkContained(t, a, b, false)
}

func TestCubioidIsSubset5(t *testing.T) {
	a := cubioid(100, 200, 54, 65, 60, 300)
	b := cubioid(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestCubioidEqual1(t *testing.T) {
	a := cubioid(1, 2, 3, 5, 3, 5)
	b := cubioid(1, 2, 3, 5, 3, 5)
	require.True(t, a.Equal(b))

	c := cubioid(1, 2, 1, 5, 3, 5)
	d := cubioid(1, 2, 1, 5, 3, 5)
	require.True(t, c.Equal(d))
}

func TestCubioidEqual2(t *testing.T) {
	c := union(
		cubioid(1, 2, 1, 5, 3, 5),
		cubioid(1, 2, 7, 9, 3, 5),
		cubioid(1, 2, 6, 7, 3, 5),
		cubioid(4, 8, 1, 9, 3, 5),
	)
	res := union(
		cubioid(4, 8, 1, 9, 3, 5),
		cubioid(1, 2, 1, 9, 3, 5),
	)
	require.True(t, res.Equal(c))

	d := union(
		cubioid(1, 2, 1, 5, 3, 5),
		cubioid(5, 6, 1, 5, 3, 5),
		cubioid(3, 4, 1, 5, 3, 5),
	)
	res2 := cubioid(1, 6, 1, 5, 3, 5)
	require.True(t, res2.Equal(d))
}

func TestCubioidBasicAddCubioid(t *testing.T) {
	a := union(
		cubioid(1, 2, 3, 5, 3, 5),
		cubioid(8, 10, 3, 5, 3, 5),
	)
	b := union(
		a,
		cubioid(1, 2, 3, 5, 3, 5),
		cubioid(6, 10, 3, 5, 3, 5),
		cubioid(1, 10, 3, 5, 3, 5),
	)
	res := cubioid(1, 10, 3, 5, 3, 5)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestCubioidBasicSubtract1(t *testing.T) {
	a := cubioid(1, 10, 3, 5, 3, 5)
	require.True(t, a.Subtract(cubioid(3, 7, 3, 5, 3, 5)).Equal(union(cubioid(1, 2, 3, 5, 3, 5), cubioid(8, 10, 3, 5, 3, 5))))
	require.True(t, a.Subtract(cubioid(3, 20, 3, 5, 3, 5)).Equal(cubioid(1, 2, 3, 5, 3, 5)))
	require.True(t, a.Subtract(cubioid(0, 20, 3, 5, 3, 5)).IsEmpty())
	require.True(t, a.Subtract(cubioid(0, 5, 3, 5, 3, 5)).Equal(cubioid(6, 10, 3, 5, 3, 5)))
	require.True(t, a.Subtract(cubioid(12, 14, 3, 5, 3, 5)).Equal(cubioid(1, 10, 3, 5, 3, 5)))
}

func TestCubioidBasicSubtract2(t *testing.T) {
	a := cubioid(1, 100, 200, 300, 3, 5).Subtract(cubioid(50, 60, 220, 300, 3, 5))
	resA := union(
		cubioid(61, 100, 200, 300, 3, 5),
		cubioid(50, 60, 200, 219, 3, 5),
		cubioid(1, 49, 200, 300, 3, 5),
	)
	require.True(t, a.Equal(resA))

	b := cubioid(1, 100, 200, 300, 3, 5).Subtract(cubioid(50, 1000, 0, 250, 3, 5))
	resB := union(
		cubioid(50, 100, 251, 300, 3, 5),
		cubioid(1, 49, 200, 300, 3, 5),
	)
	require.True(t, b.Equal(resB))

	c := union(
		cubioid(1, 100, 200, 300, 3, 5),
		cubioid(400, 700, 200, 300, 3, 5),
	).Subtract(cubioid(50, 1000, 0, 250, 3, 5))
	resC := union(
		cubioid(50, 100, 251, 300, 3, 5),
		cubioid(1, 49, 200, 300, 3, 5),
		cubioid(400, 700, 251, 300, 3, 5),
	)
	require.True(t, c.Equal(resC))

	d := cubioid(1, 100, 200, 300, 3, 5).Subtract(cubioid(50, 60, 220, 300, 3, 5))
	dRes := union(
		cubioid(1, 49, 200, 300, 3, 5),
		cubioid(50, 60, 200, 219, 3, 5),
		cubioid(61, 100, 200, 300, 3, 5),
	)
	require.True(t, d.Equal(dRes))
}

func TestCubioidAddHole2(t *testing.T) {
	c := union(
		cubioid(80, 100, 20, 300, 3, 5),
		cubioid(250, 400, 20, 300, 3, 5),
	).Subtract(cubioid(30, 300, 100, 102, 3, 5))
	d := union(
		cubioid(80, 100, 20, 99, 3, 5),
		cubioid(80, 100, 103, 300, 3, 5),
		cubioid(250, 300, 20, 99, 3, 5),
		cubioid(250, 300, 103, 300, 3, 5),
		cubioid(301, 400, 20, 300, 3, 5),
	)
	require.True(t, c.Equal(d))
}

func TestCubioidSubtractToEmpty(t *testing.T) {
	c := cubioid(1, 100, 200, 300, 3, 5).Subtract(cubioid(1, 100, 200, 300, 3, 5))
	require.True(t, c.IsEmpty())
}

func TestCubioidUnion1(t *testing.T) {
	c := union(
		cubioid(1, 100, 200, 300, 3, 5),
		cubioid(101, 200, 200, 300, 3, 5),
	)
	cExpected := cubioid(1, 200, 200, 300, 3, 5)
	require.True(t, cExpected.Equal(c))
}

func TestCubioidUnion2(t *testing.T) {
	c := union(
		cubioid(1, 100, 200, 300, 3, 5),
		cubioid(101, 200, 200, 300, 3, 5),
		cubioid(201, 300, 200, 300, 3, 5),
		cubioid(301, 400, 200, 300, 3, 5),
		cubioid(402, 500, 200, 300, 3, 5),
		cubioid(500, 600, 200, 700, 3, 5),
		cubioid(601, 700, 200, 700, 3, 5),
	)
	cExpected := union(
		cubioid(1, 400, 200, 300, 3, 5),
		cubioid(402, 500, 200, 300, 3, 5),
		cubioid(500, 700, 200, 700, 3, 5),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(cubioid(702, 800, 200, 700, 3, 5))
	dExpected := cExpected.Union(cubioid(702, 800, 200, 700, 3, 5))
	require.True(t, d.Equal(dExpected))
}

func TestCubioidIntersect(t *testing.T) {
	c := cubioid(5, 15, 3, 10, 3, 5).Intersect(cubioid(8, 30, 7, 20, 3, 5))
	d := cubioid(8, 15, 7, 10, 3, 5)
	require.True(t, c.Equal(d))
}

func TestCubioidUnionMerge(t *testing.T) {
	a := union(
		cubioid(5, 15, 3, 6, 3, 5),
		cubioid(5, 30, 7, 10, 3, 5),
		cubioid(8, 30, 11, 20, 3, 5),
	)
	excepted := union(
		cubioid(5, 15, 3, 10, 3, 5),
		cubioid(8, 30, 7, 20, 3, 5),
	)
	require.True(t, excepted.Equal(a))
}

func TestCubioidSubtract(t *testing.T) {
	g := cubioid(5, 15, 3, 10, 3, 5).Subtract(cubioid(8, 30, 7, 20, 3, 5))
	h := union(
		cubioid(5, 7, 3, 10, 3, 5),
		cubioid(8, 15, 3, 6, 3, 5),
	)
	require.True(t, g.Equal(h))
}

func TestCubioidIntersectEmpty(t *testing.T) {
	a := cubioid(5, 15, 3, 10, 3, 5)
	b := union(
		cubioid(1, 3, 7, 20, 3, 5),
		cubioid(20, 23, 7, 20, 3, 5),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestCubioidOr2(t *testing.T) {
	a := union(
		cubioid(1, 79, 10054, 10054, 3, 5),
		cubioid(80, 100, 10053, 10054, 3, 5),
		cubioid(101, 65535, 10054, 10054, 3, 5),
	)
	expected := union(
		cubioid(80, 100, 10053, 10053, 3, 5),
		cubioid(1, 65535, 10054, 10054, 3, 5),
	)
	require.True(t, expected.Equal(a))
}

type Tr = ds.Triple[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet]

func TestCubioidSwapDimensions(t *testing.T) {
	for _, s := range []ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet]{
		ds.NewRightTripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet](),
		ds.NewLeftTripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet](),
		ds.NewOuterTripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet](),
	} {
		require.True(t, ds.MapTripleSet(s, Tr.Swap12).Equal(s))

		require.True(t, ds.MapTripleSet(cubioid(1, 2, 3, 4, 5, 6), Tr.Swap12).Equal(cubioid(3, 4, 1, 2, 5, 6)))
		require.True(t, ds.MapTripleSet(cubioid(1, 2, 3, 4, 5, 6), Tr.Swap23).Equal(cubioid(1, 2, 5, 6, 3, 4)))
		require.True(t, ds.MapTripleSet(cubioid(1, 2, 3, 4, 5, 6), Tr.Swap13).Equal(cubioid(5, 6, 3, 4, 1, 2)))

		require.True(t, ds.MapTripleSet(union(
			cubioid(1, 3, 7, 20, 3, 5),
			cubioid(20, 23, 7, 20, 3, 5),
		), Tr.Swap12).Equal(union(
			cubioid(7, 20, 1, 3, 3, 5),
			cubioid(7, 20, 20, 23, 3, 5),
		)))
	}
}
