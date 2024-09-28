/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
)

func union[S ds.Set[S]](set S, sets ...S) S {
	for _, c := range sets {
		set = set.Union(c)
	}
	return set.Copy()
}

func checkContained[S ds.Set[S]](t *testing.T, a, b S, expected bool) {
	t.Helper()
	contained := a.IsSubset(b)
	require.Equal(t, expected, contained)
}

// rectangle returns a new ds.Product created from a single input rectangle
// the input rectangle is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func rectangle(s1, e1, s2, e2 int64) ds.Product[*interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianPairLeft(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet())
}

func TestRectangleProductEmpty(t *testing.T) {
	rightEmpty := ds.CartesianPairLeft(interval.New(1, 2).ToSet(), interval.NewCanonicalSet())
	require.True(t, rightEmpty.IsEmpty())
	leftEmpty := ds.CartesianPairLeft(interval.NewCanonicalSet(), interval.New(1, 2).ToSet())
	require.True(t, leftEmpty.IsEmpty())
}

func TestRectangleEqual(t *testing.T) {

	// d is of type: ds.Product[*interval.CanonicalSet, *interval.CanonicalSet] , string: {(1-100 x 1-100)}
	d := rectangle(1, 100, 1, 100)

	// e is of type: ds.Product[*interval.CanonicalSet, *interval.CanonicalSet] , string: {(1-100 x 1-100)}
	e := rectangle(1, 100, 1, 100)

	// f is of type: ds.Product[*interval.CanonicalSet, *interval.CanonicalSet] , string: {(1-100 x 1-200)}
	f := rectangle(1, 100, 1, 200)
	fmt.Println(d.String())
	fmt.Println(e.String())
	fmt.Println(f.String())

	require.True(t, d.Equal(e))
	require.True(t, e.Equal(d))

	require.False(t, d.Equal(f))
	require.False(t, f.Equal(d))

	a := union(
		rectangle(1, 2, 1, 5),
		rectangle(1, 2, 7, 9),
		rectangle(1, 2, 6, 7),
	)
	b := rectangle(1, 2, 1, 9)
	require.True(t, a.Equal(b))

	c := rectangle(1, 50, 1, 101)
	fmt.Println(c)           // {(1-50 x 1-101)}
	fmt.Println(union(c, d)) // {(1-50 x 1-101) | (51-100 x 1-100)}
	fmt.Println(union(d, c)) // {(1-50 x 1-101) | (51-100 x 1-100)}
	g := rectangle(1, 1000, 101, 101)
	h := union(d, c).Subtract(g) // {(1-100 x 1-100)}
	fmt.Println(h)
	require.True(t, h.Equal(d))
	require.True(t, e.Equal(h))

	y := rectangle(3, 50, 1, 101)
	fmt.Println(y)           // {(3-50 x 1-101)}
	fmt.Println(union(y, d)) // {(3-50 x 1-101) | (1-2,51-100 x 1-100)}
	fmt.Println(union(d, y)) // {(3-50 x 1-101) | (1-2,51-100 x 1-100)}
	fmt.Println("done")

}

// This example demonstrates the canonical representation, and the "left" impact of "product_left":
// the left dimension values determine the partitions
// all three representations are equivalent, converging to the "left" canonical representation
func TestCacnonicalRep(t *testing.T) {
	a1 := rectangle(1, 7, 1, 3)
	a2 := rectangle(8, 9, 1, 5)
	resA := union(a1, a2)
	fmt.Println(resA) // {(8-9 x 1-5) | (1-7 x 1-3)}

	b1 := rectangle(1, 7, 1, 3)
	b2 := rectangle(8, 9, 1, 3)
	b3 := rectangle(8, 9, 4, 5)
	resB := union(b1, b2, b3)
	fmt.Println(resB) // {(8-9 x 1-5) | (1-7 x 1-3)}

	c1 := rectangle(1, 9, 1, 3)
	c2 := rectangle(8, 9, 4, 5)
	resC := union(c1, c2)
	fmt.Println(resC) // {(8-9 x 1-5) | (1-7 x 1-3)}

	require.True(t, resA.Equal(resB))
	require.True(t, resA.Equal(resC))
	require.True(t, resB.Equal(resC))
}

