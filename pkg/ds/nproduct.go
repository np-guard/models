// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

import (
	"log"
	"slices"
)

// NProduct is a dynamically-bounded product of sets
type NProduct[S Set[S]] struct {
	product    *Product[S, *NProduct[S]]
	dimensions int
}

// NewCanonicalSet returns a new empty NProduct with n dimensions
func NewCanonicalSet[S Set[S]](n int) *NProduct[S] {
	return &NProduct[S]{
		product:    NewProduct[S, *NProduct[S]](),
		dimensions: n,
	}
}

// Equal return true if c equals other (same canonical form)
func (c *NProduct[S]) Equal(other *NProduct[S]) bool {
	if c == other {
		return true
	}
	if c.dimensions != other.dimensions {
		return false
	}
	return c.product.Equal(other.product)
}

const (
	hashX = 3
	hashY = 5
)

func (c *NProduct[S]) Hash() int {
	if c.dimensions == 0 {
		return 1
	}
	res := hashX
	for _, p := range c.product.Partitions() {
		res ^= hashY*p.Value.dimensions + (p.Key.Hash() << 1) ^ p.Value.Hash()
	}
	return res
}

func (c *NProduct[S]) Size() int {
	res := 0
	if c.dimensions == 0 {
		return 1
	}
	for _, v := range c.product.Partitions() {
		res += v.Key.Size() * v.Value.Size()
	}
	return res
}

// Union returns a new NProduct object that results from union of c with other
func (c *NProduct[S]) Union(other *NProduct[S]) *NProduct[S] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &NProduct[S]{dimensions: other.dimensions, product: c.product.Union(other.product)}
}

// IsEmpty returns true if c is empty
func (c *NProduct[S]) IsEmpty() bool {
	return c.product.IsEmpty()
}

// Intersect returns a new NProduct object that results from intersection of c with other
func (c *NProduct[S]) Intersect(other *NProduct[S]) *NProduct[S] {
	if c == other {
		return c.Copy()
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &NProduct[S]{dimensions: other.dimensions, product: c.product.Intersect(other.product)}
}

// Subtract returns a new NProduct object that results from subtraction other from c
func (c *NProduct[S]) Subtract(other *NProduct[S]) *NProduct[S] {
	if c == other {
		return NewCanonicalSet[S](c.dimensions)
	}
	if c.dimensions != other.dimensions {
		return nil
	}
	return &NProduct[S]{dimensions: other.dimensions, product: c.product.Subtract(other.product)}
}

// ContainedIn returns true if c contained in other
func (c *NProduct[S]) ContainedIn(other *NProduct[S]) bool {
	if c.dimensions != other.dimensions {
		log.Panic("dimensionality mismatch")
	}
	return c.product.ContainedIn(other.product)
}

// Copy returns a new NProduct object, copied from c
func (c *NProduct[S]) Copy() *NProduct[S] {
	res := NewCanonicalSet[S](c.dimensions)
	for _, p := range c.product.Partitions() {
		res.product.Insert(p.Key, p.Value.Copy())
	}
	return res
}

// Partitions returns the list of maximal partitions in c, each partition as a slice of NProduct
func (c *NProduct[S]) Partitions() [][]S {
	res := [][]S{}
	if c.dimensions == 1 {
		for _, k := range c.product.Left() {
			res = append(res, []S{k})
		}
		return res
	}
	for _, pair := range c.product.Partitions() {
		subRes := pair.Value.Partitions()
		for _, subList := range subRes {
			partition := []S{pair.Key}
			partition = append(partition, subList...)
			res = append(res, partition)
		}
	}
	return res
}

// Swap returns a new NProduct object, built from the input NProduct object,
// with dimensions dim1 and dim2 swapped
func (c *NProduct[S]) Swap(dim1, dim2 int) *NProduct[S] {
	if c.IsEmpty() || dim1 == dim2 {
		return c.Copy()
	}
	if min(dim1, dim2) < 0 || max(dim1, dim2) >= c.dimensions {
		log.Panicf("invalid dimensions: %d, %d", dim1, dim2)
	}
	res := NewCanonicalSet[S](c.dimensions)
	for _, path := range c.Partitions() {
		if !path[dim1].Equal(path[dim2]) {
			// Shallow clone should be enough, since we do shallow swap:
			path = slices.Clone(path)
			path[dim1], path[dim2] = path[dim2], path[dim1]
		}
		res = res.Union(PartitionN(path))
	}
	return res
}

// PartitionN returns a new NProduct created from a single input partition
// the input partition is a slice of NProduct, treated as ordered list of dimension values
func PartitionN[S Set[S]](partition []S) *NProduct[S] {
	if len(partition) == 0 {
		return nil
	}
	if len(partition) == 1 {
		res := NewCanonicalSet[S](1)
		res.product.Insert(partition[0], NewCanonicalSet[S](0))
		return res
	}
	res := NewCanonicalSet[S](len(partition))
	res.product.Insert(partition[0], PartitionN(partition[1:]))
	return res
}
