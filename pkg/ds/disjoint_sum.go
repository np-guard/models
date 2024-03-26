// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

// Fix: handle empty cases (usually nil values)
type Disjoint[L Set[L], R Set[R]] struct {
	left  L
	right R
}

func NewDisjoint[L Set[L], R Set[R]]() *Disjoint[L, R] {
	return &Disjoint[L, R]{}
}

func NewLeft[L Set[L], R Set[R]](left L) *Disjoint[L, R] {
	return &Disjoint[L, R]{left: left}
}

func NewRight[L Set[L], R Set[R]](right R) *Disjoint[L, R] {
	return &Disjoint[L, R]{right: right}
}

func (c *Disjoint[L, R]) Equal(other *Disjoint[L, R]) bool {
	return c.left.Equal(other.left) && c.right.Equal(other.right)
}

func (c *Disjoint[L, R]) Hash() int {
	return c.left.Hash() ^ c.right.Hash()
}
func (c *Disjoint[L, R]) Copy() *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Copy(),
		right: c.right.Copy(),
	}
}

func (c *Disjoint[L, R]) Intersect(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Intersect(other.left),
		right: c.right.Intersect(other.right),
	}
}

func (c *Disjoint[L, R]) IsEmpty() bool {
	return c.left.IsEmpty() && c.right.IsEmpty()
}

func (c *Disjoint[L, R]) Union(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Union(other.left),
		right: c.right.Union(other.right),
	}
}

func (c *Disjoint[L, R]) Subtract(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Subtract(other.left),
		right: c.right.Subtract(other.right),
	}
}

// ContainedIn returns true if c is subset of other
func (c *Disjoint[L, R]) ContainedIn(other *Disjoint[L, R]) bool {
	return c.left.ContainedIn(other.left) && c.right.ContainedIn(other.right)
}

func (c *Disjoint[L, R]) Left() L {
	return c.left.Copy()
}

func (c *Disjoint[L, R]) Right() R {
	return c.right.Copy()
}

func (c *Disjoint[L, R]) Size() int {
	return c.left.Size() + c.right.Size()
}
