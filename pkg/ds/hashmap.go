/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// HashMap is a generic hash map with keys of a Hashable type K and values of Comparable type V.
type HashMap[K Hashable[K], V Comparable[V]] struct {
	m map[int][]Pair[K, V] // map from key hash value to the list of pairs (k,v)
}

// NewHashMap creates a new empty hash map.
func NewHashMap[K Hashable[K], V Comparable[V]]() *HashMap[K, V] {
	return &HashMap[K, V]{m: map[int][]Pair[K, V]{}}
}

// Delete a key and its value from the map, if it exists.
func (m *HashMap[K, V]) Delete(k K) {
	pairs := m.m[k.Hash()]
	if pairs == nil {
		return
	}
	var res []Pair[K, V]
	for i := range pairs {
		if !pairs[i].Left.Equal(k) {
			res = append(res, pairs[i])
		}
	}
	if len(res) == 0 {
		delete(m.m, k.Hash())
	} else {
		m.m[k.Hash()] = res
	}
}

// Insert a mapping from a copy of a key k to a copy of a value v into the map.
func (m *HashMap[K, V]) Insert(k K, v V) {
	pairs := m.m[k.Hash()]
	if pairs == nil {
		m.m[k.Hash()] = []Pair[K, V]{{Left: k.Copy(), Right: v.Copy()}}
		return
	}
	for i := range pairs {
		if pairs[i].Left.Equal(k) {
			pairs[i].Right = v.Copy()
			return
		}
	}
	m.m[k.Hash()] = append(pairs, Pair[K, V]{Left: k.Copy(), Right: v.Copy()})
}

// Copy returns a deep copy of the map.
func (m *HashMap[K, V]) Copy() *HashMap[K, V] {
	res := NewHashMap[K, V]()
	for _, p := range m.Pairs() {
		res.Insert(p.Left, p.Right)
	}
	return res
}

// At returns a pair of a value and a boolean indicating whether the key exists in the map.
// The value is only valid if the boolean is true.
func (m *HashMap[K, V]) At(k K) (res V, ok bool) {
	pairs := m.m[k.Hash()]
	for i := range pairs {
		if pairs[i].Left.Equal(k) {
			return pairs[i].Right, true
		}
	}
	return res, false
}

// Pairs returns a slice of all key-value pairs in the map.
func (m *HashMap[K, V]) Pairs() []Pair[K, V] {
	var res []Pair[K, V]
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

// Keys returns a slice of all keys in the map.
func (m *HashMap[K, V]) Keys() []K {
	pairs := m.Pairs()
	res := make([]K, len(pairs))
	for i, p := range m.Pairs() {
		res[i] = p.Left
	}
	return res
}

// Values returns a slice of all (not necessarily unique) values in the map.
func (m *HashMap[K, V]) Values() []V {
	pairs := m.Pairs()
	res := make([]V, len(pairs))
	for i, p := range pairs {
		res[i] = p.Right
	}
	return res
}

// Equal returns true if the map holds the same key-value pairs as the other map.
func (m *HashMap[K, V]) Equal(other *HashMap[K, V]) bool {
	if len(m.m) != len(other.m) {
		return false
	}
	if m.IsEmpty() {
		return true
	}
	if m.Size() != other.Size() {
		return false
	}
	for _, k := range m.Keys() {
		v1, ok := m.At(k)
		if !ok {
			panic("Impossible: key not found")
		}
		v2, ok := other.At(k)
		if !ok || !v1.Equal(v2) {
			return false
		}
	}
	return true
}

// IsEmpty returns true if the map is empty.
func (m *HashMap[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

// Size returns the number of key-value pairs in the map.
func (m *HashMap[K, V]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}
