// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/hypercube"
)

func union(set *hypercube.CanonicalSet, sets ...*hypercube.CanonicalSet) *hypercube.CanonicalSet {
	for _, c := range sets {
		set = set.Union(c)
	}
	return set.Copy()
}

func TestHCBasic(t *testing.T) {
	a := hypercube.Cube(1, 100)
	b := hypercube.Cube(1, 100)
	c := hypercube.Cube(1, 200)
	d := hypercube.Cube(1, 100, 1, 100)
	e := hypercube.Cube(1, 100, 1, 100)
	f := hypercube.Cube(1, 100, 1, 200)

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
	a := hypercube.Cube(1, 100)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestString(t *testing.T) {
	require.Equal(t, "[(1-3)]", hypercube.Cube(1, 3).String())
	require.Equal(t, "[(1-3),(2-4)]", hypercube.Cube(1, 3, 2, 4).String())
}

func TestOr(t *testing.T) {
	a := hypercube.Cube(1, 100, 1, 100)
	b := hypercube.Cube(1, 90, 1, 200)
	c := a.Union(b)
	require.Equal(t, "[(1-90),(1-200)]; [(91-100),(1-100)]", c.String())
}

func TestBasic1(t *testing.T) {
	a := union(
		hypercube.Cube(1, 2),
		hypercube.Cube(5, 6),
		hypercube.Cube(3, 4),
	)
	b := hypercube.Cube(1, 6)
	require.True(t, a.Equal(b))
}

func TestBasic2(t *testing.T) {
	a := union(
		hypercube.Cube(1, 2, 1, 5),
		hypercube.Cube(1, 2, 7, 9),
		hypercube.Cube(1, 2, 6, 7),
	)
	b := hypercube.Cube(1, 2, 1, 9)
	require.True(t, a.Equal(b))
}

func TestNew(t *testing.T) {
	a := union(
		hypercube.Cube(10, 20, 10, 20, 1, 65535),
		hypercube.Cube(1, 65535, 15, 40, 1, 65535),
		hypercube.Cube(1, 65535, 100, 200, 30, 80),
	)
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
	a := hypercube.Cube(1, 100, 200, 300)
	b := hypercube.Cube(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = b.Union(hypercube.Cube(10, 200, 210, 280))
	checkContained(t, b, a, false)
}

func TestContainedIn1(t *testing.T) {
	checkContained(t, hypercube.Cube(1, 3), hypercube.Cube(2, 4), false)
	checkContained(t, hypercube.Cube(2, 4), hypercube.Cube(1, 3), false)
	checkContained(t, hypercube.Cube(1, 3), hypercube.Cube(1, 4), true)
	checkContained(t, hypercube.Cube(1, 4), hypercube.Cube(1, 3), false)
}

func TestContainedIn2(t *testing.T) {
	c := union(
		hypercube.Cube(1, 100, 200, 300),
		hypercube.Cube(150, 180, 20, 300),
		hypercube.Cube(200, 240, 200, 300),
		hypercube.Cube(241, 300, 200, 350),
	)

	a := union(
		hypercube.Cube(1, 100, 200, 300),
		hypercube.Cube(150, 180, 20, 300),
		hypercube.Cube(200, 240, 200, 300),
		hypercube.Cube(242, 300, 200, 350),
	)
	d := hypercube.Cube(210, 220, 210, 280)
	e := hypercube.Cube(210, 310, 210, 280)
	f := hypercube.Cube(210, 250, 210, 280)
	f1 := hypercube.Cube(210, 240, 210, 280)
	f2 := hypercube.Cube(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestContainedIn3(t *testing.T) {
	a := hypercube.Cube(105, 105, 54, 54)
	b := union(
		hypercube.Cube(0, 204, 0, 255),
		hypercube.Cube(205, 205, 0, 53),
		hypercube.Cube(205, 205, 55, 255),
		hypercube.Cube(206, 254, 0, 255),
	)
	checkContained(t, a, b, true)
}

func TestContainedIn4(t *testing.T) {
	a := hypercube.Cube(105, 105, 54, 54)
	b := hypercube.Cube(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestContainedIn5(t *testing.T) {
	a := hypercube.Cube(100, 200, 54, 65, 60, 300)
	b := hypercube.Cube(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestEqual1(t *testing.T) {
	a := hypercube.Cube(1, 2)
	b := hypercube.Cube(1, 2)
	require.True(t, a.Equal(b))

	c := hypercube.Cube(1, 2, 1, 5)
	d := hypercube.Cube(1, 2, 1, 5)
	require.True(t, c.Equal(d))
}

func TestEqual2(t *testing.T) {
	c := union(
		hypercube.Cube(1, 2, 1, 5),
		hypercube.Cube(1, 2, 7, 9),
		hypercube.Cube(1, 2, 6, 7),
		hypercube.Cube(4, 8, 1, 9),
	)
	res := union(
		hypercube.Cube(4, 8, 1, 9),
		hypercube.Cube(1, 2, 1, 9),
	)
	require.True(t, res.Equal(c))

	d := union(
		hypercube.Cube(1, 2, 1, 5),
		hypercube.Cube(5, 6, 1, 5),
		hypercube.Cube(3, 4, 1, 5),
	)
	res2 := hypercube.Cube(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestBasicAddCube(t *testing.T) {
	a := union(
		hypercube.Cube(1, 2),
		hypercube.Cube(8, 10),
	)
	b := union(
		a,
		hypercube.Cube(1, 2),
		hypercube.Cube(6, 10),
		hypercube.Cube(1, 10),
	)
	res := hypercube.Cube(1, 10)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestFourHoles(t *testing.T) {
	a := hypercube.Cube(1, 2, 1, 2)
	require.Equal(t, "[(1),(2)]; [(2),(1-2)]", a.Subtract(hypercube.Cube(1, 1, 1, 1)).String())
	require.Equal(t, "[(1),(1)]; [(2),(1-2)]", a.Subtract(hypercube.Cube(1, 1, 2, 2)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(2)]", a.Subtract(hypercube.Cube(2, 2, 1, 1)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(1)]", a.Subtract(hypercube.Cube(2, 2, 2, 2)).String())
}

func TestBasicSubtract1(t *testing.T) {
	a := hypercube.Cube(1, 10)
	require.True(t, a.Subtract(hypercube.Cube(3, 7)).Equal(union(hypercube.Cube(1, 2), hypercube.Cube(8, 10))))
	require.True(t, a.Subtract(hypercube.Cube(3, 20)).Equal(hypercube.Cube(1, 2)))
	require.True(t, a.Subtract(hypercube.Cube(0, 20)).IsEmpty())
	require.True(t, a.Subtract(hypercube.Cube(0, 5)).Equal(hypercube.Cube(6, 10)))
	require.True(t, a.Subtract(hypercube.Cube(12, 14)).Equal(hypercube.Cube(1, 10)))
}

func TestBasicSubtract2(t *testing.T) {
	a := hypercube.Cube(1, 100, 200, 300).Subtract(hypercube.Cube(50, 60, 220, 300))
	resA := union(
		hypercube.Cube(61, 100, 200, 300),
		hypercube.Cube(50, 60, 200, 219),
		hypercube.Cube(1, 49, 200, 300),
	)
	require.True(t, a.Equal(resA))

	b := hypercube.Cube(1, 100, 200, 300).Subtract(hypercube.Cube(50, 1000, 0, 250))
	resB := union(
		hypercube.Cube(50, 100, 251, 300),
		hypercube.Cube(1, 49, 200, 300),
	)
	require.True(t, b.Equal(resB))

	c := union(
		hypercube.Cube(1, 100, 200, 300),
		hypercube.Cube(400, 700, 200, 300),
	).Subtract(hypercube.Cube(50, 1000, 0, 250))
	resC := union(
		hypercube.Cube(50, 100, 251, 300),
		hypercube.Cube(1, 49, 200, 300),
		hypercube.Cube(400, 700, 251, 300),
	)
	require.True(t, c.Equal(resC))

	d := hypercube.Cube(1, 100, 200, 300).Subtract(hypercube.Cube(50, 60, 220, 300))
	dRes := union(
		hypercube.Cube(1, 49, 200, 300),
		hypercube.Cube(50, 60, 200, 219),
		hypercube.Cube(61, 100, 200, 300),
	)
	require.True(t, d.Equal(dRes))
}

func TestAddHole2(t *testing.T) {
	c := union(
		hypercube.Cube(80, 100, 20, 300),
		hypercube.Cube(250, 400, 20, 300),
	).Subtract(hypercube.Cube(30, 300, 100, 102))
	d := union(
		hypercube.Cube(80, 100, 20, 99),
		hypercube.Cube(80, 100, 103, 300),
		hypercube.Cube(250, 300, 20, 99),
		hypercube.Cube(250, 300, 103, 300),
		hypercube.Cube(301, 400, 20, 300),
	)
	require.True(t, c.Equal(d))
}

func TestSubtractToEmpty(t *testing.T) {
	c := hypercube.Cube(1, 100, 200, 300).Subtract(hypercube.Cube(1, 100, 200, 300))
	require.True(t, c.IsEmpty())
}

func TestUnion1(t *testing.T) {
	c := union(
		hypercube.Cube(1, 100, 200, 300),
		hypercube.Cube(101, 200, 200, 300),
	)
	cExpected := hypercube.Cube(1, 200, 200, 300)
	require.True(t, cExpected.Equal(c))
}

func TestUnion2(t *testing.T) {
	c := union(
		hypercube.Cube(1, 100, 200, 300),
		hypercube.Cube(101, 200, 200, 300),
		hypercube.Cube(201, 300, 200, 300),
		hypercube.Cube(301, 400, 200, 300),
		hypercube.Cube(402, 500, 200, 300),
		hypercube.Cube(500, 600, 200, 700),
		hypercube.Cube(601, 700, 200, 700),
	)
	cExpected := union(
		hypercube.Cube(1, 400, 200, 300),
		hypercube.Cube(402, 500, 200, 300),
		hypercube.Cube(500, 700, 200, 700),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(hypercube.Cube(702, 800, 200, 700))
	dExpected := cExpected.Union(hypercube.Cube(702, 800, 200, 700))
	require.True(t, d.Equal(dExpected))
}

func TestIntersect(t *testing.T) {
	c := hypercube.Cube(5, 15, 3, 10).Intersect(hypercube.Cube(8, 30, 7, 20))
	d := hypercube.Cube(8, 15, 7, 10)
	require.True(t, c.Equal(d))
}

func TestUnionMerge(t *testing.T) {
	a := union(
		hypercube.Cube(5, 15, 3, 6),
		hypercube.Cube(5, 30, 7, 10),
		hypercube.Cube(8, 30, 11, 20),
	)
	excepted := union(
		hypercube.Cube(5, 15, 3, 10),
		hypercube.Cube(8, 30, 7, 20),
	)
	require.True(t, excepted.Equal(a))
}

func TestSubtract(t *testing.T) {
	g := hypercube.Cube(5, 15, 3, 10).Subtract(hypercube.Cube(8, 30, 7, 20))
	h := union(
		hypercube.Cube(5, 7, 3, 10),
		hypercube.Cube(8, 15, 3, 6),
	)
	require.True(t, g.Equal(h))
}

func TestIntersectEmpty(t *testing.T) {
	a := hypercube.Cube(5, 15, 3, 10)
	b := union(
		hypercube.Cube(1, 3, 7, 20),
		hypercube.Cube(20, 23, 7, 20),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestOr2(t *testing.T) {
	a := union(
		hypercube.Cube(1, 79, 10054, 10054),
		hypercube.Cube(80, 100, 10053, 10054),
		hypercube.Cube(101, 65535, 10054, 10054),
	)
	expected := union(
		hypercube.Cube(80, 100, 10053, 10053),
		hypercube.Cube(1, 65535, 10054, 10054),
	)
	require.True(t, expected.Equal(a))
}

// Assisted by WCA for GP
// Latest GenAI contribution: granite-20B-code-instruct-v2 model
// TestSwapDimensions tests the SwapDimensions method of the CanonicalSet type.
func TestSwapDimensions(t *testing.T) {
	tests := []struct {
		name     string
		c        *hypercube.CanonicalSet
		dim1     int
		dim2     int
		expected *hypercube.CanonicalSet
	}{
		{
			name:     "empty set",
			c:        hypercube.NewCanonicalSet(2),
			dim1:     0,
			dim2:     1,
			expected: hypercube.NewCanonicalSet(2),
		},
		{
			name:     "0,0 of 1",
			c:        hypercube.Cube(1, 2),
			dim1:     0,
			dim2:     0,
			expected: hypercube.Cube(1, 2),
		},
		{
			name:     "0,1 of 2",
			c:        hypercube.Cube(1, 2, 3, 4),
			dim1:     0,
			dim2:     1,
			expected: hypercube.Cube(3, 4, 1, 2),
		},
		{
			name:     "0,1 of 2, no-op",
			c:        hypercube.Cube(1, 2, 1, 2),
			dim1:     0,
			dim2:     1,
			expected: hypercube.Cube(1, 2, 1, 2),
		},
		{
			name:     "0,1 of 3",
			c:        hypercube.Cube(1, 2, 3, 4, 5, 6),
			dim1:     0,
			dim2:     1,
			expected: hypercube.Cube(3, 4, 1, 2, 5, 6),
		},
		{
			name:     "1,2 of 3",
			c:        hypercube.Cube(1, 2, 3, 4, 5, 6),
			dim1:     1,
			dim2:     2,
			expected: hypercube.Cube(1, 2, 5, 6, 3, 4),
		},
		{
			name:     "0,2 of 3",
			c:        hypercube.Cube(1, 2, 3, 4, 5, 6),
			dim1:     0,
			dim2:     2,
			expected: hypercube.Cube(5, 6, 3, 4, 1, 2),
		},
		{
			name: "0,1 of 2, non-cube",
			c: union(
				hypercube.Cube(1, 3, 7, 20),
				hypercube.Cube(20, 23, 7, 20),
			),
			dim1: 0,
			dim2: 1,
			expected: union(
				hypercube.Cube(7, 20, 1, 3),
				hypercube.Cube(7, 20, 20, 23),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.c.SwapDimensions(tt.dim1, tt.dim2)
			require.True(t, tt.expected != actual)
			require.True(t, tt.expected.Equal(actual))
		})
	}
	require.Panics(t, func() { hypercube.Cube(1, 2).SwapDimensions(0, 1) })
	require.Panics(t, func() { hypercube.Cube(1, 2).SwapDimensions(-1, 0) })
}
