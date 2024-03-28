// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ds

type HashSet[K Hashable[K]] struct {
	m map[int][]K
}

func NewHashSet[K Hashable[K]]() HashSet[K] {
	return HashSet[K]{m: map[int][]K{}}
}

func (m *HashSet[K]) Insert(k K) {
	Keys := m.m[k.Hash()]
	if Keys == nil {
		m.m[k.Hash()] = []K{k.Copy()}
		return
	}
	for i := range Keys {
		if Keys[i].Equal(k) {
			return
		}
	}
	m.m[k.Hash()] = append(Keys, k.Copy())
}

func (m *HashSet[K]) Items() []K {
	var res []K
	for _, k := range m.m {
		res = append(res, k...)
	}
	return res
}

func (m *HashSet[K]) Equal(other *HashSet[K]) bool {
	me := m.Items()
	he := other.Items()
outer:
	for _, v := range me {
		for _, q := range he {
			if v.Equal(q) {
				continue outer
			}
		}
		return false
	}
	return true
}

func (m *HashSet[K]) IsEmpty() bool {
	return len(m.m) == 0
}

func (m *HashSet[K]) Size() int {
	res := 0
	for _, v := range m.m {
		res += len(v)
	}
	return res
}
