/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package ds_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
)

type Map = ds.HashMap[Int, Int]

func less(p1, p2 ds.Pair[Int, Int]) bool {
	if p1.Left.int == p2.Left.int {
		return p1.Right.int <= p2.Right.int
	}
	return p1.Left.int <= p2.Left.int
}

func assertMapEmpty(t *testing.T, m *Map) {
	t.Helper()
	assertEmpty(t, m)
	require.True(t, m.Equal(ds.NewHashMap[Int, Int]()))
	_, ok := m.At(Int{1})
	require.False(t, ok)

	require.Len(t, m.Pairs(), 0)
	require.Len(t, m.Keys(), 0)
	require.Len(t, m.Values(), 0)
}

func assertMapSingle(t *testing.T, m *Map, key, value int) {
	t.Helper()
	var v Int
	var ok bool
	nonexistentKey := key - 1
	assertNotEmpty(t, m)
	require.True(t, m.Size() == 1)
	v, ok = m.At(Int{key})
	require.True(t, ok)
	require.True(t, v.int == value)
	_, ok = m.At(Int{key - 1})
	require.False(t, ok)

	{
		pairs := m.Pairs()
		require.Len(t, pairs, 1)
		require.True(t, pairs[0].Left.int == key)
		require.True(t, pairs[0].Right.int == value)
	}
	{
		keys := m.Keys()
		require.Len(t, keys, 1)
		require.True(t, keys[0].int == key)
	}
	{
		values := m.Values()
		require.Len(t, values, 1)
		require.True(t, values[0].int == value)
	}

	m1 := ds.NewHashMap[Int, Int]()
	assertNotEqual(t, m, m1)

	m1.Insert(Int{key}, Int{value})
	assertEqual(t, m, m1)

	m1.Insert(Int{nonexistentKey}, Int{key})
	assertNotEqual(t, m, m1)

	m1.Delete(Int{nonexistentKey})
	assertEqual(t, m, m1)
	m1.Delete(Int{key})
	assertMapEmpty(t, m1)
}

func assertMapDouble(t *testing.T, m *Map, key1, value1, key2, value2 int) {
	t.Helper()
	var v Int
	var ok bool
	nonexistentKey := min(key1, key2) - 1

	require.False(t, m.IsEmpty())
	require.True(t, m.Size() == 2)

	v, ok = m.At(Int{key1})
	require.True(t, ok)
	require.True(t, v.int == value1)

	v, ok = m.At(Int{key2})
	require.True(t, ok)
	require.True(t, v.int == value2)

	_, ok = m.At(Int{nonexistentKey})
	require.False(t, ok)

	{
		pairs := m.Pairs()
		require.Len(t, pairs, 2)
		sort.Slice(pairs, func(i, j int) bool { return less(pairs[i], pairs[j]) })

		expectedPairs := []ds.Pair[Int, Int]{
			{Left: Int{key1}, Right: Int{value1}},
			{Left: Int{key2}, Right: Int{value2}},
		}
		sort.Slice(expectedPairs, func(i, j int) bool { return less(expectedPairs[i], expectedPairs[j]) })

		require.True(t, pairs[0].Left.Equal(expectedPairs[0].Left))
		require.True(t, pairs[0].Right.Equal(expectedPairs[0].Right))
		require.True(t, pairs[1].Left.Equal(expectedPairs[1].Left))
		require.True(t, pairs[1].Right.Equal(expectedPairs[1].Right))
	}

	{
		keys := m.Keys()
		sort.Slice(keys, func(i, j int) bool { return keys[i].int <= keys[j].int })
		k1, k2 := min(key1, key2), max(key2, key1)
		require.Len(t, keys, 2)
		require.True(t, keys[0].int == k1)
		require.True(t, keys[1].int == k2)
	}

	{
		values := m.Values()
		sort.Slice(values, func(i, j int) bool { return values[i].int <= values[j].int })
		v1, v2 := min(value1, value2), max(value1, value2)
		require.Len(t, values, 2)
		require.True(t, values[0].int == v1)
		require.True(t, values[1].int == v2)
	}

	m1 := ds.NewHashMap[Int, Int]()
	assertNotEqual(t, m, m1)

	m1.Insert(Int{key2}, Int{value2})
	assertNotEqual(t, m, m1)

	m1.Insert(Int{key1}, Int{value1})
	assertEqual(t, m, m1)

	m1.Insert(Int{nonexistentKey}, Int{3})
	assertNotEqual(t, m, m1)

	m1.Delete(Int{nonexistentKey})
	assertEqual(t, m, m1)
	m1.Delete(Int{key1})
	m1.Delete(Int{key2})
	assertMapEmpty(t, m1)
}

func TestMap(t *testing.T) {
	var m, dupl *Map

	m = ds.NewHashMap[Int, Int]()
	assertMapEmpty(t, m)
	m.Delete(Int{1})
	assertMapEmpty(t, m)

	dupl = m.Copy()
	require.False(t, dupl == m)
	assertMapEmpty(t, dupl)

	m.Insert(Int{1}, Int{3})
	assertMapSingle(t, m, 1, 3)
	m.Insert(Int{1}, Int{3})
	assertMapSingle(t, m, 1, 3)
	m.Insert(Int{1}, Int{2})
	assertMapSingle(t, m, 1, 2)
	m.Delete(Int{0})
	assertMapSingle(t, m, 1, 2)

	assertNotEqual(t, m, dupl)
	assertMapEmpty(t, dupl)
	dupl = m.Copy()
	assertEqual(t, m, dupl)
	assertMapSingle(t, dupl, 1, 2)

	m.Insert(Int{0}, Int{2})
	assertMapDouble(t, m, 1, 2, 0, 2)
	m.Insert(Int{0}, Int{3})
	assertMapDouble(t, m, 1, 2, 0, 3)
	m.Insert(Int{0}, Int{3})
	assertMapDouble(t, m, 1, 2, 0, 3)
	m.Delete(Int{2})
	assertMapDouble(t, m, 1, 2, 0, 3)

	assertNotEqual(t, m, dupl)
	assertMapSingle(t, dupl, 1, 2)
	dupl = m.Copy()
	assertEqual(t, m, dupl)
	assertMapDouble(t, dupl, 1, 2, 0, 3)

	m.Delete(Int{1})
	assertMapSingle(t, m, 0, 3)
	m.Delete(Int{2})
	assertMapSingle(t, m, 0, 3)

	assertNotEqual(t, m, dupl)
	assertMapDouble(t, dupl, 1, 2, 0, 3)
	dupl = m.Copy()
	assertEqual(t, m, dupl)
	assertMapSingle(t, dupl, 0, 3)

	m.Insert(Int{2}, Int{5})
	assertMapDouble(t, m, 0, 3, 2, 5)
	m.Insert(Int{2}, Int{4})
	assertMapDouble(t, m, 0, 3, 2, 4)
	m.Insert(Int{2}, Int{4})
	assertMapDouble(t, m, 0, 3, 2, 4)
	m.Delete(Int{1})
	assertMapDouble(t, m, 0, 3, 2, 4)

	assertNotEqual(t, m, dupl)
	assertMapSingle(t, dupl, 0, 3)
	dupl = m.Copy()
	assertEqual(t, m, dupl)
	assertMapDouble(t, dupl, 0, 3, 2, 4)
}
