// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

type HashMap[K Hashable[K], V Comparable[V]] struct {
	m map[int][]Pair[K, V]
}

func NewMap[K Hashable[K], V Comparable[V]]() *HashMap[K, V] {
	return &HashMap[K, V]{m: map[int][]Pair[K, V]{}}
}

// Insert mapping from a copy of k to a copy of v
func (m *HashMap[K, V]) Insert(k K, v V) {
	Pairs := m.m[k.Hash()]
	if Pairs == nil {
		m.m[k.Hash()] = []Pair[K, V]{{Key: k.Copy(), Value: v.Copy()}}
		return
	}
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			Pairs[i].Value = v.Copy()
			return
		}
	}
	m.m[k.Hash()] = append(Pairs, Pair[K, V]{Key: k.Copy(), Value: v.Copy()})
}

func (m *HashMap[K, V]) Copy() *HashMap[K, V] {
	res := NewMap[K, V]()
	for _, p := range m.Pairs() {
		res.Insert(p.Key.Copy(), p.Value.Copy())
	}
	return res
}

func (m *HashMap[K, V]) At(k K) (res V, ok bool) {
	Pairs := m.m[k.Hash()]
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			return Pairs[i].Value, true
		}
	}
	return res, false
}

func (m *HashMap[K, V]) Pairs() []Pair[K, V] {
	res := []Pair[K, V]{}
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

func (m *HashMap[K, V]) Keys() []K {
	res := []K{}
	for _, p := range m.Pairs() {
		res = append(res, p.Key)
	}
	return res
}

func (m *HashMap[K, V]) Values() []V {
	res := []V{}
	for _, p := range m.Pairs() {
		res = append(res, p.Value)
	}
	return res
}

func (m *HashMap[K, V]) Equal(other *HashMap[K, V]) bool {
	me := m.Pairs()
	he := other.Pairs()
outer:
	for _, v := range me {
		for _, q := range he {
			if v.Key.Equal(q.Key) {
				if v.Value.Equal(q.Value) {
					continue outer
				}
				break
			}
		}
		return false
	}
	return true
}

func (m *HashMap[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

func (m *HashMap[K, V]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}
