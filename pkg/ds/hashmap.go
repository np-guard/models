// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ds

type HashMap[K Hashable[K], V Comparable[V]] struct {
	m map[int][]Pair[K, V]
}

func NewMap[K Hashable[K], V Comparable[V]]() *HashMap[K, V] {
	return &HashMap[K, V]{m: map[int][]Pair[K, V]{}}
}

// Delete k from the map
func (m *HashMap[K, V]) Delete(k K) {
	Pairs := m.m[k.Hash()]
	if Pairs == nil {
		return
	}
	var res []Pair[K, V]
	for i := range Pairs {
		if !Pairs[i].Left.Equal(k) {
			res = append(res, Pairs[i])
		}
	}
	m.m[k.Hash()] = res
}

// Insert a mapping from a copy of k to a copy of v
func (m *HashMap[K, V]) Insert(k K, v V) {
	Pairs := m.m[k.Hash()]
	if Pairs == nil {
		m.m[k.Hash()] = []Pair[K, V]{{Left: k.Copy(), Right: v.Copy()}}
		return
	}
	for i := range Pairs {
		if Pairs[i].Left.Equal(k) {
			Pairs[i].Right = v.Copy()
			return
		}
	}
	m.m[k.Hash()] = append(Pairs, Pair[K, V]{Left: k.Copy(), Right: v.Copy()})
}

func (m *HashMap[K, V]) Copy() *HashMap[K, V] {
	res := NewMap[K, V]()
	for _, p := range m.Pairs() {
		res.Insert(p.Left.Copy(), p.Right.Copy())
	}
	return res
}

func (m *HashMap[K, V]) At(k K) (res V, ok bool) {
	Pairs := m.m[k.Hash()]
	for i := range Pairs {
		if Pairs[i].Left.Equal(k) {
			return Pairs[i].Right, true
		}
	}
	return res, false
}

func (m *HashMap[K, V]) Pairs() []Pair[K, V] {
	var res []Pair[K, V]
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

func (m *HashMap[K, V]) Keys() []K {
	pairs := m.Pairs()
	res := make([]K, len(pairs))
	for i, p := range m.Pairs() {
		res[i] = p.Left
	}
	return res
}

func (m *HashMap[K, V]) Values() []V {
	pairs := m.Pairs()
	res := make([]V, len(pairs))
	for i, p := range pairs {
		res[i] = p.Right
	}
	return res
}

func (m *HashMap[K, V]) Equal(other *HashMap[K, V]) bool {
	me := m.Pairs()
	he := other.Pairs()
outer:
	for _, v := range me {
		for _, q := range he {
			if v.Left.Equal(q.Left) {
				if v.Right.Equal(q.Right) {
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