func TestRectangleCopy(t *testing.T) {
	a := rectangle(1, 100, 2, 200)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestRectangleIsSubset(t *testing.T) {
	a := rectangle(1, 100, 200, 300)
	b := rectangle(10, 80, 210, 280) // b is subset of a
	checkContained(t, b, a, true)
	// print the delta
	fmt.Println(a.Subtract(b)) // {(1-9,81-100 x 200-300) | (10-80 x 200-209,281-300)}
	require.True(t, b.Subtract(a).IsEmpty())

	b1 := b.Union(rectangle(10, 200, 210, 280)) // b1 is not a subset of a (first dimension exceeds a's range)
	checkContained(t, b1, a, false)
	checkContained(t, a, b1, false) // a is not a subset of b1
	// print the delta
	fmt.Println(b1.Subtract(a)) // {(101-200 x 210-280)}
	fmt.Println(a.Subtract(b1)) // {(10-100 x 200-209,281-300) | (1-9 x 200-300)}

	b2 := b.Union(rectangle(10, 100, 210, 310)) // b2 is not a subset of a (second dimension exceeds a's range)
	checkContained(t, b2, a, false)
	checkContained(t, a, b2, false) // b2 is not a subset of a
	// print the detla
	fmt.Println(b2.Subtract(a)) // {(10-100 x 301-310)}
	fmt.Println(a.Subtract(b2)) // {(10-100 x 200-209) | (1-9 x 200-300)}

	b3 := b.Union(rectangle(99, 99, 201, 300)) // b3 is a subset of a
	checkContained(t, b3, a, true)
	// print the delta (a minus b3):
	// {(99 x 200) | (10-80 x 200-209,281-300) | (1-9,81-98,100 x 200-300)}
	fmt.Println(a.Subtract(b3).String())

	c := union(
		rectangle(1, 100, 200, 300),
		rectangle(150, 180, 20, 300),
		rectangle(200, 240, 200, 300),
		rectangle(241, 300, 200, 350),
	)

	d := rectangle(210, 220, 210, 280)
	e := rectangle(210, 310, 210, 280)
	f := rectangle(210, 250, 210, 280)
	f1 := rectangle(210, 240, 210, 280)
	f2 := rectangle(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)

	a = union(
		rectangle(1, 100, 200, 300),
		rectangle(150, 180, 20, 300),
		rectangle(200, 240, 200, 300),
		rectangle(242, 300, 200, 350),
	)
	checkContained(t, f, a, false)
}

func TestRectangleIsSubset3(t *testing.T) {
	a := rectangle(105, 105, 54, 54)
	b := union(
		rectangle(0, 204, 0, 255),
		rectangle(205, 205, 0, 53),
		rectangle(205, 205, 55, 255),
		rectangle(206, 254, 0, 255),
	)
	checkContained(t, a, b, true)
}

func TestRectangleIsSubset4(t *testing.T) {
	a := rectangle(105, 105, 54, 54)
	b := rectangle(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestRectangleEqual1(t *testing.T) {
	c := rectangle(1, 2, 1, 5)
	d := rectangle(1, 2, 1, 5)
	require.True(t, c.Equal(d))
}

func TestRectangleEqual2(t *testing.T) {
	c := union(
		rectangle(1, 2, 1, 5),
		rectangle(1, 2, 7, 9),
		rectangle(1, 2, 6, 7),
		rectangle(4, 8, 1, 9),
	)
	res := union(
		rectangle(4, 8, 1, 9),
		rectangle(1, 2, 1, 9),
	)
	require.True(t, res.Equal(c))

	d := union(
		rectangle(1, 2, 1, 5),
		rectangle(5, 6, 1, 5),
		rectangle(3, 4, 1, 5),
	)
	res2 := rectangle(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestRectangleBasicAddCube(t *testing.T) {
	a := union(
		rectangle(1, 2, 3, 4),
		rectangle(8, 10, 3, 4),
	)
	b := union(
		a,
		rectangle(1, 2, 3, 4),
		rectangle(6, 10, 3, 4),
		rectangle(1, 10, 3, 4),
	)
	res := rectangle(1, 10, 3, 4)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestRectangleBasicSubtract(t *testing.T) {
	a := rectangle(1, 100, 200, 300).Subtract(rectangle(50, 60, 220, 300))
	resA := union(
		rectangle(61, 100, 200, 300),
		rectangle(50, 60, 200, 219),
		rectangle(1, 49, 200, 300),
	)
	require.True(t, a.Equal(resA))

	b := rectangle(1, 100, 200, 300).Subtract(rectangle(50, 1000, 0, 250))
	resB := union(
		rectangle(50, 100, 251, 300),
		rectangle(1, 49, 200, 300),
	)
	require.True(t, b.Equal(resB))

	c := union(
		rectangle(1, 100, 200, 300),
		rectangle(400, 700, 200, 300),
	).Subtract(rectangle(50, 1000, 0, 250))
	resC := union(
		rectangle(50, 100, 251, 300),
		rectangle(1, 49, 200, 300),
		rectangle(400, 700, 251, 300),
	)
	require.True(t, c.Equal(resC))

	d := rectangle(1, 100, 200, 300).Subtract(rectangle(50, 60, 220, 300))
	dRes := union(
		rectangle(1, 49, 200, 300),
		rectangle(50, 60, 200, 219),
		rectangle(61, 100, 200, 300),
	)
	require.True(t, d.Equal(dRes))
}

func TestRectangleAddHole(t *testing.T) {
	c := union(
		rectangle(80, 100, 20, 300),
		rectangle(250, 400, 20, 300),
	).Subtract(rectangle(30, 300, 100, 102))
	d := union(
		rectangle(80, 100, 20, 99),
		rectangle(80, 100, 103, 300),
		rectangle(250, 300, 20, 99),
		rectangle(250, 300, 103, 300),
		rectangle(301, 400, 20, 300),
	)
	require.True(t, c.Equal(d))
}

func TestRectangleSubtractToEmpty(t *testing.T) {
	c := rectangle(1, 100, 200, 300).Subtract(rectangle(1, 100, 200, 300))
	require.True(t, c.IsEmpty())
}

func TestRectangleUnion1(t *testing.T) {
	c := union(
		rectangle(1, 100, 200, 300),
		rectangle(101, 200, 200, 300),
	)
	cExpected := rectangle(1, 200, 200, 300)
	require.True(t, cExpected.Equal(c))
}

func TestRectangleUnion2(t *testing.T) {
	c := union(
		rectangle(1, 100, 200, 300),
		rectangle(101, 200, 200, 300),
		rectangle(201, 300, 200, 300),
		rectangle(301, 400, 200, 300),
		rectangle(402, 500, 200, 300),
		rectangle(500, 600, 200, 700),
		rectangle(601, 700, 200, 700),
	)
	cExpected := union(
		rectangle(1, 400, 200, 300),
		rectangle(402, 500, 200, 300),
		rectangle(500, 700, 200, 700),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(rectangle(702, 800, 200, 700))
	dExpected := cExpected.Union(rectangle(702, 800, 200, 700))
	require.True(t, d.Equal(dExpected))
}

func TestRectangleIntersect(t *testing.T) {
	c := rectangle(5, 15, 3, 10).Intersect(rectangle(8, 30, 7, 20))
	d := rectangle(8, 15, 7, 10)
	require.True(t, c.Equal(d))
}

func TestRectangleUnionMerge(t *testing.T) {
	a := union(
		rectangle(5, 15, 3, 6),
		rectangle(5, 30, 7, 10),
		rectangle(8, 30, 11, 20),
	)
	excepted := union(
		rectangle(5, 15, 3, 10),
		rectangle(8, 30, 7, 20),
	)
	require.True(t, excepted.Equal(a))
}

func TestRectangleSubtract(t *testing.T) {
	g := rectangle(5, 15, 3, 10).Subtract(rectangle(8, 30, 7, 20))
	h := union(
		rectangle(5, 7, 3, 10),
		rectangle(8, 15, 3, 6),
	)
	require.True(t, g.Equal(h))
}

func TestRectangleIntersectEmpty(t *testing.T) {
	a := rectangle(5, 15, 3, 10)
	b := union(
		rectangle(1, 3, 7, 20),
		rectangle(20, 23, 7, 20),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestRectangleUnion3(t *testing.T) {
	a := union(
		rectangle(1, 79, 10054, 10054),
		rectangle(80, 100, 10053, 10054),
		rectangle(101, 65535, 10054, 10054),
	)
	expected := union(
		rectangle(80, 100, 10053, 10053),
		rectangle(1, 65535, 10054, 10054),
	)
	require.True(t, expected.Equal(a))
}

func TestRectangleSwapDimensions(t *testing.T) {
	require.True(t, rectangle(1, 2, 3, 4).Swap().Equal(rectangle(3, 4, 1, 2)))
	require.True(t, rectangle(1, 2, 1, 2).Swap().Equal(rectangle(1, 2, 1, 2)))

	require.True(t, union(
		rectangle(1, 3, 7, 20),
		rectangle(20, 23, 7, 20),
	).Swap().Equal(union(
		rectangle(7, 20, 1, 3),
		rectangle(7, 20, 20, 23),
	)))
}

func TestAdditionalInterfaceFuncs(t *testing.T) {
	// NumPartitions()
	a := union(rectangle(1, 9, 1, 3), rectangle(8, 9, 4, 5)) // {(8-9 x 1-5) | (1-7 x 1-3)}
	require.Equal(t, 2, a.NumPartitions())
	require.Equal(t, 31, a.Size()) // 10 + 21 = 31

	// Left() , Right()
	leftSet := a.Left(interval.NewCanonicalSet())
	rightSet := a.Right(interval.NewCanonicalSet())
	require.True(t, interval.New(1, 9).ToSet().Equal(leftSet))
	require.True(t, interval.New(1, 5).ToSet().Equal(rightSet))
	fmt.Println(a.Left(interval.NewCanonicalSet()))  // 1-9
	fmt.Println(a.Right(interval.NewCanonicalSet())) // 1-5

	// NewProductLeft()
	z1 := ds.NewProductLeft[*interval.CanonicalSet, *interval.CanonicalSet]()
	require.True(t, z1.IsEmpty())
	fmt.Println(z1)

	// CartesianPairLeft()
	z2 := ds.CartesianPairLeft(interval.NewCanonicalSet(), interval.New(1, 9).ToSet())
	z3 := ds.CartesianPairLeft(interval.New(1, 9).ToSet(), interval.NewCanonicalSet())
	z4 := ds.CartesianPairLeft(interval.New(1, 9).ToSet(), interval.New(1, 9).ToSet())
	require.True(t, z2.IsEmpty())
	require.True(t, z3.IsEmpty())
	require.True(t, !z4.IsEmpty())

	fmt.Printf("done")
}

func TestUnion(t *testing.T) {
	a1 := rectangle(1, 3, 1, 5)
	a2 := rectangle(3, 3, 6, 7)
	a3 := a1.Union(a2)
	fmt.Println(a3) // {(3 x 1-7) | (1-2 x 1-5)}
	require.Equal(t, 2, a3.NumPartitions())
	a4 := rectangle(1, 2, 6, 7)
	// union below, identifies that the two paritions now have common 2-nd dimension sets, thus results in one parition object res
	fmt.Println(union(a4, a3)) // {(1-3 x 1-7)}
	require.Equal(t, 1, union(a4, a3).NumPartitions())
}

func TestSubtract(t *testing.T) {
	a1 := rectangle(1, 3, 1, 5)
	a2 := rectangle(3, 3, 1, 10)
	res1 := a1.Subtract(a2) // nothing remains on [3] in left set
	require.True(t, res1.Left(interval.NewCanonicalSet()).Equal(interval.New(1, 2).ToSet()))
	res2 := res1.Subtract(rectangle(1, 2, 1, 10))
	require.True(t, res2.IsEmpty())
	fmt.Println(res1) // {(1-2 x 1-5)}
	fmt.Println("done")

}
