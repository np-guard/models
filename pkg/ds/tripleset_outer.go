/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

import (
	"fmt"
	"strings"
)

// OuterTripleSet is an outer-associative 3-product of sets (S1 x S3) x S2,  created as LeftTripleSet[S1, S3, S2] (Product[Product[S1, S3], S2])
type OuterTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	m *LeftTripleSet[S1, S3, S2]
}

func NewOuterTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]]() *OuterTripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: NewLeftTripleSet[S1, S3, S2]()}
}

func CartesianOuterTriple[S1 Set[S1], S2 Set[S2], S3 Set[S3]](s1 S1, s2 S2, s3 S3) *OuterTripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: CartesianLeftTriple(s1, s3, s2)}
}

func (c *OuterTripleSet[S1, S2, S3]) Equal(other TripleSet[S1, S2, S3]) bool {
	return c.m.Equal(asOuterTripleSet(other).m)
}

func (c *OuterTripleSet[S1, S2, S3]) Copy() TripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: c.m.Copy().(*LeftTripleSet[S1, S3, S2])}
}

func (c *OuterTripleSet[S1, S2, S3]) Hash() int {
	return c.m.Hash()
}

func (c *OuterTripleSet[S1, S2, S3]) IsEmpty() bool {
	return c.m.IsEmpty()
}

func (c *OuterTripleSet[S1, S2, S3]) Size() int {
	return c.m.Size()
}

func asOuterTripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]](other TripleSet[S1, S2, S3]) *OuterTripleSet[S1, S2, S3] {
	r, ok := other.(*OuterTripleSet[S1, S2, S3])
	if ok {
		return r
	}
	var res TripleSet[S1, S2, S3] = NewOuterTripleSet[S1, S2, S3]()
	for _, p := range other.Partitions() {
		res = res.Union(CartesianOuterTriple(p.S1, p.S2, p.S3))
	}
	return res.(*OuterTripleSet[S1, S2, S3])
}

// IsSubset returns true if c is subset of other
func (c *OuterTripleSet[S1, S2, S3]) IsSubset(other TripleSet[S1, S2, S3]) bool {
	return c.m.IsSubset(asOuterTripleSet(other).m)
}

func (c *OuterTripleSet[S1, S2, S3]) Union(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: c.m.Union(asOuterTripleSet(other).m).(*LeftTripleSet[S1, S3, S2])}
}

func (c *OuterTripleSet[S1, S2, S3]) Intersect(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: c.m.Intersect(asOuterTripleSet(other).m).(*LeftTripleSet[S1, S3, S2])}
}

func (c *OuterTripleSet[S1, S2, S3]) Subtract(other TripleSet[S1, S2, S3]) TripleSet[S1, S2, S3] {
	return &OuterTripleSet[S1, S2, S3]{m: c.m.Subtract(asOuterTripleSet(other).m).(*LeftTripleSet[S1, S3, S2])}
}

func (c *OuterTripleSet[S1, S2, S3]) Partitions() []Triple[S1, S2, S3] {
	return partitionsMap(c.m, Triple[S1, S3, S2].Swap23)
}

func (c *OuterTripleSet[S1, S2, S3]) String() string {
	partitions := c.Partitions()
	partitionsStrings := make([]string, len(partitions))
	for i, triple := range partitions {
		partitionsStrings[i] = fmt.Sprintf("(%s x %s x %s)", triple.S1.String(), triple.S2.String(), triple.S3.String())
	}
	return "{" + strings.Join(partitionsStrings, " | ") + "}"
}

/*func (c *OuterTripleSet[S1, S2, S3]) NumPartitions() int {
	res := 0
	for _, p := range c.m.Partitions() {
		res += p.Left.NumPartitions()
	}
	return res
	return 0
}*/
