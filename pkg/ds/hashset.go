/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

// HashSet is a generic set with elements of a Hashable type V.
type HashSet[V Hashable[V]] struct {
	m map[int][]V
}

// NewHashSet creates a new hash set with the given elements. Repeated elements are ignored.
func NewHashSet[V Hashable[V]](items ...V) *HashSet[V] {
	res := &HashSet[V]{m: map[int][]V{}}
	for _, item := range items {
		res.Insert(item)
	}
	return res
}

// Insert a value into the set, if it does not already exist.
func (m *HashSet[V]) Insert(v V) {
	Keys := m.m[v.Hash()]
	if Keys == nil {
		m.m[v.Hash()] = []V{v.Copy()}
		return
	}
	for i := range Keys {
		if Keys[i].Equal(v) {
			return
		}
	}
	m.m[v.Hash()] = append(Keys, v.Copy())
}

// Delete a value from the set, if it exists.
func (m *HashSet[V]) Delete(v V) {
	keys := m.m[v.Hash()]
	if keys == nil {
		return
	}
	var res []V
	for i := range keys {
		if !keys[i].Equal(v) {
			res = append(res, keys[i])
		}
	}
	if len(res) == 0 {
		delete(m.m, v.Hash())
	} else {
		m.m[v.Hash()] = res
	}
}

// Contains returns true if the value exists in the set.
func (m *HashSet[V]) Contains(v V) bool {
	keys := m.m[v.Hash()]
	for i := range keys {
		if keys[i].Equal(v) {
			return true
		}
	}
	return false
}

// Items returns a slice of all (non-repeated) values in the set.
func (m *HashSet[V]) Items() []V {
	var res []V
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

// Copy returns a deep copy of the set.
func (m *HashSet[V]) Copy() *HashSet[V] {
	res := NewHashSet[V]()
	for _, p := range m.Items() {
		res.Insert(p)
	}
	return res
}

// Equal returns true if the set contains exactly the same elements as the other set.
func (m *HashSet[V]) Equal(other *HashSet[V]) bool {
	if len(m.m) != len(other.m) {
		return false
	}
	if m.IsEmpty() {
		return true
	}
	if m.Size() != other.Size() {
		return false
	}
	for _, v := range m.Items() {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

// IsEmpty returns true if the set is empty.
func (m *HashSet[V]) IsEmpty() bool {
	return len(m.m) == 0
}

// Size returns the number of elements in the set.
func (m *HashSet[V]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}
