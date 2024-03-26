// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

// MultiPair is an element of a multimap: a pair of (key, set[Value])
type MultiPair[K Hashable[K], V Hashable[V]] struct {
	Key   K
	Value HashSet[V]
}

// MultiMap is a mapping from keys to sets of values
type MultiMap[K Hashable[K], V Hashable[V]] struct {
	m map[int][]MultiPair[K, V]
}

func NewMultiMap[K Hashable[K], V Hashable[V]]() MultiMap[K, V] {
	return MultiMap[K, V]{m: map[int][]MultiPair[K, V]{}}
}

func (m *MultiMap[K, V]) Insert(k K, v V) {
	vs := NewHashSet[V]()
	vs.Insert(v)
	pairs := m.m[k.Hash()]
	if pairs == nil {
		m.m[k.Hash()] = []MultiPair[K, V]{{Key: k.Copy(), Value: vs}}
		return
	}
	for i := range pairs {
		if pairs[i].Key.Equal(k) {
			pairs[i].Value.Insert(v)
			return
		}
	}
	m.m[k.Hash()] = append(pairs, MultiPair[K, V]{Key: k.Copy(), Value: vs})
}

func (m *MultiMap[K, V]) At(k K) HashSet[V] {
	Pairs := m.m[k.Hash()]
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			return Pairs[i].Value
		}
	}
	return NewHashSet[V]()
}

func (m *MultiMap[K, V]) Pairs() []MultiPair[K, V] {
	res := []MultiPair[K, V]{}
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

func (m *MultiMap[K, V]) Keys() []K {
	res := []K{}
	for _, p := range m.Pairs() {
		res = append(res, p.Key)
	}
	return res
}

func (m *MultiMap[K, V]) Values() []V {
	res := []V{}
	for _, p := range m.Pairs() {
		res = append(res, p.Value.Items()...)
	}
	return res
}

func (m *MultiMap[K, V]) Equal(other *MultiMap[K, V]) bool {
	me := m.Pairs()
	he := other.Pairs()
outer:
	for _, v := range me {
		for i := range he {
			if v.Key.Equal(he[i].Key) {
				if v.Value.Equal(&he[i].Value) {
					continue outer
				}
				break
			}
		}
		return false
	}
	return true
}

func (m *MultiMap[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

func (m *MultiMap[K, V]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}

func InverseMap[K Hashable[K], V Hashable[V]](m *HashMap[K, V]) *MultiMap[V, K] {
	inverse := NewMultiMap[V, K]()
	for _, p := range m.Pairs() {
		inverse.Insert(p.Value, p.Key)
	}
	return &inverse
}
