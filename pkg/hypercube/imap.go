// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hypercube

type IMap[S Set[S], V Set[V]] struct {
	m *Map[S, V]
}

func NewIMap[S Set[S], V Set[V]]() *IMap[S, V] {
	return &IMap[S, V]{m: NewMap[S, V]()}
}

func (m *IMap[S, V]) Insert(s S, v V) {
	m.m.Insert(s, v)
	m.canonicalize()
}

// Union returns a new IMap object that results from union of m with other
func (m *IMap[S, V]) Union(other *IMap[S, V]) *IMap[S, V] {
	if m == other {
		return m.Copy()
	}
	if m.IsEmpty() {
		return other.Copy()
	}
	if other.IsEmpty() {
		return m.Copy()
	}
	res := NewIMap[S, V]()
	remainingFromOther := NewMap[S, S]()
	for _, k := range other.Keys() {
		remainingFromOther.Insert(k, k.Copy())
	}
	for _, pair := range m.Pairs() {
		remainingFromSelf := pair.Key.Copy()
		for _, otherPair := range other.Pairs() {
			commonElem := pair.Key.Intersect(otherPair.Key)
			if commonElem.IsEmpty() {
				continue
			}
			if v, ok := remainingFromOther.At(otherPair.Key); ok {
				remainingFromOther.Insert(otherPair.Key, v.Subtract(commonElem))
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonElem)
			newSubElem := pair.Value.Union(otherPair.Value)
			res.Insert(commonElem, newSubElem)
		}
		if !remainingFromSelf.IsEmpty() {
			res.Insert(remainingFromSelf, pair.Value.Copy())
		}
	}
	for _, pair := range remainingFromOther.Pairs() {
		if !pair.Value.IsEmpty() {
			if otherValue, ok := other.At(pair.Key); ok {
				res.Insert(pair.Value, otherValue.Copy())
			}
		}
	}
	res.canonicalize()
	return res
}

// Intersect returns a new IMap object that results from intersection of m with other
func (m *IMap[S, V]) Intersect(other *IMap[S, V]) *IMap[S, V] {
	if m == other {
		return m.Copy()
	}
	res := NewIMap[S, V]()
	for _, pair := range m.Pairs() {
		for _, otherPair := range other.Pairs() {
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

// Subtract returns a new IMap object that results from subtraction other from m
func (m *IMap[S, V]) Subtract(other *IMap[S, V]) *IMap[S, V] {
	if m == other {
		return NewIMap[S, V]()
	}
	if other.IsEmpty() {
		return m.Copy()
	}
	res := NewIMap[S, V]()
	for _, pair := range m.Pairs() {
		remainingFromSelf := pair.Key.Copy()
		for _, otherPair := range other.Pairs() {
			commonELem := pair.Key.Intersect(otherPair.Key)
			if commonELem.IsEmpty() {
				continue
			}
			remainingFromSelf = remainingFromSelf.Subtract(commonELem)
			newSubElem := pair.Value.Subtract(otherPair.Value)
			if !newSubElem.IsEmpty() {
				res.Insert(commonELem, newSubElem)
			}
		}
		if !remainingFromSelf.IsEmpty() {
			res.Insert(remainingFromSelf, pair.Value.Copy())
		}
	}
	res.canonicalize()
	return res
}

// ContainedIn returns true if m contained in other
func (m *IMap[S, V]) ContainedIn(other *IMap[S, V]) bool {
	subsetCount := 0
	for _, pair := range m.Pairs() {
		LeftoverKey := pair.Key.Copy()
		for _, otherPair := range other.Pairs() {
			commonKey := otherPair.Key.Intersect(LeftoverKey)
			if commonKey.IsEmpty() {
				continue
			}
			subContainment := pair.Value.ContainedIn(otherPair.Value)
			if !subContainment {
				return false
			}
			LeftoverKey = LeftoverKey.Subtract(commonKey)
			if LeftoverKey.IsEmpty() {
				subsetCount += 1
				break
			}
		}
	}
	return subsetCount == m.Size()
}

func (m *IMap[S, V]) Copy() *IMap[S, V] {
	return &IMap[S, V]{m: m.m.Copy()}
}

func (m *IMap[S, V]) At(s S) (res V, ok bool) {
	return m.m.At(s)
}

func (m *IMap[S, V]) Pairs() []Pair[S, V] {
	return m.m.Pairs()
}

func (m *IMap[S, V]) Keys() []S {
	return m.m.Keys()
}

func (m *IMap[S, V]) Values() []V {
	return m.m.Values()
}

func (m *IMap[S, V]) Equal(other *IMap[S, V]) bool {
	return m.m.Equal(other.m)
}

func (m *IMap[S, V]) IsEmpty() bool {
	return m.m.IsEmpty()
}

func (m *IMap[S, V]) Size() int {
	// maybe inappropriate
	return m.m.Size()
}

func (m *IMap[S, V]) canonicalize() {
	newM := NewMap[S, V]()
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
