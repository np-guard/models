// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ds

// ProductLeft is a cartesian product of two sets
// The implementation represents the sets succinctly, merging keys with equivalent values to a single key
type ProductLeft[K Set[K], V Set[V]] struct {
	m *HashMap[K, V]
}

func NewProductLeft[K Set[K], V Set[V]]() *ProductLeft[K, V] {
	return &ProductLeft[K, V]{m: NewMap[K, V]()}
}

func CartesianPairLeft[K Set[K], V Set[V]](k K, v V) *ProductLeft[K, V] {
	m := NewProductLeft[K, V]()
	if !k.IsEmpty() && !v.IsEmpty() {
		m.m.Insert(k, v)
	}
	return m
}

func (m *ProductLeft[K, V]) Left() []K {
	return m.m.Keys()
}

func (m *ProductLeft[K, V]) Right() []V {
	return m.m.Values()
}

func asLeftProduct[K Set[K], V Set[V]](m Product[K, V]) *ProductLeft[K, V] {
	p, ok := m.(*ProductLeft[K, V])
	if ok {
		return p
	}
	res := NewProductLeft[K, V]()
	for _, pair := range m.Partitions() {
		res.m.Insert(pair.Left, pair.Right)
	}
	res.canonicalize()
	return res
}

func (m *ProductLeft[K, V]) Equal(other Product[K, V]) bool {
	return m.m.Equal(asLeftProduct(other).m)
}

func (m *ProductLeft[K, V]) Copy() Product[K, V] {
	return &ProductLeft[K, V]{m: m.m.Copy()}
}

func (m *ProductLeft[K, V]) Hash() int {
	const rrr = 5
	res := rrr
	for _, p := range m.Partitions() {
		res ^= (p.Left.Hash() << 1) ^ p.Right.Hash()
	}
	return res
}

func (m *ProductLeft[K, V]) IsEmpty() bool {
	return m.m.IsEmpty()
}

func (m *ProductLeft[K, V]) Size() int {
	res := 0
	for _, p := range m.m.Pairs() {
		res += p.Left.Size() * p.Right.Size()
	}
	return res
}

// IsSubset returns true if m contained in other
func (m *ProductLeft[K, V]) IsSubset(other Product[K, V]) bool {
	subsetCount := 0
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Left.Copy()
		for _, otherPair := range other.Partitions() {
			commonKey := otherPair.Left.Intersect(LeftoverKey)
			if commonKey.IsEmpty() {
				continue
			}
			if !pair.Right.IsSubset(otherPair.Right) {
				return false
			}
			LeftoverKey = LeftoverKey.Subtract(commonKey)
			if LeftoverKey.IsEmpty() {
				subsetCount += 1
				break
			}
		}
	}
	return subsetCount == m.m.Size()
}

// Union returns a new Product object that results from union of m with other
func (m *ProductLeft[K, V]) Union(other Product[K, V]) Product[K, V] {
	if m == other {
		return m.Copy()
	}
	if m.IsEmpty() {
		return other.Copy()
	}
	if other.IsEmpty() {
		return m.Copy()
	}
	remainingFromOther := NewMap[K, K]()
	for _, k := range other.Left() {
		remainingFromOther.Insert(k, k)
	}
	res := NewProductLeft[K, V]()
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Left // copy will happen upon insertion
		for _, otherPair := range other.Partitions() {
			commonElem := pair.Left.Intersect(otherPair.Left)
			if commonElem.IsEmpty() {
				continue
			}
			if v, ok := remainingFromOther.At(otherPair.Left); ok {
				remainingFromOther.Insert(otherPair.Left, v.Subtract(commonElem))
			}
			LeftoverKey = LeftoverKey.Subtract(commonElem)
			newSubElem := pair.Right.Union(otherPair.Right)
			res.m.Insert(commonElem, newSubElem)
		}
		if !LeftoverKey.IsEmpty() {
			res.m.Insert(LeftoverKey, pair.Right)
		}
	}
	for _, pair := range remainingFromOther.Pairs() {
		if !pair.Right.IsEmpty() {
			if otherValue, ok := asLeftProduct(other).m.At(pair.Left); ok {
				res.m.Insert(pair.Right, otherValue)
			}
		}
	}
	res.canonicalize()
	return res
}

// Intersect returns a new Product object that results from intersection of m with other
func (m *ProductLeft[K, V]) Intersect(other Product[K, V]) Product[K, V] {
	if m == other {
		return m.Copy()
	}
	res := NewProductLeft[K, V]()
	for _, pair := range m.Partitions() {
		for _, otherPair := range other.Partitions() {
			commonELem := pair.Left.Intersect(otherPair.Left)
			if commonELem.IsEmpty() {
				continue
			}
			newSubElem := pair.Right.Intersect(otherPair.Right)
			if !newSubElem.IsEmpty() {
				res.m.Insert(commonELem, newSubElem)
			}
		}
	}
	res.canonicalize()
	return res
}

// Subtract returns a new Product object that results from subtraction other from m
func (m *ProductLeft[K, V]) Subtract(other Product[K, V]) Product[K, V] {
	if m == other {
		return NewProductLeft[K, V]()
	}
	if other.IsEmpty() {
		return m.Copy()
	}
	res := NewProductLeft[K, V]()
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Left // copy will happen upon insertion
		for _, otherPair := range other.Partitions() {
			commonELem := pair.Left.Intersect(otherPair.Left)
			if commonELem.IsEmpty() {
				continue
			}
			LeftoverKey = LeftoverKey.Subtract(commonELem)
			newSubElem := pair.Right.Subtract(otherPair.Right)
			if !newSubElem.IsEmpty() {
				res.m.Insert(commonELem, newSubElem)
			}
		}
		if !LeftoverKey.IsEmpty() {
			res.m.Insert(LeftoverKey, pair.Right)
		}
	}
	res.canonicalize()
	return res
}

func (m *ProductLeft[K, V]) canonicalize() {
	for _, k := range m.m.Keys() {
		if k.IsEmpty() {
			m.m.Delete(k)
		}
	}
	newM := NewMap[K, V]()
	for _, p := range InverseMap(m.m).Pairs() {
		items := p.Value.Items()
		if len(items) == 0 {
			continue
		}
		newKey := items[0]
		for _, v := range items[1:] {
			newKey = newKey.Union(v)
		}
		newM.Insert(newKey, p.Key)
	}
	m.m = newM
}

func (m *ProductLeft[K, V]) Partitions() []Pair[K, V] {
	return m.m.Pairs()
}

// Swap returns a new Product object, built from the input Product object,
// with left and right swapped
func (m *ProductLeft[K, V]) Swap() Product[V, K] {
	if m.IsEmpty() {
		return NewProductLeft[V, K]()
	}
	var res Product[V, K] = NewProductLeft[V, K]()
	for _, pair := range m.Partitions() {
		res = res.Union(CartesianPairLeft(pair.Right, pair.Left))
	}
	res.(*ProductLeft[V, K]).canonicalize()
	return res
}
