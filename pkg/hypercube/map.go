// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube

type Comparable[T any] interface {
	Equal(T) bool
}

type Hashable[T any] interface {
	Comparable[T]
	Copy() T
	Hash() int
}

type Pair[K, V any] struct {
	Key   K
	Value V
}

type Map[K Hashable[K], V Comparable[V]] struct {
	m map[int][]Pair[K, V]
}

func NewMap[K Hashable[K], V Comparable[V]]() Map[K, V] {
	return Map[K, V]{m: map[int][]Pair[K, V]{}}
}

func (m *Map[K, V]) Insert(k K, v V) {
	Pairs := m.m[k.Hash()]
	if Pairs == nil {
		m.m[k.Hash()] = []Pair[K, V]{{Key: k.Copy(), Value: v}}
		return
	}
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			Pairs[i].Value = v
			return
		}
	}
	m.m[k.Hash()] = append(Pairs, Pair[K, V]{Key: k.Copy(), Value: v})
}

func (m *Map[K, V]) At(k K) (res V, ok bool) {
	Pairs := m.m[k.Hash()]
	for i := range Pairs {
		if Pairs[i].Key.Equal(k) {
			return Pairs[i].Value, true
		}
	}
	return res, false
}

func (m *Map[K, V]) Pairs() []Pair[K, V] {
	res := []Pair[K, V]{}
	for _, v := range m.m {
		res = append(res, v...)
	}
	return res
}

func (m *Map[K, V]) Keys() []K {
	res := []K{}
	for _, p := range m.Pairs() {
		res = append(res, p.Key)
	}
	return res
}

func (m *Map[K, V]) Values() []V {
	res := []V{}
	for _, p := range m.Pairs() {
		res = append(res, p.Value)
	}
	return res
}

func (m *Map[K, V]) Equal(other *Map[K, V]) bool {
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

func (m *Map[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

func (m *Map[K, V]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}
