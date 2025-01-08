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

// cubioidLeft returns a new ds.NProduct created from a single input cubioidLeft
// the input cubioidLeft is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func cubioidLeft(s1, e1, s2, e2, s3, e3 int64) ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianLeftTriple(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet(), interval.New(s3, e3).ToSet())
}

//nolint:unparam //  `e3` always receives `100` on current test examples
func cubioidRight(s1, e1, s2, e2, s3, e3 int64) ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianRightTriple(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet(), interval.New(s3, e3).ToSet())
}

func cubioidOuter(s1, e1, s2, e2, s3, e3 int64) ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet] {
	return ds.CartesianOuterTriple(interval.New(s1, e1).ToSet(), interval.New(s2, e2).ToSet(), interval.New(s3, e3).ToSet())
}

func TestCubioidEqual(t *testing.T) {
	d := cubioidLeft(1, 100, 1, 100, 2, 100)
	e := cubioidLeft(1, 100, 1, 100, 2, 100)
	f := cubioidLeft(1, 100, 1, 200, 2, 100)

	fmt.Println(d) // {(1-100 x 1-100 x 2-100)}
	fmt.Println(e) // {(1-100 x 1-100 x 2-100)}
	fmt.Println(f) // {(1-100 x 1-200 x 2-100)}

	require.True(t, d.Equal(e))
	require.True(t, e.Equal(d))

	require.False(t, d.Equal(f))
	require.False(t, f.Equal(d))
}

func TestCubioidLeftRightOuter(t *testing.T) {
	// simple case - one cube
	d1 := cubioidLeft(1, 100, 3, 100, 2, 100)
	d2 := cubioidRight(1, 100, 3, 100, 2, 100)
	d3 := cubioidOuter(1, 100, 3, 100, 2, 100)
	require.True(t, d1.Equal(d2))
	require.True(t, d1.Equal(d3))
	require.True(t, d2.Equal(d3))
	require.Equal(t, d1.String(), d2.String())
	require.Equal(t, d1.String(), d3.String())
	require.Equal(t, d2.String(), d3.String())

	d1Left := d1.(*ds.LeftTripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *interval.CanonicalSet])
	require.Equal(t, d1Left.S1(interval.NewCanonicalSet()), interval.New(1, 100).ToSet())
	require.Equal(t, d1Left.S2(interval.NewCanonicalSet()), interval.New(3, 100).ToSet())
	require.Equal(t, d1Left.S3(interval.NewCanonicalSet()), interval.New(2, 100).ToSet())

	// union with dimensions overlaps
	e1 := cubioidLeft(1, 100, 500, 500, 600, 600)
	d1e1 := union(d1, e1)
	fmt.Println(d1e1) // {(1-100 x 500 x 600) | (1-100 x 3-100 x 2-100)}
	require.Equal(t, 2, len(d1e1.Partitions()))

	e2 := cubioidLeft(20, 80, 101, 101, 1, 100)
	d1e2 := union(d1, e2)
	fmt.Println(d1e2) // {(20-80 x 101 x 1-100) | (1-100 x 3-100 x 2-100)}
	require.Equal(t, 2, len(d1e2.Partitions()))
	// note in the previous version this would be split to 3 cubes:
	//  [(1-19,81-100),(3-100),(2-100)]; [(20-80),(101),(1-100)]; [(20-80),(3-100),(2-100)]
	// in this version there is an overlap in the first dimension of the two cubes
	// that is because in the internal representation the tripleset is a map from product-set to set

	e3 := cubioidLeft(1, 100, 3, 100, 101, 101)
	d1e3 := union(d1, e3)
	fmt.Println(d1e3) // {(1-100 x 3-100 x 2-101)}
	require.Equal(t, 1, len(d1e3.Partitions()))

	e3Swapped := cubioidLeft(3, 100, 101, 101, 1, 100)
	d1Swapped := cubioidLeft(3, 100, 2, 100, 1, 100)
	d1e3Swapped := union(d1Swapped, e3Swapped)
	fmt.Println(d1e3Swapped) // {(3-100 x 2-101 x 1-100)}
	require.Equal(t, 1, len(d1e3Swapped.Partitions()))

	// e4 has a sub-cube on dims S1,S2 that is contained in d1's sub-cube, thus causing a split
	e4 := cubioidLeft(20, 80, 20, 80, 101, 101)
	d1e4 := union(d1, e4)
	fmt.Println(d1e4) // (1-19,81-100 x 3-100 x 2-100) | (20-80 x 3-19,81-100 x 2-100) | (20-80 x 20-80 x 2-101)}
	require.Equal(t, 3, len(d1e4.Partitions()))

	// right vs left example
	// same cubes, different partitions sets for right and left
	r1 := cubioidRight(3, 100, 2, 100, 1, 100)
	r2 := cubioidRight(101, 101, 20, 80, 1, 100)
	r3 := r1.Union(r2)
	fmt.Println(r3) // {(3-100 x 2-19,81-100 x 1-100) | (3-101 x 20-80 x 1-100)}

	l1 := cubioidLeft(3, 100, 2, 100, 1, 100)
	l2 := cubioidLeft(101, 101, 20, 80, 1, 100)
	l3 := l1.Union(l2)
	fmt.Println(l3) // {(101 x 20-80 x 1-100) | (3-100 x 2-100 x 1-100)}

	require.NotEqual(t, l3.String(), r3.String())
	require.True(t, l3.Equal(r3))
	require.True(t, r3.Equal(l3))

	// outer vs left, right example ("right" should be keeping orig cubes in this case)
	o1 := cubioidOuter(2, 100, 3, 100, 1, 100)
	o2 := cubioidOuter(20, 80, 101, 101, 1, 100)
	o3 := o1.Union(o2)
	fmt.Println(o3) // {(2-19,81-100 x 3-100 x 1-100) | (20-80 x 3-101 x 1-100)}

	l1 = cubioidLeft(2, 100, 3, 100, 1, 100)
	l2 = cubioidLeft(20, 80, 101, 101, 1, 100)
	l3 = l1.Union(l2)
	fmt.Println(l3) // {(2-19,81-100 x 3-100 x 1-100) | (20-80 x 3-101 x 1-100)}

	r1 = cubioidRight(2, 100, 3, 100, 1, 100)
	r2 = cubioidRight(20, 80, 101, 101, 1, 100)
	r3 = r1.Union(r2)
	fmt.Println(r3) // {(2-100 x 3-100 x 1-100) | (20-80 x 101 x 1-100)}

	require.NotEqual(t, o3.String(), r3.String())
	require.NotEqual(t, r3.String(), l3.String())
	require.Equal(t, o3.String(), l3.String())
	require.True(t, o3.Equal(r3))
	require.True(t, o3.Equal(l3))
	require.True(t, r3.Equal(l3))
}

