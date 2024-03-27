// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
)

// square returns a new ds.Product created from a single input square
// the input square is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func square(s1, e1, s2, e2 int64) *ds.Product[*interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianPair(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet())
}

func TestCubeBasic(t *testing.T) {
	d := square(1, 100, 1, 100)
	e := square(1, 100, 1, 100)
	f := square(1, 100, 1, 200)

	require.True(t, d.Equal(e))
	require.True(t, e.Equal(d))

	require.False(t, d.Equal(f))
	require.False(t, f.Equal(d))

	a := union(
		square(1, 2, 1, 5),
		square(1, 2, 7, 9),
		square(1, 2, 6, 7),
	)
	b := square(1, 2, 1, 9)
	require.True(t, a.Equal(b))
}

func TestCubeCopy(t *testing.T) {
	a := square(1, 100, 2, 200)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestCubeIsSubset(t *testing.T) {
	a := square(1, 100, 200, 300)
	b := square(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = b.Union(square(10, 200, 210, 280))
	checkContained(t, b, a, false)

	c := union(
		square(1, 100, 200, 300),
		square(150, 180, 20, 300),
		square(200, 240, 200, 300),
		square(241, 300, 200, 350),
	)

	d := square(210, 220, 210, 280)
	e := square(210, 310, 210, 280)
	f := square(210, 250, 210, 280)
	f1 := square(210, 240, 210, 280)
	f2 := square(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)

	a = union(
		square(1, 100, 200, 300),
		square(150, 180, 20, 300),
		square(200, 240, 200, 300),
		square(242, 300, 200, 350),
	)
	checkContained(t, f, a, false)
}

func TestCubeIsSubset3(t *testing.T) {
	a := square(105, 105, 54, 54)
	b := union(
		square(0, 204, 0, 255),
		square(205, 205, 0, 53),
		square(205, 205, 55, 255),
		square(206, 254, 0, 255),
	)
	checkContained(t, a, b, true)
}

func TestCubeIsSubset4(t *testing.T) {
	a := square(105, 105, 54, 54)
	b := square(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestCubeEqual1(t *testing.T) {
	c := square(1, 2, 1, 5)
	d := square(1, 2, 1, 5)
	require.True(t, c.Equal(d))
}

func TestCubeEqual2(t *testing.T) {
	c := union(
		square(1, 2, 1, 5),
		square(1, 2, 7, 9),
		square(1, 2, 6, 7),
		square(4, 8, 1, 9),
	)
	res := union(
		square(4, 8, 1, 9),
		square(1, 2, 1, 9),
	)
	require.True(t, res.Equal(c))

	d := union(
		square(1, 2, 1, 5),
		square(5, 6, 1, 5),
		square(3, 4, 1, 5),
	)
	res2 := square(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestCubeBasicAddCube(t *testing.T) {
	a := union(
		square(1, 2, 3, 4),
		square(8, 10, 3, 4),
	)
	b := union(
		a,
		square(1, 2, 3, 4),
		square(6, 10, 3, 4),
		square(1, 10, 3, 4),
	)
	res := square(1, 10, 3, 4)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestCubeBasicSubtract(t *testing.T) {
	a := square(1, 100, 200, 300).Subtract(square(50, 60, 220, 300))
	resA := union(
		square(61, 100, 200, 300),
		square(50, 60, 200, 219),
		square(1, 49, 200, 300),
	)
	require.True(t, a.Equal(resA))

	b := square(1, 100, 200, 300).Subtract(square(50, 1000, 0, 250))
	resB := union(
		square(50, 100, 251, 300),
		square(1, 49, 200, 300),
	)
	require.True(t, b.Equal(resB))

	c := union(
		square(1, 100, 200, 300),
		square(400, 700, 200, 300),
	).Subtract(square(50, 1000, 0, 250))
	resC := union(
		square(50, 100, 251, 300),
		square(1, 49, 200, 300),
		square(400, 700, 251, 300),
	)
	require.True(t, c.Equal(resC))

	d := square(1, 100, 200, 300).Subtract(square(50, 60, 220, 300))
	dRes := union(
		square(1, 49, 200, 300),
		square(50, 60, 200, 219),
		square(61, 100, 200, 300),
	)
	require.True(t, d.Equal(dRes))
}

func TestCubeAddHole(t *testing.T) {
	c := union(
		square(80, 100, 20, 300),
		square(250, 400, 20, 300),
	).Subtract(square(30, 300, 100, 102))
	d := union(
		square(80, 100, 20, 99),
		square(80, 100, 103, 300),
		square(250, 300, 20, 99),
		square(250, 300, 103, 300),
		square(301, 400, 20, 300),
	)
	require.True(t, c.Equal(d))
}

func TestCubeSubtractToEmpty(t *testing.T) {
	c := square(1, 100, 200, 300).Subtract(square(1, 100, 200, 300))
	require.True(t, c.IsEmpty())
}

func TestCubeUnion1(t *testing.T) {
	c := union(
		square(1, 100, 200, 300),
		square(101, 200, 200, 300),
	)
	cExpected := square(1, 200, 200, 300)
	require.True(t, cExpected.Equal(c))
}

func TestCubeUnion2(t *testing.T) {
	c := union(
		square(1, 100, 200, 300),
		square(101, 200, 200, 300),
		square(201, 300, 200, 300),
		square(301, 400, 200, 300),
		square(402, 500, 200, 300),
		square(500, 600, 200, 700),
		square(601, 700, 200, 700),
	)
	cExpected := union(
		square(1, 400, 200, 300),
		square(402, 500, 200, 300),
		square(500, 700, 200, 700),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(square(702, 800, 200, 700))
	dExpected := cExpected.Union(square(702, 800, 200, 700))
	require.True(t, d.Equal(dExpected))
}

func TestCubeIntersect(t *testing.T) {
	c := square(5, 15, 3, 10).Intersect(square(8, 30, 7, 20))
	d := square(8, 15, 7, 10)
	require.True(t, c.Equal(d))
}

func TestCubeUnionMerge(t *testing.T) {
	a := union(
		square(5, 15, 3, 6),
		square(5, 30, 7, 10),
		square(8, 30, 11, 20),
	)
	excepted := union(
		square(5, 15, 3, 10),
		square(8, 30, 7, 20),
	)
	require.True(t, excepted.Equal(a))
}

func TestCubeSubtract(t *testing.T) {
	g := square(5, 15, 3, 10).Subtract(square(8, 30, 7, 20))
	h := union(
		square(5, 7, 3, 10),
		square(8, 15, 3, 6),
	)
	require.True(t, g.Equal(h))
}

func TestCubeIntersectEmpty(t *testing.T) {
	a := square(5, 15, 3, 10)
	b := union(
		square(1, 3, 7, 20),
		square(20, 23, 7, 20),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestCubeUnion3(t *testing.T) {
	a := union(
		square(1, 79, 10054, 10054),
		square(80, 100, 10053, 10054),
		square(101, 65535, 10054, 10054),
	)
	expected := union(
		square(80, 100, 10053, 10053),
		square(1, 65535, 10054, 10054),
	)
	require.True(t, expected.Equal(a))
}

func TestCubeSwapDimensions(t *testing.T) {
	require.True(t, square(1, 2, 3, 4).Swap().Equal(square(3, 4, 1, 2)))
	require.True(t, square(1, 2, 1, 2).Swap().Equal(square(1, 2, 1, 2)))

	require.True(t, union(
		square(1, 3, 7, 20),
		square(20, 23, 7, 20),
	).Swap().Equal(union(
		square(7, 20, 1, 3),
		square(7, 20, 20, 23),
	)))
}
