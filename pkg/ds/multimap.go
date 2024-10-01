/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// MultiMap is a mapping from keys to sets of values
type MultiMap[K Hashable[K], V Hashable[V]] struct {
	// m holds a map from the hash to the slice of matching keys
	m *HashMap[K, *HashSet[V]]
}

// NewMultiMap creates a new empty multimap
func NewMultiMap[K Hashable[K], V Hashable[V]]() *MultiMap[K, V] {
	return &MultiMap[K, V]{m: NewHashMap[K, *HashSet[V]]()}
}

// Delete k and all its values from the map
func (m *MultiMap[K, V]) Delete(k K) {
	m.m.Delete(k)
}

// Insert a mapping from a key k to a value v into the multimap
func (m *MultiMap[K, V]) Insert(k K, v V) {
	vs, ok := m.m.At(k)
	if !ok {
		m.m.Insert(k, NewHashSet[V](v))
		return
	}
	// vs is not a copy, so we can simply insert v
	vs.Insert(v)
}

// At returns a copy of the set of values for a key k
func (m *MultiMap[K, V]) At(k K) *HashSet[V] {
	vs, ok := m.m.At(k)
	if !ok {
		return NewHashSet[V]()
	}
	return vs.Copy()
}

// MultiPairs returns a slice of (key, set of values) pairs in the multimap. The keys are unique.
func (m *MultiMap[K, V]) MultiPairs() []Pair[K, *HashSet[V]] {
	var res = make([]Pair[K, *HashSet[V]], m.m.Size())
	for i, v := range m.m.Pairs() {
		res[i] = Pair[K, *HashSet[V]]{Left: v.Left, Right: v.Right.Copy()}
	}
	return res
}

// Pairs returns a slice of (key, value) pairs in the multimap. It is a flattened version of MultiPairs.
func (m *MultiMap[K, V]) Pairs() []Pair[K, V] {
	var res []Pair[K, V]
	for _, mp := range m.MultiPairs() {
		for _, item := range mp.Right.Items() {
			res = append(res, Pair[K, V]{Left: mp.Left, Right: item})
		}
	}
	return res
}

// Keys returns a slice of unique keys in the multimap
func (m *MultiMap[K, V]) Keys() []K {
	return m.m.Keys()
}

// Values returns a slice of all (possibly repeating) values in the multimap
func (m *MultiMap[K, V]) Values() []V {
	var res []V
	for _, p := range m.MultiPairs() {
		res = append(res, p.Right.Items()...)
	}
	return res
}

// Equal returns true if the multimap is equal to the other multimap. That is, if they have the same set of Pairs().
func (m *MultiMap[K, V]) Equal(other *MultiMap[K, V]) bool {
	return m.m.Equal(other.m)
}

// Copy returns a deep copy of the multimap
func (m *MultiMap[K, V]) Copy() *MultiMap[K, V] {
	return &MultiMap[K, V]{m: m.m.Copy()}
}

// IsEmpty returns true if the multimap is empty
func (m *MultiMap[K, V]) IsEmpty() bool {
	return m.m.IsEmpty()
}

// Size returns the number of key-value pairs in the multimap, that is the length of Pairs()
func (m *MultiMap[K, V]) Size() int {
	res := 0
	for _, p := range m.MultiPairs() {
		res += p.Right.Size()
	}
	return res
}

// InverseMap take a HashMap and returns a new multimap where every value in the original map is a key in the new map.
func InverseMap[K Hashable[K], V Hashable[V]](m *HashMap[K, V]) *MultiMap[V, K] {
	inverse := NewMultiMap[V, K]()
	for _, p := range m.Pairs() {
		inverse.Insert(p.Right, p.Left)
	}
	return inverse
}
