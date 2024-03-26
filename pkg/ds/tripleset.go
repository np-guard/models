// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

type TripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m *Product[S1, *Product[S2, S3]]
}

func NewTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *TripleSet[S1, S2, S3] {
	return &TripleSet[S1, S2, S3]{m: NewIMap[S1, *Product[S2, S3]]()}
}

func Path[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *TripleSet[S1, S2, S3] {
	return &TripleSet[S1, S2, S3]{m: IPath(s1, IPath(s2, s3))}
}

type Triple[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	S1 S1
	S2 S2
	S3 S3
}

func (c *TripleSet[S1, S2, S3]) Equal(other *TripleSet[S1, S2, S3]) bool {
	return c.m.Equal(other.m)
}

func (c *TripleSet[S1, S2, S3]) Hash() int {
	return c.m.Hash()
}

func (c *TripleSet[S1, S2, S3]) Copy() *TripleSet[S1, S2, S3] {
	return &TripleSet[S1, S2, S3]{
		m: c.m.Copy(),
	}
}

func (c *TripleSet[S1, S2, S3]) Intersect(other *TripleSet[S1, S2, S3]) *TripleSet[S1, S2, S3] {
	return &TripleSet[S1, S2, S3]{m: c.m.Intersect(other.m)}
}

func (c *TripleSet[S1, S2, S3]) IsEmpty() bool {
	return c.m.IsEmpty()
}

func (c *TripleSet[S1, S2, S3]) Union(other *TripleSet[S1, S2, S3]) *TripleSet[S1, S2, S3] {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &TripleSet[S1, S2, S3]{
		m: c.m.Union(other.m),
	}
}

func (c *TripleSet[S1, S2, S3]) Subtract(other *TripleSet[S1, S2, S3]) *TripleSet[S1, S2, S3] {
	if c.IsEmpty() {
		return NewTripleSet[S1, S2, S3]()
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	return &TripleSet[S1, S2, S3]{m: c.m.Subtract(other.m)}
}

// ContainedIn returns true if c is subset of other
func (c *TripleSet[S1, S2, S3]) ContainedIn(other *TripleSet[S1, S2, S3]) bool {
	return c.m.ContainedIn(other.m)
}

func (c *TripleSet[S1, S2, S3]) Triples() []Triple[S1, S2, S3] {
	res := []Triple[S1, S2, S3]{}
	for _, outer := range c.m.Pairs() {
		for _, inner := range outer.Value.m.Pairs() {
			res = append(res, Triple[S1, S2, S3]{
				S1: outer.Key.Copy(),
				S2: inner.Key.Copy(),
				S3: inner.Value.Copy(),
			})
		}
	}
	return res
}
