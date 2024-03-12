package hypercube_test

import (
	"fmt"
	"testing"

	"github.com/np-guard/models/pkg/hypercube"
	"github.com/np-guard/models/pkg/interval"
)

func TestHCBasic(t *testing.T) {
	cube1 := []*interval.CanonicalSet{interval.FromInterval(1, 100)}
	cube2 := []*interval.CanonicalSet{interval.FromInterval(1, 100)}
	cube3 := []*interval.CanonicalSet{interval.FromInterval(1, 200)}
	cube4 := []*interval.CanonicalSet{interval.FromInterval(1, 100), interval.FromInterval(1, 100)}
	cube5 := []*interval.CanonicalSet{interval.FromInterval(1, 100), interval.FromInterval(1, 100)}
	cube6 := []*interval.CanonicalSet{interval.FromInterval(1, 100), interval.FromInterval(1, 200)}

	a := hypercube.FromCube(cube1)
	b := hypercube.FromCube(cube2)
	c := hypercube.FromCube(cube3)
	d := hypercube.FromCube(cube4)
	e := hypercube.FromCube(cube5)
	f := hypercube.FromCube(cube6)

	if !a.Equal(b) {
		t.FailNow()
	}
	if a.Equal(c) {
		t.FailNow()
	}
	if b.Equal(c) {
		t.FailNow()
	}
	//nolint:all
	/*if !c.Equal(c) {
		t.FailNow()
	}
	if !a.Equal(a) {
		t.FailNow()
	}
	if !b.Equal(b) {
		t.FailNow()
	}*/
	if !d.Equal(e) {
		t.FailNow()
	}
	if !e.Equal(d) {
		t.FailNow()
	}
	if d.Equal(f) {
		t.FailNow()
	}
	if f.Equal(d) {
		t.FailNow()
	}
}

func TestCopy(t *testing.T) {
	cube1 := []*interval.CanonicalSet{interval.FromInterval(1, 100)}
	a := hypercube.FromCube(cube1)
	b := a.Copy()
	if !a.Equal(b) {
		t.FailNow()
	}
	if !b.Equal(a) {
		t.FailNow()
	}
	if a == b {
		t.FailNow()
	}
}

func TestString(t *testing.T) {
	cube1 := []*interval.CanonicalSet{interval.FromInterval(1, 100)}
	cube2 := []*interval.CanonicalSet{interval.FromInterval(1, 100), interval.FromInterval(1, 100)}
	a := hypercube.FromCube(cube1)
	b := hypercube.FromCube(cube2)
	fmt.Println(a.String())
	fmt.Println(b.String())
	fmt.Println("done")
}

func TestOr(t *testing.T) {
	cube1 := []*interval.CanonicalSet{interval.FromInterval(1, 100), interval.FromInterval(1, 100)}
	cube2 := []*interval.CanonicalSet{interval.FromInterval(1, 90), interval.FromInterval(1, 200)}
	a := hypercube.FromCube(cube1)
	b := hypercube.FromCube(cube2)
	c := a.Union(b)
	fmt.Println(a.String())
	fmt.Println(b.String())
	fmt.Println(c.String())
	fmt.Println("done")
}

func addCube1Dim(o *hypercube.CanonicalSet, start, end int64) *hypercube.CanonicalSet {
	cube := []*interval.CanonicalSet{interval.FromInterval(start, end)}
	a := hypercube.FromCube(cube)
	return o.Union(a)
}

func addCube2Dim(o *hypercube.CanonicalSet, start1, end1, start2, end2 int64) *hypercube.CanonicalSet {
	cube := []*interval.CanonicalSet{interval.FromInterval(start1, end1), interval.FromInterval(start2, end2)}
	a := hypercube.FromCube(cube)
	return o.Union(a)
}

func addCube3Dim(o *hypercube.CanonicalSet, s1, e1, s2, e2, s3, e3 int64) *hypercube.CanonicalSet {
	cube := []*interval.CanonicalSet{
		interval.FromInterval(s1, e1),
		interval.FromInterval(s2, e2),
		interval.FromInterval(s3, e3)}
	a := hypercube.FromCube(cube)
	return o.Union(a)
}

func TestBasic(t *testing.T) {
	a := hypercube.NewCanonicalSet(1)
	a = addCube1Dim(a, 1, 2)
	a = addCube1Dim(a, 5, 6)
	a = addCube1Dim(a, 3, 4)
	b := hypercube.NewCanonicalSet(1)
	b = addCube1Dim(b, 1, 6)
	if !a.Equal(b) {
		t.FailNow()
	}
}

func TestBasic2(t *testing.T) {
	a := hypercube.NewCanonicalSet(2)
	a = addCube2Dim(a, 1, 2, 1, 5)
	a = addCube2Dim(a, 1, 2, 7, 9)
	a = addCube2Dim(a, 1, 2, 6, 7)
	b := hypercube.NewCanonicalSet(2)
	b = addCube2Dim(b, 1, 2, 1, 9)
	if !a.Equal(b) {
		t.FailNow()
	}
}

