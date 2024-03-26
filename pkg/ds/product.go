// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

// Product is a cartesian product of two sets
// The implementation represents the sets succinctly, merging keys with equivalent values to a single key
type Product[K Set[K], V Set[V]] struct {
	m *HashMap[K, V]
}

func NewProduct[K Set[K], V Set[V]]() *Product[K, V] {
	return &Product[K, V]{m: NewMap[K, V]()}
}

func CartesianPair[K Set[K], V Set[V]](k K, v V) *Product[K, V] {
	m := NewProduct[K, V]()
	m.Insert(k, v)
	return m
}

// Insert mapping from a copy of k to a copy of v
func (m *Product[K, V]) Insert(k K, v V) {
	m.m.Insert(k, v)
	m.canonicalize()
}

// Union returns a new Product object that results from union of m with other
func (m *Product[K, V]) Union(other *Product[K, V]) *Product[K, V] {
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
	res := NewProduct[K, V]()
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Key // copy will happen upon insertion
		for _, otherPair := range other.Partitions() {
			commonElem := pair.Key.Intersect(otherPair.Key)
			if commonElem.IsEmpty() {
				continue
			}
			if v, ok := remainingFromOther.At(otherPair.Key); ok {
				remainingFromOther.Insert(otherPair.Key, v.Subtract(commonElem))
			}
			LeftoverKey = LeftoverKey.Subtract(commonElem)
			newSubElem := pair.Value.Union(otherPair.Value)
			res.Insert(commonElem, newSubElem)
		}
		if !LeftoverKey.IsEmpty() {
			res.Insert(LeftoverKey, pair.Value)
		}
	}
	for _, pair := range remainingFromOther.Pairs() {
		if !pair.Value.IsEmpty() {
			if otherValue, ok := other.At(pair.Key); ok {
				res.Insert(pair.Value, otherValue)
			}
		}
	}
	res.canonicalize()
	return res
}

// Intersect returns a new Product object that results from intersection of m with other
func (m *Product[K, V]) Intersect(other *Product[K, V]) *Product[K, V] {
	if m == other {
		return m.Copy()
	}
	res := NewProduct[K, V]()
	for _, pair := range m.Partitions() {
		for _, otherPair := range other.Partitions() {
			commonELem := pair.Key.Intersect(otherPair.Key)
			if commonELem.IsEmpty() {
				continue
			}
			newSubElem := pair.Value.Intersect(otherPair.Value)
			if !newSubElem.IsEmpty() {
				res.Insert(commonELem, newSubElem)
			}
		}
	}
	res.canonicalize()
	return res
}

// Subtract returns a new Product object that results from subtraction other from m
func (m *Product[K, V]) Subtract(other *Product[K, V]) *Product[K, V] {
	if m == other {
		return NewProduct[K, V]()
	}
	if other.IsEmpty() {
		return m.Copy()
	}
	res := NewProduct[K, V]()
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Key // copy will happen upon insertion
		for _, otherPair := range other.Partitions() {
			commonELem := pair.Key.Intersect(otherPair.Key)
			if commonELem.IsEmpty() {
				continue
			}
			LeftoverKey = LeftoverKey.Subtract(commonELem)
			newSubElem := pair.Value.Subtract(otherPair.Value)
			if !newSubElem.IsEmpty() {
				res.Insert(commonELem, newSubElem)
			}
		}
		if !LeftoverKey.IsEmpty() {
			res.Insert(LeftoverKey, pair.Value)
		}
	}
	res.canonicalize()
	return res
}

// ContainedIn returns true if m contained in other
func (m *Product[K, V]) ContainedIn(other *Product[K, V]) bool {
	subsetCount := 0
	for _, pair := range m.Partitions() {
		LeftoverKey := pair.Key.Copy()
		for _, otherPair := range other.Partitions() {
			commonKey := otherPair.Key.Intersect(LeftoverKey)
			if commonKey.IsEmpty() {
				continue
			}
			if !pair.Value.ContainedIn(otherPair.Value) {
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

func (m *Product[K, V]) Copy() *Product[K, V] {
	return &Product[K, V]{m: m.m.Copy()}
}

func (m *Product[K, V]) At(k K) (res V, ok bool) {
	return m.m.At(k)
}

func (m *Product[K, V]) Partitions() []Pair[K, V] {
	return m.m.Pairs()
}

func (m *Product[K, V]) Left() []K {
	return m.m.Keys()
}

func (m *Product[K, V]) Right() []V {
	return m.m.Values()
}

func (m *Product[K, V]) Equal(other *Product[K, V]) bool {
	return m.m.Equal(other.m)
}

const rrr = 5

func (c *Product[K, V]) Hash() int {
	res := rrr
	for _, p := range c.Partitions() {
		res ^= (p.Key.Hash() << 1) ^ p.Value.Hash()
	}
	return res
}

func (m *Product[K, V]) IsEmpty() bool {
	return m.m.IsEmpty()
}

func (m *Product[K, V]) Size() int {
	res := 0
	for _, p := range m.m.Pairs() {
		res += p.Key.Size() * p.Value.Size()
	}
	return res
}

func (m *Product[K, V]) canonicalize() {
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

// Swap returns a new Product object, built from the input Product object,
// with left and right swapped
func (m *Product[K, V]) Swap() *Product[V, K] {
	if m.IsEmpty() {
		return NewProduct[V, K]()
	}
	res := NewProduct[V, K]()
	for _, pair := range m.Partitions() {
		res.Insert(pair.Value, pair.Key)
	}
	return res
}