func TestCanonicalRepCubioids(t *testing.T) {
	// S1,S2 are with different components, but equiv union sets of rectangles in few variations
	a1 := cubioidLeft(1, 7, 1, 3, 1, 1)
	a2 := cubioidLeft(8, 9, 1, 5, 1, 1)
	resA := union(a1, a2)
	fmt.Println(resA) // {(1-7 x 1-3 x 1) | (8-9 x 1-5 x 1)}

	b1 := cubioidLeft(1, 7, 1, 3, 1, 1)
	b2 := cubioidLeft(8, 9, 1, 3, 1, 1)
	b3 := cubioidLeft(8, 9, 4, 5, 1, 1)
	resB := union(b1, b2, b3)
	fmt.Println(resB) // {(8-9 x 1-5 x 1) | (1-7 x 1-3 x 1)}

	c1 := cubioidLeft(1, 9, 1, 3, 1, 1)
	c2 := cubioidLeft(8, 9, 4, 5, 1, 1)
	resC := union(c1, c2)
	fmt.Println(resC) // {(8-9 x 1-5 x 1) | (1-7 x 1-3 x 1)}

	require.True(t, resA.Equal(resB))
	require.True(t, resA.Equal(resC))
	require.True(t, resB.Equal(resC))

	// S1,S3 are with different components, but equiv union sets of rectangles in few variations
	a1 = cubioidLeft(1, 7, 1, 1, 1, 3)
	a2 = cubioidLeft(8, 9, 1, 1, 1, 5)
	resA = union(a1, a2)
	fmt.Println(resA) // {(8-9 x 1 x 1-5) | (1-7 x 1 x 1-3)}

	b1 = cubioidLeft(1, 7, 1, 1, 1, 3)
	b2 = cubioidLeft(8, 9, 1, 1, 1, 3)
	b3 = cubioidLeft(8, 9, 1, 1, 4, 5)
	resB = union(b1, b2, b3)
	fmt.Println(resB) // {(8-9 x 1 x 1-5) | (1-7 x 1 x 1-3)}

	c1 = cubioidLeft(1, 9, 1, 1, 1, 3)
	c2 = cubioidLeft(8, 9, 1, 1, 4, 5)
	resC = union(c1, c2)
	fmt.Println(resC) // {(8-9 x 1 x 1-5) | (1-7 x 1 x 1-3)}

	require.True(t, resA.Equal(resB))
	require.True(t, resA.Equal(resC))
	require.True(t, resB.Equal(resC))

	// S2,S3 are with different components, but equiv union sets of rectangles in few variations
	a1 = cubioidLeft(1, 1, 1, 7, 1, 3)
	a2 = cubioidLeft(1, 1, 8, 9, 1, 5)
	resA = union(a1, a2)
	fmt.Println(resA) // {(1 x 8-9 x 1-5) | (1 x 1-7 x 1-3)}

	b1 = cubioidLeft(1, 1, 1, 7, 1, 3)
	b2 = cubioidLeft(1, 1, 8, 9, 1, 3)
	b3 = cubioidLeft(1, 1, 8, 9, 4, 5)
	resB = union(b1, b2, b3)
	fmt.Println(resB) // {(1 x 8-9 x 1-5) | (1 x 1-7 x 1-3)}

	c1 = cubioidLeft(1, 1, 1, 9, 1, 3)
	c2 = cubioidLeft(1, 1, 8, 9, 4, 5)
	resC = union(c1, c2)
	fmt.Println(resC) // {(1 x 8-9 x 1-5) | (1 x 1-7 x 1-3)}

	require.True(t, resA.Equal(resB))
	require.True(t, resA.Equal(resC))
	require.True(t, resB.Equal(resC))
}