func TestNew(t *testing.T) {
	a := hypercube.NewCanonicalSet(3)
	a = addCube3Dim(a, 10, 20, 10, 20, 1, 65535)
	a = addCube3Dim(a, 1, 65535, 15, 40, 1, 65535)
	a = addCube3Dim(a, 1, 65535, 100, 200, 30, 80)
	expectedStr := "[(1-9,21-65535),(100-200),(30-80)]; [(1-9,21-65535),(15-40),(1-65535)]"
	expectedStr += "; [(10-20),(10-40),(1-65535)]; [(10-20),(100-200),(30-80)]"
	actualStr := a.String()
	if actualStr != expectedStr {
		t.FailNow()
	}
	fmt.Println(a.String())
	fmt.Println("done")
}

func checkContained(t *testing.T, a, b *hypercube.CanonicalSet, expected bool) {
	t.Helper()
	contained, err := a.ContainedIn(b)
	if contained != expected || err != nil {
		t.FailNow()
	}
}

func checkEqual(t *testing.T, a, b *hypercube.CanonicalSet, expected bool) {
	t.Helper()
	res := a.Equal(b)
	if res != expected {
		t.FailNow()
	}
}

func TestContainedIn(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 200, 300)
	b := hypercube.FromCubeShort(10, 80, 210, 280)
	checkContained(t, b, a, true)
	b = addCube2Dim(b, 10, 200, 210, 280)
	checkContained(t, b, a, false)
}

