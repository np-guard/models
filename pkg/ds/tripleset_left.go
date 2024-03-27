// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

// LeftTripleSet is a left-associative 3-product of sets (S1 x S2) x S3
type LeftTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m *Product[*Product[S1, S2], S3]
}

func NewLeftTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *LeftTripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: NewProduct[*Product[S1, S2], S3]()}
}

func CartesianLeftTriple[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *LeftTripleSet[S1, S2, S3] {
	return &LeftTripleSet[S1, S2, S3]{m: CartesianPair(CartesianPair(s1, s2), s3)}
}

func (c *LeftTripleSet[S1, S2, S3]) Equal(other TripleSet[S1, S2, S3]) bool {
	r, ok := other.(*LeftTripleSet[S1, S2, S3])
	if !ok {
		return false
	}
	return c.m.Equal(r.m)
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
	res := NewLeftTripleSet[S1, S2, S3]()
	for _, p := range other.Partitions() {
		res = res.Union(CartesianLeftTriple(p.S1, p.S2, p.S3)).(*LeftTripleSet[S1, S2, S3])
	}
	return res
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
	res := []Triple[S1, S2, S3]{}
	for _, outer := range c.m.Partitions() {
		for _, inner := range outer.Key.Partitions() {
			res = append(res, Triple[S1, S2, S3]{
				S1: inner.Key.Copy(),
				S2: inner.Value.Copy(),
				S3: outer.Value.Copy(),
			})
		}
	}
	return res
}

// Swap returns a new Product object, built from the input Product object,
// with left and right swapped
func (c *LeftTripleSet[S1, S2, S3]) Swap23() TripleSet[S1, S3, S2] {
	res := NewLeftTripleSet[S1, S3, S2]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianLeftTriple(triple.S1, triple.S3, triple.S2)).(*LeftTripleSet[S1, S3, S2])
	}
	return res
}

func (c *LeftTripleSet[S1, S2, S3]) Swap12() TripleSet[S2, S1, S3] {
	res := NewLeftTripleSet[S2, S1, S3]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianLeftTriple(triple.S2, triple.S1, triple.S3)).(*LeftTripleSet[S2, S1, S3])
	}
	return res
}

func (c *LeftTripleSet[S1, S2, S3]) Swap13() TripleSet[S3, S2, S1] {
	res := NewLeftTripleSet[S3, S2, S1]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianLeftTriple(triple.S3, triple.S2, triple.S1)).(*LeftTripleSet[S3, S2, S1])
	}
	return res
}
