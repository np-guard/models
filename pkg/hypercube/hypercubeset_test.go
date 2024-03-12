// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/hypercube"
)

func TestHCBasic(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100)
	b := hypercube.FromCubeShort(1, 100)
	c := hypercube.FromCubeShort(1, 200)
	d := hypercube.FromCubeShort(1, 100, 1, 100)
	e := hypercube.FromCubeShort(1, 100, 1, 100)
	f := hypercube.FromCubeShort(1, 100, 1, 200)

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
	a := hypercube.FromCubeShort(1, 100)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestString(t *testing.T) {
	require.Equal(t, "[(1-3)]", hypercube.FromCubeShort(1, 3).String())
	require.Equal(t, "[(1-3),(2-4)]", hypercube.FromCubeShort(1, 3, 2, 4).String())
}

func TestOr(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 1, 100)
	b := hypercube.FromCubeShort(1, 90, 1, 200)
	c := a.Union(b)
	require.Equal(t, "[(1-90),(1-200)]; [(91-100),(1-100)]", c.String())
}

func addCube(o *hypercube.CanonicalSet, bounds ...int64) *hypercube.CanonicalSet {
	return o.Union(hypercube.FromCubeShort(bounds...))
}

func TestBasic(t *testing.T) {
	a := hypercube.NewCanonicalSet(1)
	a = addCube(a, 1, 2)
	a = addCube(a, 5, 6)
	a = addCube(a, 3, 4)
	b := hypercube.FromCubeShort(1, 6)
	require.True(t, a.Equal(b))
}

func TestBasic2(t *testing.T) {
	a := hypercube.NewCanonicalSet(2)
	a = addCube(a, 1, 2, 1, 5)
	a = addCube(a, 1, 2, 7, 9)
	a = addCube(a, 1, 2, 6, 7)
	b := hypercube.NewCanonicalSet(2)
	b = addCube(b, 1, 2, 1, 9)
	require.True(t, a.Equal(b))
}

func TestNew(t *testing.T) {
	a := hypercube.NewCanonicalSet(3)
	a = addCube(a, 10, 20, 10, 20, 1, 65535)
	a = addCube(a, 1, 65535, 15, 40, 1, 65535)
	a = addCube(a, 1, 65535, 100, 200, 30, 80)
	expectedStr := "[(1-9,21-65535),(100-200),(30-80)]; " +
		"[(1-9,21-65535),(15-40),(1-65535)]; " +
		"[(10-20),(10-40),(1-65535)]; " +
		"[(10-20),(100-200),(30-80)]"
	require.Equal(t, expectedStr, a.String())
}

func checkContained(t *testing.T, a, b *hypercube.CanonicalSet, expected bool) {
	t.Helper()
	contained, err := a.ContainedIn(b)
	require.Nil(t, err)
	require.Equal(t, expected, contained)
}

func TestContainedIn(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 200, 300)
	b := hypercube.FromCubeShort(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = addCube(b, 10, 200, 210, 280)
	checkContained(t, b, a, false)
}

