/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// RightTripleSet is a right-associative 3-product of sets S1 x (S2 x S3),
// created as LeftTripleSet[S2, S3, S1] (Product[Product[S2, S3], S1])
type RightTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m *LeftTripleSet[S2, S3, S1]
}

func NewRightTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *RightTripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: NewLeftTripleSet[S2, S3, S1]()}
}

func CartesianRightTriple[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *RightTripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: CartesianLeftTriple(s2, s3, s1)}
}

func (c *RightTripleSet[S1, S2, S3]) Equal(other TripleSet[S1, S2, S3]) bool {
	return c.m.Equal(asRightTripleSet(other).m)
}

func (c *RightTripleSet[S1, S2, S3]) Copy() TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Copy().(*LeftTripleSet[S2, S3, S1])}
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
	var res TripleSet[S1, S2, S3] = NewRightTripleSet[S1, S2, S3]()
	for _, p := range other.Partitions() {
		res = res.Union(CartesianRightTriple(p.S1, p.S2, p.S3))
	}
	return res.(*RightTripleSet[S1, S2, S3])
}

// IsSubset returns true if c is subset of other
func (c *RightTripleSet[S1, S2, S3]) IsSubset(other TripleSet[S1, S2, S3]) bool {
	return c.m.IsSubset(asRightTripleSet(other).m)
}

func (c *RightTripleSet[S1, S2, S3]) Union(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Union(asRightTripleSet(other).m).(*LeftTripleSet[S2, S3, S1])}
}

func (c *RightTripleSet[S1, S2, S3]) Intersect(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Intersect(asRightTripleSet(other).m).(*LeftTripleSet[S2, S3, S1])}
}

func (c *RightTripleSet[S1, S2, S3]) Subtract(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &RightTripleSet[S1, S2, S3]{m: c.m.Subtract(asRightTripleSet(other).m).(*LeftTripleSet[S2, S3, S1])}
}

func (c *RightTripleSet[S1, S2, S3]) Partitions() []Triple[S1, S2, S3] {
	return partitionsMap(c.m, Triple[S2, S3, S1].ShiftRight)
}

func (c *RightTripleSet[S1, S2, S3]) String() string {
	partitions := c.Partitions()
	partitionsStrings := make([]string, len(partitions))
	for i, triple := range partitions {
		partitionsStrings[i] = tupleString(triple.S1.String(), triple.S2.String(), triple.S3.String())
	}
	return setString(partitionsStrings)
}

/*func (c *RightTripleSet[S1, S2, S3]) NumPartitions() int {
	res := 0
	for _, p := range c.m.Partitions() {
		res += p.Left().NumPartitions()
	}
	return res
}*/
