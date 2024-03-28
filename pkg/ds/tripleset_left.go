// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ds

// LeftTripleSet is a left-associative 3-product of sets (S1 x S2) x S3
type LeftTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m Product[Product[S1, S2], S3]
}

func NewLeftTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *LeftTripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: NewProductLeft[Product[S1, S2], S3]()}
}

func CartesianLeftTriple[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *LeftTripleSet[S1, S2, S3] {
	var l Product[S1, S2] = CartesianPairLeft(s1, s2)
	var r Product[Product[S1, S2], S3] = CartesianPairLeft(l, s3)
	return &LeftTripleSet[S1, S2, S3]{m: r}
}

func (c *LeftTripleSet[S1, S2, S3]) Equal(other TripleSet[S1, S2, S3]) bool {
	return c.m.Equal(asLeftTripleSet(other).m)
}

func (c *LeftTripleSet[S1, S2, S3]) Copy() TripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: c.m.Copy()}
}

func (c *LeftTripleSet[S1, S2, S3]) Hash() int {
	return c.m.Hash()
}

func (c *LeftTripleSet[S1, S2, S3]) IsEmpty() bool {
	return c.m.IsEmpty()
}

func (c *LeftTripleSet[S1, S2, S3]) Size() int {
	return c.m.Size()
}

func asLeftTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]](other TripleSet[S1, S2, S3]) *LeftTripleSet[S1, S2, S3] {
	r, ok := other.(*LeftTripleSet[S1, S2, S3])
	if ok {
		return r
	}
	var res TripleSet[S1, S2, S3] = NewLeftTripleSet[S1, S2, S3]()
	for _, p := range other.Partitions() {
		res = res.Union(CartesianLeftTriple(p.S1, p.S2, p.S3))
	}
	return res.(*LeftTripleSet[S1, S2, S3])
}

// IsSubset returns true if c is subset of other
func (c *LeftTripleSet[S1, S2, S3]) IsSubset(other TripleSet[S1, S2, S3]) bool {
	return c.m.IsSubset(asLeftTripleSet(other).m)
}

func (c *LeftTripleSet[S1, S2, S3]) Union(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: c.m.Union(asLeftTripleSet(other).m)}
}

func (c *LeftTripleSet[S1, S2, S3]) Intersect(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: c.m.Intersect(asLeftTripleSet(other).m)}
}

func (c *LeftTripleSet[S1, S2, S3]) Subtract(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: c.m.Subtract(asLeftTripleSet(other).m)}
}

func (c *LeftTripleSet[S1, S2, S3]) Partitions() []Triple[S1, S2, S3] {
	var res []Triple[S1, S2, S3]
	for _, outer := range c.m.Partitions() {
		for _, inner := range outer.Left.Partitions() {
			res = append(res, Triple[S1, S2, S3]{
				S1: inner.Left.Copy(),
				S2: inner.Right.Copy(),
				S3: outer.Right.Copy(),
			})
		}
	}
	return res
}

// Swap12 returns a new TripleSet object, built from the input Product object,
// with S1 and S2 swapped
func (c *LeftTripleSet[S1, S2, S3]) Swap12() TripleSet[S2, S1, S3] {
	var res TripleSet[S2, S1, S3] = NewLeftTripleSet[S2, S1, S3]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S2, triple.S1, triple.S3))
	}
	return res
}

// Swap23 returns a new TripleSet object, built from the input Product object,
// with S2 and S3 swapped
func (c *LeftTripleSet[S1, S2, S3]) Swap23() TripleSet[S1, S3, S2] {
	var res TripleSet[S1, S3, S2] = NewLeftTripleSet[S1, S3, S2]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S1, triple.S3, triple.S2))
	}
	return res
}

// Swap13 returns a new TripleSet object, built from the input Product object,
// with S1 and S3 swapped
func (c *LeftTripleSet[S1, S2, S3]) Swap13() TripleSet[S3, S2, S1] {
	var res TripleSet[S3, S2, S1] = NewLeftTripleSet[S3, S2, S1]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S3, triple.S2, triple.S1))
	}
	return res
}