func TestContainedIn2(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube(c, 150, 180, 20, 300)
	c = addCube(c, 200, 240, 200, 300)
	c = addCube(c, 241, 300, 200, 350)

	a := hypercube.FromCubeShort(1, 100, 200, 300)
	a = addCube(a, 150, 180, 20, 300)
	a = addCube(a, 200, 240, 200, 300)
	a = addCube(a, 242, 300, 200, 350)

	d := hypercube.FromCubeShort(210, 220, 210, 280)
	e := hypercube.FromCubeShort(210, 310, 210, 280)
	f := hypercube.FromCubeShort(210, 250, 210, 280)
	f1 := hypercube.FromCubeShort(210, 240, 210, 280)
	f2 := hypercube.FromCubeShort(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestContainedIn3(t *testing.T) {
	a := hypercube.FromCubeShort(105, 105, 54, 54)
	b := hypercube.FromCubeShort(0, 204, 0, 255)
	b = addCube(b, 205, 205, 0, 53)
	b = addCube(b, 205, 205, 55, 255)
	b = addCube(b, 206, 254, 0, 255)
	checkContained(t, a, b, true)
}

func TestContainedIn4(t *testing.T) {
	a := hypercube.FromCubeShort(105, 105, 54, 54)
	b := hypercube.FromCubeShort(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestContainedIn5(t *testing.T) {
	a := hypercube.FromCubeShort(100, 200, 54, 65, 60, 300)
	b := hypercube.FromCubeShort(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestEquals(t *testing.T) {
	a := hypercube.FromCubeShort(1, 2)
	b := hypercube.FromCubeShort(1, 2)
	require.True(t, a.Equal(b))
	c := hypercube.FromCubeShort(1, 2, 1, 5)
	d := hypercube.FromCubeShort(1, 2, 1, 5)
	require.True(t, c.Equal(d))
	c = addCube(c, 1, 2, 7, 9)
	c = addCube(c, 1, 2, 6, 7)
	c = addCube(c, 4, 8, 1, 9)
	res := hypercube.FromCubeShort(4, 8, 1, 9)
	res = addCube(res, 1, 2, 1, 9)
	require.True(t, res.Equal(c))

	a = addCube(a, 5, 6)
	a = addCube(a, 3, 4)
	res1 := hypercube.FromCubeShort(1, 6)
	require.True(t, res1.Equal(a))

	d = addCube(d, 1, 2, 1, 5)
	d = addCube(d, 5, 6, 1, 5)
	d = addCube(d, 3, 4, 1, 5)
	res2 := hypercube.FromCubeShort(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestBasicAddCube(t *testing.T) {
	a := hypercube.FromCubeShort(1, 2)
	a = addCube(a, 8, 10)
	b := a
	a = addCube(a, 1, 2)
	a = addCube(a, 6, 10)
	a = addCube(a, 1, 10)
	res := hypercube.FromCubeShort(1, 10)
	require.True(t, res.Equal(a))
	require.NotEqual(t, res, b)
}
func TestBasicAddHole(t *testing.T) {
	a := hypercube.FromCubeShort(1, 10)
	b := a.Subtract(hypercube.FromCubeShort(3, 20))
	c := a.Subtract(hypercube.FromCubeShort(0, 20))
	d := a.Subtract(hypercube.FromCubeShort(0, 5))
	e := a.Subtract(hypercube.FromCubeShort(12, 14))
	a = a.Subtract(hypercube.FromCubeShort(3, 7))
	f := hypercube.FromCubeShort(1, 2)
	f = addCube(f, 8, 10)
	require.True(t, a.Equal(f))
	require.True(t, b.Equal(hypercube.FromCubeShort(1, 2)))
	require.True(t, c.Equal(hypercube.NewCanonicalSet(1)))
	require.True(t, d.Equal(hypercube.FromCubeShort(6, 10)))
	require.True(t, e.Equal(hypercube.FromCubeShort(1, 10)))
}

func TestAddHoleBasic20(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 200, 300).Subtract(hypercube.FromCubeShort(50, 60, 220, 300))
	resA := hypercube.FromCubeShort(61, 100, 200, 300)
	resA = addCube(resA, 50, 60, 200, 219)
	resA = addCube(resA, 1, 49, 200, 300)
	require.True(t, a.Equal(resA), fmt.Sprintf("%v != %v", a, resA))
}

func TestAddHoleBasic21(t *testing.T) {
	b := hypercube.FromCubeShort(1, 100, 200, 300).Subtract(hypercube.FromCubeShort(50, 1000, 0, 250))
	resB := hypercube.FromCubeShort(50, 100, 251, 300)
	resB = addCube(resB, 1, 49, 200, 300)
	require.True(t, b.Equal(resB), fmt.Sprintf("%v != %v", b, resB))
}

func TestAddHoleBasic22(t *testing.T) {
	a := hypercube.FromCubeShort(1, 2, 1, 2)
	require.Equal(t, "[(1),(2)]; [(2),(1-2)]", a.Subtract(hypercube.FromCubeShort(1, 1, 1, 1)).String())
	require.Equal(t, "[(1),(1)]; [(2),(1-2)]", a.Subtract(hypercube.FromCubeShort(1, 1, 2, 2)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(2)]", a.Subtract(hypercube.FromCubeShort(2, 2, 1, 1)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(1)]", a.Subtract(hypercube.FromCubeShort(2, 2, 2, 2)).String())
}

func TestAddHoleBasic23(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 200, 300)
	a = addCube(a, 400, 700, 200, 300)
	require.Equal(t, "[(1-100,400-700),(200-300)]", a.String())
	a = a.Subtract(hypercube.FromCubeShort(50, 1000, 0, 250))
	require.Equal(t, "[(1-49),(200-300)]; [(50-100,400-700),(251-300)]", a.String())
}

func TestAddHole(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = c.Subtract(hypercube.FromCubeShort(50, 60, 220, 300))
	d := hypercube.FromCubeShort(1, 49, 200, 300)
	d = addCube(d, 50, 60, 200, 219)
	d = addCube(d, 61, 100, 200, 300)
	require.True(t, c.Equal(d))
}

func TestAddHole2(t *testing.T) {
	c := hypercube.FromCubeShort(80, 100, 20, 300)
	c = addCube(c, 250, 400, 20, 300)
	c = c.Subtract(hypercube.FromCubeShort(30, 300, 100, 102))
	d := hypercube.FromCubeShort(80, 100, 20, 99)
	d = addCube(d, 80, 100, 103, 300)
	d = addCube(d, 250, 300, 20, 99)
	d = addCube(d, 250, 300, 103, 300)
	d = addCube(d, 301, 400, 20, 300)
	require.True(t, c.Equal(d))
}

func TestAddHole3(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = c.Subtract(hypercube.FromCubeShort(1, 100, 200, 300))
	require.Equal(t, c, hypercube.NewCanonicalSet(2))
}

func TestIntervalsUnion(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube(c, 101, 200, 200, 300)
	d := hypercube.FromCubeShort(1, 200, 200, 300)
	require.True(t, c.Equal(d))
	if c.String() != d.String() {
		t.FailNow()
	}
}

func TestIntervalsUnion2(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube(c, 101, 200, 200, 300)
	c = addCube(c, 201, 300, 200, 300)
	c = addCube(c, 301, 400, 200, 300)
	c = addCube(c, 402, 500, 200, 300)
	c = addCube(c, 500, 600, 200, 700)
	c = addCube(c, 601, 700, 200, 700)

	d := c.Copy()
	d = addCube(d, 702, 800, 200, 700)

	cExpected := hypercube.FromCubeShort(1, 400, 200, 300)
	cExpected = addCube(cExpected, 402, 500, 200, 300)
	cExpected = addCube(cExpected, 500, 700, 200, 700)
	dExpected := cExpected.Copy()
	dExpected = addCube(dExpected, 702, 800, 200, 700)
	require.True(t, c.Equal(cExpected))
	require.True(t, d.Equal(dExpected))
}

func TestAndSubOr(t *testing.T) {
	a := hypercube.FromCubeShort(5, 15, 3, 10)
	b := hypercube.FromCubeShort(8, 30, 7, 20)

	c := a.Intersect(b)
	d := hypercube.FromCubeShort(8, 15, 7, 10)
	require.True(t, c.Equal(d))

	f := a.Union(b)
	e := hypercube.FromCubeShort(5, 15, 3, 6)
	e = addCube(e, 5, 30, 7, 10)
	e = addCube(e, 8, 30, 11, 20)
	require.True(t, e.Equal(f))

	g := a.Subtract(b)
	h := hypercube.FromCubeShort(5, 7, 3, 10)
	h = addCube(h, 8, 15, 3, 6)
	require.True(t, g.Equal(h))
}

func TestAnd2(t *testing.T) {
	a := hypercube.FromCubeShort(5, 15, 3, 10)
	b := hypercube.FromCubeShort(1, 3, 7, 20)
	b = addCube(b, 20, 23, 7, 20)
	c := a.Intersect(b)
	require.Equal(t, c, hypercube.NewCanonicalSet(2))
}

func TestOr2(t *testing.T) {
	a := hypercube.FromCubeShort(80, 100, 10053, 10053)
	b := hypercube.FromCubeShort(1, 65535, 10054, 10054)
	a = a.Union(b)
	expected := hypercube.FromCubeShort(1, 79, 10054, 10054)
	expected = addCube(expected, 80, 100, 10053, 10054)
	expected = addCube(expected, 101, 65535, 10054, 10054)
	require.True(t, a.Equal(expected))
}
