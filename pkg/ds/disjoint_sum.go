/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// Disjoint is the union of two disjoint (tagged) sets L and R.
type Disjoint[L Set[L], R Set[R]] struct {
	left  L
	right R
}

// NewDisjoint creates a new Disjoint set from two sets.
// This is the only way to create a Disjoint set, since Go does not support creating an empty value of a generic type.
func NewDisjoint[L Set[L], R Set[R]](left L, right R) *Disjoint[L, R] {
	return &Disjoint[L, R]{left: left.Copy(), right: right.Copy()}
}

// Left returns the left-tagged set of the Disjoint set.
func (c *Disjoint[L, R]) Left() L {
	return c.left.Copy()
}

// Right returns the right-tagged set of the Disjoint set.
func (c *Disjoint[L, R]) Right() R {
	return c.right.Copy()
}

// Does not work: yields nil values for the other
// func NewLeft[L Set[L], R Set[R]](left L) *Disjoint[L, R] {
// 	return &Disjoint[L, R]{left: left}
// }
//
// func NewRight[L Set[L], R Set[R]](right R) *Disjoint[L, R] {
// 	return &Disjoint[L, R]{right: right}
// }

// Equal returns true if both left and right sets are equal to the other's left and right sets.
func (c *Disjoint[L, R]) Equal(other *Disjoint[L, R]) bool {
	return c.left.Equal(other.left) && c.right.Equal(other.right)
}

// Copy returns a deep copy of the Disjoint set.
func (c *Disjoint[L, R]) Copy() *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Copy(),
		right: c.right.Copy(),
	}
}

// Hash returns the hash value of the Disjoint set.
func (c *Disjoint[L, R]) Hash() int {
	return c.left.Hash() ^ c.right.Hash()
}

// IsEmpty returns true if both left and right sets are empty.
func (c *Disjoint[L, R]) IsEmpty() bool {
	return c.left.IsEmpty() && c.right.IsEmpty()
}

// Size returns the sum of the sizes of the left and right sets.
func (c *Disjoint[L, R]) Size() int {
	return c.left.Size() + c.right.Size()
}

// IsSubset returns true if both left and right sets are subsets of the other's left and right sets.
func (c *Disjoint[L, R]) IsSubset(other *Disjoint[L, R]) bool {
	return c.left.IsSubset(other.left) && c.right.IsSubset(other.right)
}

// Union returns the union of the left set with the other's left set and the right set with the other's right set.
func (c *Disjoint[L, R]) Union(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Union(other.left),
		right: c.right.Union(other.right),
	}
}

// Intersect returns the intersection of the left set with the other's left set and the right set with the other's right set.
func (c *Disjoint[L, R]) Intersect(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Intersect(other.left),
		right: c.right.Intersect(other.right),
	}
}

// Subtract returns the subtraction of the left set with the other's left set and the right set with the other's right set.
func (c *Disjoint[L, R]) Subtract(other *Disjoint[L, R]) *Disjoint[L, R] {
	return &Disjoint[L, R]{
		left:  c.left.Subtract(other.left),
		right: c.right.Subtract(other.right),
	}
}
