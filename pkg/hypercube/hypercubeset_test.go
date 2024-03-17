package hypercube_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/hypercube"
	"github.com/np-guard/models/pkg/interval"
)

// cube returns a new hypercube.CanonicalSet created from a single input cube
// the input cube is given as an ordered list of integer values, where each two values
// represent the range (start,end) for a dimension value
func cube(values ...int64) *hypercube.CanonicalSet {
	cube := []*interval.CanonicalSet{}
	for i := 0; i < len(values); i += 2 {
		cube = append(cube, interval.FromInterval(values[i], values[i+1]))
	}
	return hypercube.FromCube(cube)
}

func union(set *hypercube.CanonicalSet, sets ...*hypercube.CanonicalSet) *hypercube.CanonicalSet {
	for _, c := range sets {
		set = set.Union(c)
	}
	return set.Copy()
}

func TestHCBasic(t *testing.T) {
	a := cube(1, 100)
	b := cube(1, 100)
	c := cube(1, 200)
	d := cube(1, 100, 1, 100)
	e := cube(1, 100, 1, 100)
	f := cube(1, 100, 1, 200)

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
	a := cube(1, 100)
	b := a.Copy()
	require.True(t, a.Equal(b))
	require.True(t, b.Equal(a))
	require.True(t, a != b)
}

func TestString(t *testing.T) {
	require.Equal(t, "[(1-3)]", cube(1, 3).String())
	require.Equal(t, "[(1-3),(2-4)]", cube(1, 3, 2, 4).String())
}

func TestOr(t *testing.T) {
	a := cube(1, 100, 1, 100)
	b := cube(1, 90, 1, 200)
	c := a.Union(b)
	require.Equal(t, "[(1-90),(1-200)]; [(91-100),(1-100)]", c.String())
}

func TestBasic1(t *testing.T) {
	a := union(
		cube(1, 2),
		cube(5, 6),
		cube(3, 4),
	)
	b := cube(1, 6)
	require.True(t, a.Equal(b))
}

func TestBasic2(t *testing.T) {
	a := union(
		cube(1, 2, 1, 5),
		cube(1, 2, 7, 9),
		cube(1, 2, 6, 7),
	)
	b := cube(1, 2, 1, 9)
	require.True(t, a.Equal(b))
}