func TestCubioidCopy(t *testing.T) {
	a := cubioidLeft(1, 100, 3, 4, 5, 6)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestCubioidBasic1(t *testing.T) {
	a := union(
		cubioidLeft(1, 2, 3, 4, 5, 6),
		cubioidLeft(5, 6, 3, 4, 5, 6),
		cubioidLeft(3, 4, 3, 4, 5, 6),
	)
	b := cubioidLeft(1, 6, 3, 4, 5, 6)
	require.True(t, a.Equal(b))
}

func TestCubioidBasic2(t *testing.T) {
	a := union(
		cubioidLeft(1, 2, 1, 5, 0, 3),
		cubioidLeft(1, 2, 7, 9, 0, 3),
		cubioidLeft(1, 2, 6, 7, 0, 3),
	)
	b := cubioidLeft(1, 2, 1, 9, 0, 3)
	require.True(t, a.Equal(b))
}

func TestCubioidIsSubset(t *testing.T) {
	a := cubioidLeft(1, 100, 200, 300, 0, 3)
	b := cubioidLeft(10, 80, 210, 280, 0, 3)
	checkContained(t, b, a, true)
	checkContained(t, a, b, false)

	// check the delta
	c := a.Subtract(b)
	fmt.Println(c)           // {(1-9,81-100 x 200-300 x 0-3) | (10-80 x 200-209,281-300 x 0-3)}
	fmt.Println(union(c, b)) // {(1-100 x 200-300 x 0-3)}
	require.Equal(t, 2, len(c.Partitions()))
	require.True(t, union(c, b).Equal(a))

	b = b.Union(cubioidLeft(10, 200, 210, 280, 0, 3))
	checkContained(t, b, a, false)
	checkContained(t, a, b, false)

	// check the delta
	d1 := a.Subtract(b) // b is now cubioidLeft(10, 80, 210, 280, 0, 3) union cubioidLeft(10, 200, 210, 280, 0, 3)
	d2 := b.Subtract(a)
	fmt.Println(d1) // {(10-100 x 200-209,281-300 x 0-3) | (1-9 x 200-300 x 0-3)}
	fmt.Println(d2) // {(101-200 x 210-280 x 0-3)}
	e := a.Intersect(b)
	fmt.Println(e) // {(10-100 x 210-280 x 0-3)}
	require.True(t, union(d1, e).Equal(a))
	require.True(t, union(d2, e).Equal(b))
}

func TestCubioidIsSubset1(t *testing.T) {
	checkContained(t, cubioidLeft(1, 3, 0, 3, 3, 5), cubioidLeft(2, 4, 0, 3, 3, 5), false)
	checkContained(t, cubioidLeft(2, 4, 0, 3, 3, 5), cubioidLeft(1, 3, 0, 3, 3, 5), false)
	checkContained(t, cubioidLeft(1, 3, 0, 3, 3, 5), cubioidLeft(1, 4, 0, 3, 3, 5), true)
	checkContained(t, cubioidLeft(1, 4, 0, 3, 3, 5), cubioidLeft(1, 3, 0, 3, 3, 5), false)
}

func TestCubioidIsSubset2(t *testing.T) {
	c := union(
		cubioidLeft(1, 100, 200, 300, 3, 5),
		cubioidLeft(150, 180, 20, 300, 3, 5),
		cubioidLeft(200, 240, 200, 300, 3, 5),
		cubioidLeft(241, 300, 200, 350, 3, 5),
	)

	a := union(
		cubioidLeft(1, 100, 200, 300, 3, 5),
		cubioidLeft(150, 180, 20, 300, 3, 5),
		cubioidLeft(200, 240, 200, 300, 3, 5),
		cubioidLeft(242, 300, 200, 350, 3, 5),
	)
	d := cubioidLeft(210, 220, 210, 280, 3, 5)
	e := cubioidLeft(210, 310, 210, 280, 3, 5)
	f := cubioidLeft(210, 250, 210, 280, 3, 5)
	f1 := cubioidLeft(210, 240, 210, 280, 3, 5)
	f2 := cubioidLeft(241, 250, 210, 280, 3, 5)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestCubioidIsSubset3(t *testing.T) {
	a := cubioidLeft(105, 105, 54, 54, 3, 5)
	b := union(
		cubioidLeft(0, 204, 0, 255, 3, 5),
		cubioidLeft(205, 205, 0, 53, 3, 5),
		cubioidLeft(205, 205, 55, 255, 3, 5),
		cubioidLeft(206, 254, 0, 255, 3, 5),
	)
	checkContained(t, a, b, true)
}

func TestCubioidIsSubset4(t *testing.T) {
	a := cubioidLeft(105, 105, 54, 54, 3, 5)
	b := cubioidLeft(200, 204, 0, 255, 3, 5)
	checkContained(t, a, b, false)
}

func TestCubioidIsSubset5(t *testing.T) {
	a := cubioidLeft(100, 200, 54, 65, 60, 300)
	b := cubioidLeft(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestCubioidEqual1(t *testing.T) {
	a := cubioidLeft(1, 2, 3, 5, 3, 5)
	b := cubioidLeft(1, 2, 3, 5, 3, 5)
	require.True(t, a.Equal(b))

	c := cubioidLeft(1, 2, 1, 5, 3, 5)
	d := cubioidLeft(1, 2, 1, 5, 3, 5)
	require.True(t, c.Equal(d))
}

func TestCubioidEqual2(t *testing.T) {
	c := union(
		cubioidLeft(1, 2, 1, 5, 3, 5),
		cubioidLeft(1, 2, 7, 9, 3, 5),
		cubioidLeft(1, 2, 6, 7, 3, 5),
		cubioidLeft(4, 8, 1, 9, 3, 5),
	)
	res := union(
		cubioidLeft(4, 8, 1, 9, 3, 5),
		cubioidLeft(1, 2, 1, 9, 3, 5),
	)
	require.True(t, res.Equal(c))

	d := union(
		cubioidLeft(1, 2, 1, 5, 3, 5),
		cubioidLeft(5, 6, 1, 5, 3, 5),
		cubioidLeft(3, 4, 1, 5, 3, 5),
	)
	res2 := cubioidLeft(1, 6, 1, 5, 3, 5)
	require.True(t, res2.Equal(d))
}

func TestCubioidBasicAddCubioid(t *testing.T) {
	a := union(
		cubioidLeft(1, 2, 3, 5, 3, 5),
		cubioidLeft(8, 10, 3, 5, 3, 5),
	)
	b := union(
		a,
		cubioidLeft(1, 2, 3, 5, 3, 5),
		cubioidLeft(6, 10, 3, 5, 3, 5),
		cubioidLeft(1, 10, 3, 5, 3, 5),
	)
	res := cubioidLeft(1, 10, 3, 5, 3, 5)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestCubioidBasicSubtract1(t *testing.T) {
	a := cubioidLeft(1, 10, 3, 5, 3, 5)
	require.True(t, a.Subtract(cubioidLeft(3, 7, 3, 5, 3, 5)).Equal(union(cubioidLeft(1, 2, 3, 5, 3, 5), cubioidLeft(8, 10, 3, 5, 3, 5))))
	require.True(t, a.Subtract(cubioidLeft(3, 20, 3, 5, 3, 5)).Equal(cubioidLeft(1, 2, 3, 5, 3, 5)))
	require.True(t, a.Subtract(cubioidLeft(0, 20, 3, 5, 3, 5)).IsEmpty())
	require.True(t, a.Subtract(cubioidLeft(0, 5, 3, 5, 3, 5)).Equal(cubioidLeft(6, 10, 3, 5, 3, 5)))
	require.True(t, a.Subtract(cubioidLeft(12, 14, 3, 5, 3, 5)).Equal(cubioidLeft(1, 10, 3, 5, 3, 5)))
}

func TestCubioidBasicSubtract2(t *testing.T) {
	a := cubioidLeft(1, 100, 200, 300, 3, 5).Subtract(cubioidLeft(50, 60, 220, 300, 3, 5))
	resA := union(
		cubioidLeft(61, 100, 200, 300, 3, 5),
		cubioidLeft(50, 60, 200, 219, 3, 5),
		cubioidLeft(1, 49, 200, 300, 3, 5),
	)
	require.True(t, a.Equal(resA))

	b := cubioidLeft(1, 100, 200, 300, 3, 5).Subtract(cubioidLeft(50, 1000, 0, 250, 3, 5))
	resB := union(
		cubioidLeft(50, 100, 251, 300, 3, 5),
		cubioidLeft(1, 49, 200, 300, 3, 5),
	)
	require.True(t, b.Equal(resB))

	c := union(
		cubioidLeft(1, 100, 200, 300, 3, 5),
		cubioidLeft(400, 700, 200, 300, 3, 5),
	).Subtract(cubioidLeft(50, 1000, 0, 250, 3, 5))
	resC := union(
		cubioidLeft(50, 100, 251, 300, 3, 5),
		cubioidLeft(1, 49, 200, 300, 3, 5),
		cubioidLeft(400, 700, 251, 300, 3, 5),
	)
	require.True(t, c.Equal(resC))

	d := cubioidLeft(1, 100, 200, 300, 3, 5).Subtract(cubioidLeft(50, 60, 220, 300, 3, 5))
	dRes := union(
		cubioidLeft(1, 49, 200, 300, 3, 5),
		cubioidLeft(50, 60, 200, 219, 3, 5),
		cubioidLeft(61, 100, 200, 300, 3, 5),
	)
	require.True(t, d.Equal(dRes))
}

func TestCubioidAddHole2(t *testing.T) {
	c := union(
		cubioidLeft(80, 100, 20, 300, 3, 5),
		cubioidLeft(250, 400, 20, 300, 3, 5),
	).Subtract(cubioidLeft(30, 300, 100, 102, 3, 5))
	d := union(
		cubioidLeft(80, 100, 20, 99, 3, 5),
		cubioidLeft(80, 100, 103, 300, 3, 5),
		cubioidLeft(250, 300, 20, 99, 3, 5),
		cubioidLeft(250, 300, 103, 300, 3, 5),
		cubioidLeft(301, 400, 20, 300, 3, 5),
	)
	require.True(t, c.Equal(d))
}

func TestCubioidSubtractToEmpty(t *testing.T) {
	c := cubioidLeft(1, 100, 200, 300, 3, 5).Subtract(cubioidLeft(1, 100, 200, 300, 3, 5))
	require.True(t, c.IsEmpty())
}

func TestCubioidUnion1(t *testing.T) {
	c := union(
		cubioidLeft(1, 100, 200, 300, 3, 5),
		cubioidLeft(101, 200, 200, 300, 3, 5),
	)
	cExpected := cubioidLeft(1, 200, 200, 300, 3, 5)
	require.True(t, cExpected.Equal(c))
}

func TestCubioidUnion2(t *testing.T) {
	c := union(
		cubioidLeft(1, 100, 200, 300, 3, 5),
		cubioidLeft(101, 200, 200, 300, 3, 5),
		cubioidLeft(201, 300, 200, 300, 3, 5),
		cubioidLeft(301, 400, 200, 300, 3, 5),
		cubioidLeft(402, 500, 200, 300, 3, 5),
		cubioidLeft(500, 600, 200, 700, 3, 5),
		cubioidLeft(601, 700, 200, 700, 3, 5),
	)
	cExpected := union(
		cubioidLeft(1, 400, 200, 300, 3, 5),
		cubioidLeft(402, 500, 200, 300, 3, 5),
		cubioidLeft(500, 700, 200, 700, 3, 5),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(cubioidLeft(702, 800, 200, 700, 3, 5))
	dExpected := cExpected.Union(cubioidLeft(702, 800, 200, 700, 3, 5))
	require.True(t, d.Equal(dExpected))
}

func TestCubioidIntersect(t *testing.T) {
	c := cubioidLeft(5, 15, 3, 10, 3, 5).Intersect(cubioidLeft(8, 30, 7, 20, 3, 5))
	d := cubioidLeft(8, 15, 7, 10, 3, 5)
	require.True(t, c.Equal(d))
}

func TestCubioidUnionMerge(t *testing.T) {
	a := union(
		cubioidLeft(5, 15, 3, 6, 3, 5),
		cubioidLeft(5, 30, 7, 10, 3, 5),
		cubioidLeft(8, 30, 11, 20, 3, 5),
	)
	excepted := union(
		cubioidLeft(5, 15, 3, 10, 3, 5),
		cubioidLeft(8, 30, 7, 20, 3, 5),
	)
	require.True(t, excepted.Equal(a))
}

func TestCubioidSubtract(t *testing.T) {
	g := cubioidLeft(5, 15, 3, 10, 3, 5).Subtract(cubioidLeft(8, 30, 7, 20, 3, 5))
	h := union(
		cubioidLeft(5, 7, 3, 10, 3, 5),
		cubioidLeft(8, 15, 3, 6, 3, 5),
	)
	require.True(t, g.Equal(h))
}

func TestCubioidIntersectEmpty(t *testing.T) {
	a := cubioidLeft(5, 15, 3, 10, 3, 5)
	b := union(
		cubioidLeft(1, 3, 7, 20, 3, 5),
		cubioidLeft(20, 23, 7, 20, 3, 5),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestCubioidOr2(t *testing.T) {
	a := union(
		cubioidLeft(1, 79, 10054, 10054, 3, 5),
		cubioidLeft(80, 100, 10053, 10054, 3, 5),
		cubioidLeft(101, 65535, 10054, 10054, 3, 5),
	)
	expected := union(
		cubioidLeft(80, 100, 10053, 10053, 3, 5),
		cubioidLeft(1, 65535, 10054, 10054, 3, 5),
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

		require.True(t, ds.MapTripleSet(cubioidLeft(1, 2, 3, 4, 5, 6), Tr.Swap12).Equal(cubioidLeft(3, 4, 1, 2, 5, 6)))
		require.True(t, ds.MapTripleSet(cubioidLeft(1, 2, 3, 4, 5, 6), Tr.Swap23).Equal(cubioidLeft(1, 2, 5, 6, 3, 4)))
		require.True(t, ds.MapTripleSet(cubioidLeft(1, 2, 3, 4, 5, 6), Tr.Swap13).Equal(cubioidLeft(5, 6, 3, 4, 1, 2)))

		require.True(t, ds.MapTripleSet(union(
			cubioidLeft(1, 3, 7, 20, 3, 5),
			cubioidLeft(20, 23, 7, 20, 3, 5),
		), Tr.Swap12).Equal(union(
			cubioidLeft(7, 20, 1, 3, 3, 5),
			cubioidLeft(7, 20, 20, 23, 3, 5),
		)))
	}
}

func TestCartesianLeftTriple(t *testing.T) {
	z1 := ds.CartesianLeftTriple(interval.NewCanonicalSet(), interval.New(1, 9).ToSet(), interval.New(1, 9).ToSet())
	z2 := ds.CartesianLeftTriple(interval.New(1, 9).ToSet(), interval.NewCanonicalSet(), interval.New(1, 9).ToSet())
	z3 := ds.CartesianLeftTriple(interval.New(1, 9).ToSet(), interval.New(1, 9).ToSet(), interval.NewCanonicalSet())
	z4 := ds.CartesianLeftTriple(interval.New(1, 9).ToSet(), interval.New(1, 9).ToSet(), interval.New(1, 9).ToSet())

	require.True(t, z1.IsEmpty())
	require.True(t, z2.IsEmpty())
	require.True(t, z3.IsEmpty())
	require.True(t, !z4.IsEmpty())
	fmt.Println(z1) // {}
}
