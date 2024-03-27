// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

// RightTripleSet is a right-associative 3-product of sets S1 x (S2 x S3)
type RightTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m *Product[S1, *Product[S2, S3]]
}

func NewRightTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *RightTripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: NewProduct[S1, *Product[S2, S3]]()}
}

func CartesianRightTriple[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *RightTripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: CartesianPair(s1, CartesianPair(s2, s3))}
}

func (c *RightTripleSet[S1, S2, S3]) Equal(other TripleSet[S1, S2, S3]) bool {
	r, ok := other.(*RightTripleSet[S1, S2, S3])
	if !ok {
		return false
	}
	return c.m.Equal(r.m)
}

func (c *RightTripleSet[S1, S2, S3]) Copy() TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Copy()}
}

func (c *RightTripleSet[S1, S2, S3]) Hash() int {
	return c.m.Hash()
}

func (c *RightTripleSet[S1, S2, S3]) IsEmpty() bool {
	return c.m.IsEmpty()
}

func (c *RightTripleSet[S1, S2, S3]) Size() int {
	return c.m.Size()
}

func asRightTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]](other TripleSet[S1, S2, S3]) *RightTripleSet[S1, S2, S3] {
	r, ok := other.(*RightTripleSet[S1, S2, S3])
	if ok {
		return r
	}
	res := NewRightTripleSet[S1, S2, S3]()
	for _, p := range other.Partitions() {
		res = res.Union(CartesianRightTriple(p.S1, p.S2, p.S3)).(*RightTripleSet[S1, S2, S3])
	}
	return res
}

// IsSubset returns true if c is subset of other
func (c *RightTripleSet[S1, S2, S3]) IsSubset(other TripleSet[S1, S2, S3]) bool {
	return c.m.IsSubset(asRightTripleSet(other).m)
}

func (c *RightTripleSet[S1, S2, S3]) Union(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Union(asRightTripleSet(other).m)}
}

func (c *RightTripleSet[S1, S2, S3]) Intersect(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Intersect(asRightTripleSet(other).m)}
}

func (c *RightTripleSet[S1, S2, S3]) Subtract(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Subtract(asRightTripleSet(other).m)}
}

func (c *RightTripleSet[S1, S2, S3]) Partitions() []Triple[S1, S2, S3] {
	res := []Triple[S1, S2, S3]{}
	for _, outer := range c.m.Partitions() {
		for _, inner := range outer.Value.Partitions() {
			res = append(res, Triple[S1, S2, S3]{
				S1: outer.Key.Copy(),
				S2: inner.Key.Copy(),
				S3: inner.Value.Copy(),
			})
		}
	}
	return res
}

// Swap returns a new Product object, built from the input Product object,
// with left and right swapped
func (c *RightTripleSet[S1, S2, S3]) Swap23() TripleSet[S1, S3, S2] {
	res := NewRightTripleSet[S1, S3, S2]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S1, triple.S3, triple.S2)).(*RightTripleSet[S1, S3, S2])
	}
	return res
}

func (c *RightTripleSet[S1, S2, S3]) Swap12() TripleSet[S2, S1, S3] {
	res := NewRightTripleSet[S2, S1, S3]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S2, triple.S1, triple.S3)).(*RightTripleSet[S2, S1, S3])
	}
	return res
}

func (c *RightTripleSet[S1, S2, S3]) Swap13() TripleSet[S3, S2, S1] {
	res := NewRightTripleSet[S3, S2, S1]()
	for _, triple := range c.Partitions() {
		res = res.Union(CartesianRightTriple(triple.S3, triple.S2, triple.S1)).(*RightTripleSet[S3, S2, S1])
	}
	return res
}
