// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ds

// ProductRight is a cartesian product of two sets
// The implementation simply holds the dual ProductLeft
type ProductRight[K Set[K], V Set[V]] struct {
	p *ProductLeft[V, K]
}

func NewProductRight[K Set[K], V Set[V]]() *ProductRight[K, V] {
	return &ProductRight[K, V]{p: NewProductLeft[V, K]()}
}

func CartesianPairRight[K Set[K], V Set[V]](k K, v V) *ProductRight[K, V] {
	return &ProductRight[K, V]{p: CartesianPairLeft[V, K](v, k)}
}

func (m *ProductRight[K, V]) Left() []K {
	return m.p.Right()
}

func (m *ProductRight[K, V]) Right() []V {
	return m.p.Left()
}

func asRightProduct[K Set[K], V Set[V]](m Product[K, V]) *ProductRight[K, V] {
	res, ok := m.(*ProductRight[K, V])
	if ok {
		return res
	}
	var p Product[V, K] = NewProductLeft[V, K]()
	for _, pair := range m.Partitions() {
		p = p.Union(CartesianPairLeft(pair.Right, pair.Left))
	}
	return &ProductRight[K, V]{p: p.(*ProductLeft[V, K])}
}

func (m *ProductRight[K, V]) Equal(other Product[K, V]) bool {
	return m.p.Equal(asRightProduct(other).p)
}

func (m *ProductRight[K, V]) Copy() Product[K, V] {
	return &ProductRight[K, V]{p: asRightProduct(m).p.Copy().(*ProductLeft[V, K])}
}

func (m *ProductRight[K, V]) Hash() int {
	return m.p.Hash()
}

func (m *ProductRight[K, V]) IsEmpty() bool {
	return m.p.IsEmpty()
}

func (m *ProductRight[K, V]) Size() int {
	return m.p.Size()
}

// IsSubset returns true if m contained in other
func (m *ProductRight[K, V]) IsSubset(other Product[K, V]) bool {
	return m.p.IsSubset(asRightProduct(other).p)
}

// Union returns a new Product object that results from union of m with other
func (m *ProductRight[K, V]) Union(other Product[K, V]) Product[K, V] {
	return &ProductRight[K, V]{p: m.p.Union(asRightProduct(other).p).(*ProductLeft[V, K])}
}

// Intersect returns a new Product object that results from intersection of m with other
func (m *ProductRight[K, V]) Intersect(other Product[K, V]) Product[K, V] {
	return &ProductRight[K, V]{p: m.p.Intersect(asRightProduct(other).p).(*ProductLeft[V, K])}
}

// Subtract returns a new Product object that results from subtraction other from m
func (m *ProductRight[K, V]) Subtract(other Product[K, V]) Product[K, V] {
	return &ProductRight[K, V]{p: m.p.Subtract(asRightProduct(other).p).(*ProductLeft[V, K])}
}

func (m *ProductRight[K, V]) Partitions() []Pair[K, V] {
	partitions := m.p.Partitions()
	res := make([]Pair[K, V], len(partitions))
	for i, pair := range partitions {
		res[i] = pair.Swap()
	}
	return res
}

// Swap returns a new Product object, built from the input Product object,
// with left and right swapped
func (m *ProductRight[K, V]) Swap() Product[V, K] {
	return &ProductRight[V, K]{p: m.p.Swap().(*ProductLeft[K, V])}
}