func TestNew(t *testing.T) {
	a := union(
		cube(10, 20, 10, 20, 1, 65535),
		cube(1, 65535, 15, 40, 1, 65535),
		cube(1, 65535, 100, 200, 30, 80),
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
	a := cube(1, 100, 200, 300)
	b := cube(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = b.Union(cube(10, 200, 210, 280))
	checkContained(t, b, a, false)
}

func TestContainedIn2(t *testing.T) {
	c := union(
		cube(1, 100, 200, 300),
		cube(150, 180, 20, 300),
		cube(200, 240, 200, 300),
		cube(241, 300, 200, 350),
	)

	a := union(
		cube(1, 100, 200, 300),
		cube(150, 180, 20, 300),
		cube(200, 240, 200, 300),
		cube(242, 300, 200, 350),
	)
	d := cube(210, 220, 210, 280)
	e := cube(210, 310, 210, 280)
	f := cube(210, 250, 210, 280)
	f1 := cube(210, 240, 210, 280)
	f2 := cube(241, 250, 210, 280)

	checkContained(t, d, c, true)
	checkContained(t, e, c, false)
	checkContained(t, f1, c, true)
	checkContained(t, f2, c, true)
	checkContained(t, f, c, true)
	checkContained(t, f, a, false)
}

func TestContainedIn3(t *testing.T) {
	a := cube(105, 105, 54, 54)
	b := union(
		cube(0, 204, 0, 255),
		cube(205, 205, 0, 53),
		cube(205, 205, 55, 255),
		cube(206, 254, 0, 255),
	)
	checkContained(t, a, b, true)
}

func TestContainedIn4(t *testing.T) {
	a := cube(105, 105, 54, 54)
	b := cube(200, 204, 0, 255)
	checkContained(t, a, b, false)
}

func TestContainedIn5(t *testing.T) {
	a := cube(100, 200, 54, 65, 60, 300)
	b := cube(110, 120, 0, 10, 0, 255)
	checkContained(t, b, a, false)
}

func TestEqual1(t *testing.T) {
	a := cube(1, 2)
	b := cube(1, 2)
	require.True(t, a.Equal(b))

	c := cube(1, 2, 1, 5)
	d := cube(1, 2, 1, 5)
	require.True(t, c.Equal(d))
}

func TestEqual2(t *testing.T) {
	c := union(
		cube(1, 2, 1, 5),
		cube(1, 2, 7, 9),
		cube(1, 2, 6, 7),
		cube(4, 8, 1, 9),
	)
	res := union(
		cube(4, 8, 1, 9),
		cube(1, 2, 1, 9),
	)
	require.True(t, res.Equal(c))

	d := union(
		cube(1, 2, 1, 5),
		cube(5, 6, 1, 5),
		cube(3, 4, 1, 5),
	)
	res2 := cube(1, 6, 1, 5)
	require.True(t, res2.Equal(d))
}

func TestBasicAddCube(t *testing.T) {
	a := union(
		cube(1, 2),
		cube(8, 10),
	)
	b := union(
		a,
		cube(1, 2),
		cube(6, 10),
		cube(1, 10),
	)
	res := cube(1, 10)
	require.False(t, res.Equal(a))
	require.True(t, res.Equal(b))
}

func TestFourHoles(t *testing.T) {
	a := cube(1, 2, 1, 2)
	require.Equal(t, "[(1),(2)]; [(2),(1-2)]", a.Subtract(cube(1, 1, 1, 1)).String())
	require.Equal(t, "[(1),(1)]; [(2),(1-2)]", a.Subtract(cube(1, 1, 2, 2)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(2)]", a.Subtract(cube(2, 2, 1, 1)).String())
	require.Equal(t, "[(1),(1-2)]; [(2),(1)]", a.Subtract(cube(2, 2, 2, 2)).String())
}

func TestBasicSubtract1(t *testing.T) {
	a := cube(1, 10)
	require.True(t, a.Subtract(cube(3, 7)).Equal(union(cube(1, 2), cube(8, 10))))
	require.True(t, a.Subtract(cube(3, 20)).Equal(cube(1, 2)))
	require.True(t, a.Subtract(cube(0, 20)).IsEmpty())
	require.True(t, a.Subtract(cube(0, 5)).Equal(cube(6, 10)))
	require.True(t, a.Subtract(cube(12, 14)).Equal(cube(1, 10)))
}

func TestBasicSubtract2(t *testing.T) {
	a := cube(1, 100, 200, 300).Subtract(cube(50, 60, 220, 300))
	resA := union(
		cube(61, 100, 200, 300),
		cube(50, 60, 200, 219),
		cube(1, 49, 200, 300),
	)
	require.True(t, a.Equal(resA))

	b := cube(1, 100, 200, 300).Subtract(cube(50, 1000, 0, 250))
	resB := union(
		cube(50, 100, 251, 300),
		cube(1, 49, 200, 300),
	)
	require.True(t, b.Equal(resB))

	c := union(
		cube(1, 100, 200, 300),
		cube(400, 700, 200, 300),
	).Subtract(cube(50, 1000, 0, 250))
	resC := union(
		cube(50, 100, 251, 300),
		cube(1, 49, 200, 300),
		cube(400, 700, 251, 300),
	)
	require.True(t, c.Equal(resC))

	d := cube(1, 100, 200, 300).Subtract(cube(50, 60, 220, 300))
	dRes := union(
		cube(1, 49, 200, 300),
		cube(50, 60, 200, 219),
		cube(61, 100, 200, 300),
	)
	require.True(t, d.Equal(dRes))
}

func TestAddHole2(t *testing.T) {
	c := union(
		cube(80, 100, 20, 300),
		cube(250, 400, 20, 300),
	).Subtract(cube(30, 300, 100, 102))
	d := union(
		cube(80, 100, 20, 99),
		cube(80, 100, 103, 300),
		cube(250, 300, 20, 99),
		cube(250, 300, 103, 300),
		cube(301, 400, 20, 300),
	)
	require.True(t, c.Equal(d))
}

func TestSubtractToEmpty(t *testing.T) {
	c := cube(1, 100, 200, 300).Subtract(cube(1, 100, 200, 300))
	require.True(t, c.IsEmpty())
}

func TestUnion1(t *testing.T) {
	c := union(
		cube(1, 100, 200, 300),
		cube(101, 200, 200, 300),
	)
	cExpected := cube(1, 200, 200, 300)
	require.True(t, cExpected.Equal(c))
}

func TestUnion2(t *testing.T) {
	c := union(
		cube(1, 100, 200, 300),
		cube(101, 200, 200, 300),
		cube(201, 300, 200, 300),
		cube(301, 400, 200, 300),
		cube(402, 500, 200, 300),
		cube(500, 600, 200, 700),
		cube(601, 700, 200, 700),
	)
	cExpected := union(
		cube(1, 400, 200, 300),
		cube(402, 500, 200, 300),
		cube(500, 700, 200, 700),
	)
	require.True(t, c.Equal(cExpected))

	d := c.Union(cube(702, 800, 200, 700))
	dExpected := cExpected.Union(cube(702, 800, 200, 700))
	require.True(t, d.Equal(dExpected))
}

func TestIntersect(t *testing.T) {
	c := cube(5, 15, 3, 10).Intersect(cube(8, 30, 7, 20))
	d := cube(8, 15, 7, 10)
	require.True(t, c.Equal(d))
}

func TestUnionMerge(t *testing.T) {
	a := union(
		cube(5, 15, 3, 6),
		cube(5, 30, 7, 10),
		cube(8, 30, 11, 20),
	)
	excepted := union(
		cube(5, 15, 3, 10),
		cube(8, 30, 7, 20),
	)
	require.True(t, excepted.Equal(a))
}

func TestSubtract(t *testing.T) {
	g := cube(5, 15, 3, 10).Subtract(cube(8, 30, 7, 20))
	h := union(
		cube(5, 7, 3, 10),
		cube(8, 15, 3, 6),
	)
	require.True(t, g.Equal(h))
}

func TestIntersectEmpty(t *testing.T) {
	a := cube(5, 15, 3, 10)
	b := union(
		cube(1, 3, 7, 20),
		cube(20, 23, 7, 20),
	)
	c := a.Intersect(b)
	require.True(t, c.IsEmpty())
}

func TestOr2(t *testing.T) {
	a := union(
		cube(1, 79, 10054, 10054),
		cube(80, 100, 10053, 10054),
		cube(101, 65535, 10054, 10054),
	)
	expected := union(
		cube(80, 100, 10053, 10053),
		cube(1, 65535, 10054, 10054),
	)
	require.True(t, expected.Equal(a))
}
