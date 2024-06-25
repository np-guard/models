/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// MultiPair is an element of a multimap: a pair of (key, set[Value])
type MultiPair[K Hashable[K], V Hashable[V]] struct {
	Key   K
	Value *HashSet[V]
}

// MultiMap is a mapping from keys to sets of values
type MultiMap[K Hashable[K], V Hashable[V]] struct {
	// m holds a map from the hash to the slice of matching keys
	m map[int][]MultiPair[K, V]
}

// NewMultiMap creates a new empty multimap
func NewMultiMap[K Hashable[K], V Hashable[V]]() *MultiMap[K, V] {
	return &MultiMap[K, V]{m: map[int][]MultiPair[K, V]{}}
}

// Delete k and all its values from the map
func (m *MultiMap[K, V]) Delete(k K) {
	pairs := m.m[k.Hash()]
	if pairs == nil {
		return
	}
	var res []MultiPair[K, V]
	for i := range pairs {
		if !pairs[i].Key.Equal(k) {
			res = append(res, pairs[i])
		}
	}
	if len(res) == 0 {
		delete(m.m, k.Hash())
	} else {
		m.m[k.Hash()] = res
	}
}

// Insert a mapping from a key k to a value v into the multimap
func (m *MultiMap[K, V]) Insert(k K, v V) {
	pairs := m.m[k.Hash()]
	if pairs == nil {
		pairs = []MultiPair[K, V]{}
	} else {
		for i := range pairs {
			if pairs[i].Key.Equal(k) {
				pairs[i].Value.Insert(v)
				return
			}
		}
	}
	m.m[k.Hash()] = append(pairs, MultiPair[K, V]{Key: k.Copy(), Value: NewHashSet[V](v)})
}

// At returns a copy of the set of values for a key k
func (m *MultiMap[K, V]) At(k K) *HashSet[V] {
	Pairs := m.m[k.Hash()]
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			return Pairs[i].Value.Copy()
		}
	}
	return NewHashSet[V]()
}

// MultiPairs returns a slice of (key, set of values) pairs in the multimap. The keys are unique.
func (m *MultiMap[K, V]) MultiPairs() []MultiPair[K, V] {
	var res []MultiPair[K, V]
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

// Pairs returns a slice of (key, value) pairs in the multimap. It is a flattened version of MultiPairs.
func (m *MultiMap[K, V]) Pairs() []Pair[K, V] {
	var res []Pair[K, V]
	for _, mp := range m.MultiPairs() {
		for _, item := range mp.Value.Items() {
			res = append(res, Pair[K, V]{Left: mp.Key, Right: item})
		}
	}
	return res
}

// Keys returns a slice of unique keys in the multimap
func (m *MultiMap[K, V]) Keys() []K {
	pairs := m.MultiPairs()
	res := make([]K, len(pairs))
	for i, p := range pairs {
		res[i] = p.Key
	}
	return res
}

// Values returns a slice of all (possibly repeating) values in the multimap
func (m *MultiMap[K, V]) Values() []V {
	var res []V
	for _, p := range m.MultiPairs() {
		res = append(res, p.Value.Items()...)
	}
	return res
}

// Equal returns true if the multimap is equal to the other multimap. That is, if they have the same set of Pairs().
func (m *MultiMap[K, V]) Equal(other *MultiMap[K, V]) bool {
	if len(m.m) != len(other.m) {
		return false
	}
	if m.IsEmpty() {
		return true
	}
	if m.Size() != other.Size() {
		return false
	}
	for _, v := range m.MultiPairs() {
		if !other.At(v.Key).Equal(v.Value) {
			return false
		}
	}
	return true
}

// Copy returns a deep copy of the multimap
func (m *MultiMap[K, V]) Copy() *MultiMap[K, V] {
	res := NewMultiMap[K, V]()
	for _, p := range m.Pairs() {
		res.Insert(p.Left, p.Right)
	}
	return res
}

// IsEmpty returns true if the multimap is empty
func (m *MultiMap[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

// Size returns the number of key-value pairs in the multimap, that is the length of Pairs()
func (m *MultiMap[K, V]) Size() int {
	return len(m.Pairs())
}

// InverseMap take a HashMap and returns a new multimap where every value in the original map is a key in the new map.
func InverseMap[K Hashable[K], V Hashable[V]](m *HashMap[K, V]) *MultiMap[V, K] {
	inverse := NewMultiMap[V, K]()
	for _, p := range m.Pairs() {
		inverse.Insert(p.Right, p.Left)
	}
	return inverse
}