func TestContainedIn2(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube2Dim(c, 150, 180, 20, 300)
	c = addCube2Dim(c, 200, 240, 200, 300)
	c = addCube2Dim(c, 241, 300, 200, 350)

	a := hypercube.FromCubeShort(1, 100, 200, 300)
	a = addCube2Dim(a, 150, 180, 20, 300)
	a = addCube2Dim(a, 200, 240, 200, 300)
	a = addCube2Dim(a, 242, 300, 200, 350)

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
	b = addCube2Dim(b, 205, 205, 0, 53)
	b = addCube2Dim(b, 205, 205, 55, 255)
	b = addCube2Dim(b, 206, 254, 0, 255)
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

func TestEqual(t *testing.T) {
	a := hypercube.FromCubeShort(1, 2)
	b := hypercube.FromCubeShort(1, 2)
	checkEqual(t, a, b, true)
	c := hypercube.FromCubeShort(1, 2, 1, 5)
	d := hypercube.FromCubeShort(1, 2, 1, 5)
	checkEqual(t, c, d, true)
	c = addCube2Dim(c, 1, 2, 7, 9)
	c = addCube2Dim(c, 1, 2, 6, 7)
	c = addCube2Dim(c, 4, 8, 1, 9)
	res := hypercube.FromCubeShort(4, 8, 1, 9)
	res = addCube2Dim(res, 1, 2, 1, 9)
	checkEqual(t, res, c, true)

	a = addCube1Dim(a, 5, 6)
	a = addCube1Dim(a, 3, 4)
	res1 := hypercube.FromCubeShort(1, 6)
	checkEqual(t, res1, a, true)

	d = addCube2Dim(d, 1, 2, 1, 5)
	d = addCube2Dim(d, 5, 6, 1, 5)
	d = addCube2Dim(d, 3, 4, 1, 5)
	res2 := hypercube.FromCubeShort(1, 6, 1, 5)
	checkEqual(t, res2, d, true)
}

func TestBasicAddCube(t *testing.T) {
	a := hypercube.FromCubeShort(1, 2)
	a = addCube1Dim(a, 8, 10)
	b := a
	a = addCube1Dim(a, 1, 2)
	a = addCube1Dim(a, 6, 10)
	a = addCube1Dim(a, 1, 10)
	res := hypercube.FromCubeShort(1, 10)
	checkEqual(t, res, a, true)
	checkEqual(t, res, b, false)
}
func TestBasicAddHole(t *testing.T) {
	a := hypercube.FromCubeShort(1, 10)
	b := a.Subtract(hypercube.FromCubeShort(3, 20))
	c := a.Subtract(hypercube.FromCubeShort(0, 20))
	d := a.Subtract(hypercube.FromCubeShort(0, 5))
	e := a.Subtract(hypercube.FromCubeShort(12, 14))
	a = a.Subtract(hypercube.FromCubeShort(3, 7))
	f := hypercube.FromCubeShort(1, 2)
	f = addCube1Dim(f, 8, 10)
	checkEqual(t, a, f, true)
	checkEqual(t, b, hypercube.FromCubeShort(1, 2), true)
	checkEqual(t, c, hypercube.NewCanonicalSet(1), true)
	checkEqual(t, d, hypercube.FromCubeShort(6, 10), true)
	checkEqual(t, e, hypercube.FromCubeShort(1, 10), true)
}

func TestAddHoleBasic2(t *testing.T) {
	a := hypercube.FromCubeShort(1, 100, 200, 300)
	b := a.Copy()
	c := a.Copy()
	a = a.Subtract(hypercube.FromCubeShort(50, 60, 220, 300))
	resA := hypercube.FromCubeShort(61, 100, 200, 300)
	resA = addCube2Dim(resA, 50, 60, 200, 219)
	resA = addCube2Dim(resA, 1, 49, 200, 300)
	checkEqual(t, a, resA, true)

	b = b.Subtract(hypercube.FromCubeShort(50, 1000, 0, 250))
	resB := hypercube.FromCubeShort(50, 100, 251, 300)
	resB = addCube2Dim(resB, 1, 49, 200, 300)
	checkEqual(t, b, resB, true)

	c = addCube2Dim(c, 400, 700, 200, 300)
	c = c.Subtract(hypercube.FromCubeShort(50, 1000, 0, 250))
	resC := hypercube.FromCubeShort(50, 100, 251, 300)
	resC = addCube2Dim(resC, 1, 49, 200, 300)
	resC = addCube2Dim(resC, 400, 700, 251, 300)
	checkEqual(t, c, resC, true)
}

func TestAddHole(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = c.Subtract(hypercube.FromCubeShort(50, 60, 220, 300))
	d := hypercube.FromCubeShort(1, 49, 200, 300)
	d = addCube2Dim(d, 50, 60, 200, 219)
	d = addCube2Dim(d, 61, 100, 200, 300)
	checkEqual(t, c, d, true)
}

func TestAddHole2(t *testing.T) {
	c := hypercube.FromCubeShort(80, 100, 20, 300)
	c = addCube2Dim(c, 250, 400, 20, 300)
	c = c.Subtract(hypercube.FromCubeShort(30, 300, 100, 102))
	d := hypercube.FromCubeShort(80, 100, 20, 99)
	d = addCube2Dim(d, 80, 100, 103, 300)
	d = addCube2Dim(d, 250, 300, 20, 99)
	d = addCube2Dim(d, 250, 300, 103, 300)
	d = addCube2Dim(d, 301, 400, 20, 300)
	checkEqual(t, c, d, true)
}
func TestAddHole3(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = c.Subtract(hypercube.FromCubeShort(1, 100, 200, 300))
	checkEqual(t, c, hypercube.NewCanonicalSet(2), true)
}

func TestIntervalsUnion(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube2Dim(c, 101, 200, 200, 300)
	d := hypercube.FromCubeShort(1, 200, 200, 300)
	checkEqual(t, c, d, true)
	if c.String() != d.String() {
		t.FailNow()
	}
}

func TestIntervalsUnion2(t *testing.T) {
	c := hypercube.FromCubeShort(1, 100, 200, 300)
	c = addCube2Dim(c, 101, 200, 200, 300)
	c = addCube2Dim(c, 201, 300, 200, 300)
	c = addCube2Dim(c, 301, 400, 200, 300)
	c = addCube2Dim(c, 402, 500, 200, 300)
	c = addCube2Dim(c, 500, 600, 200, 700)
	c = addCube2Dim(c, 601, 700, 200, 700)

	d := c.Copy()
	d = addCube2Dim(d, 702, 800, 200, 700)

	cExpected := hypercube.FromCubeShort(1, 400, 200, 300)
	cExpected = addCube2Dim(cExpected, 402, 500, 200, 300)
	cExpected = addCube2Dim(cExpected, 500, 700, 200, 700)
	dExpected := cExpected.Copy()
	dExpected = addCube2Dim(dExpected, 702, 800, 200, 700)
	checkEqual(t, c, cExpected, true)
	checkEqual(t, d, dExpected, true)
}

func TestAndSubOr(t *testing.T) {
	a := hypercube.FromCubeShort(5, 15, 3, 10)
	b := hypercube.FromCubeShort(8, 30, 7, 20)

	c := a.Intersect(b)
	d := hypercube.FromCubeShort(8, 15, 7, 10)
	checkEqual(t, c, d, true)

	f := a.Union(b)
	e := hypercube.FromCubeShort(5, 15, 3, 6)
	e = addCube2Dim(e, 5, 30, 7, 10)
	e = addCube2Dim(e, 8, 30, 11, 20)
	checkEqual(t, e, f, true)

	g := a.Subtract(b)
	h := hypercube.FromCubeShort(5, 7, 3, 10)
	h = addCube2Dim(h, 8, 15, 3, 6)
	checkEqual(t, g, h, true)
}

func TestAnd2(t *testing.T) {
	a := hypercube.FromCubeShort(5, 15, 3, 10)
	b := hypercube.FromCubeShort(1, 3, 7, 20)
	b = addCube2Dim(b, 20, 23, 7, 20)
	c := a.Intersect(b)
	checkEqual(t, c, hypercube.NewCanonicalSet(2), true)
}

func TestOr2(t *testing.T) {
	a := hypercube.FromCubeShort(80, 100, 10053, 10053)
	b := hypercube.FromCubeShort(1, 65535, 10054, 10054)
	a = a.Union(b)
	expected := hypercube.FromCubeShort(1, 79, 10054, 10054)
	expected = addCube2Dim(expected, 80, 100, 10053, 10054)
	expected = addCube2Dim(expected, 101, 65535, 10054, 10054)
	checkEqual(t, a, expected, true)
}
