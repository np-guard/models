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

type Set = ds.HashSet[Int]

func assertSetEmpty(t *testing.T, s *Set) {
	t.Helper()
	assertEmpty(t, s)
	require.True(t, s.Equal(ds.NewHashSet[Int]()))
	require.False(t, s.Contains(Int{1}))

	require.Len(t, s.Items(), 0)
}

func assertSingleSet(t *testing.T, s *Set, value int) {
	t.Helper()
	nonexistentValue := value - 1
	assertNotEmpty(t, s)
	require.True(t, s.Size() == 1)
	require.True(t, s.Contains(Int{value}))
	require.False(t, s.Contains(Int{value - 1}))

	{
		items := s.Items()
		require.Len(t, items, 1)
		require.True(t, items[0].int == value)
	}

	m1 := ds.NewHashSet[Int]()
	assertNotEqual(t, s, m1)

	m1.Insert(Int{value})
	assertEqual(t, s, m1)

	m1.Insert(Int{nonexistentValue})
	assertNotEqual(t, s, m1)

	m1.Delete(Int{nonexistentValue})
	assertEqual(t, s, m1)
	m1.Delete(Int{value})
	assertSetEmpty(t, m1)
}

func assertDoubleSet(t *testing.T, s *Set, value1, value2 int) {
	t.Helper()
	nonexistentValue := min(value1, value2) - 1

	require.False(t, s.IsEmpty())
	require.True(t, s.Size() == 2)

	require.True(t, s.Contains(Int{value1}))
	require.True(t, s.Contains(Int{value2}))
	require.False(t, s.Contains(Int{nonexistentValue}))

	{
		values := s.Items()
		sort.Slice(values, func(i, j int) bool { return values[i].int <= values[j].int })
		v1, v2 := min(value1, value2), max(value1, value2)
		require.Len(t, values, 2)
		require.True(t, values[0].int == v1)
		require.True(t, values[1].int == v2)
	}

	m1 := ds.NewHashSet[Int]()
	assertNotEqual(t, s, m1)

	m1.Insert(Int{value2})
	assertNotEqual(t, s, m1)

	m1.Insert(Int{value1})
	assertEqual(t, s, m1)

	m1.Insert(Int{nonexistentValue})
	assertNotEqual(t, s, m1)

	m1.Delete(Int{nonexistentValue})
	assertEqual(t, s, m1)
	m1.Delete(Int{value1})
	m1.Delete(Int{value2})
	assertSetEmpty(t, m1)
}

func TestSet(t *testing.T) {
	var s, dupl *Set

	s = ds.NewHashSet[Int]()
	assertSetEmpty(t, s)
	s.Delete(Int{1})
	assertSetEmpty(t, s)

	dupl = s.Copy()
	require.False(t, dupl == s)
	assertSetEmpty(t, dupl)

	s.Insert(Int{1})
	assertSingleSet(t, s, 1)
	s.Insert(Int{1})
	assertSingleSet(t, s, 1)
	s.Delete(Int{0})
	assertSingleSet(t, s, 1)

	assertNotEqual(t, s, dupl)
	assertSetEmpty(t, dupl)
	dupl = s.Copy()
	assertEqual(t, s, dupl)
	assertSingleSet(t, dupl, 1)

	s.Insert(Int{2})
	assertDoubleSet(t, s, 1, 2)
	s.Insert(Int{1})
	assertDoubleSet(t, s, 2, 1)
	s.Delete(Int{3})
	assertDoubleSet(t, s, 1, 2)

	assertNotEqual(t, s, dupl)
	assertSingleSet(t, dupl, 1)
	dupl = s.Copy()
	assertEqual(t, s, dupl)
	assertDoubleSet(t, dupl, 1, 2)

	s.Delete(Int{1})
	assertSingleSet(t, s, 2)
	s.Delete(Int{3})
	assertSingleSet(t, s, 2)

	assertNotEqual(t, s, dupl)
	assertDoubleSet(t, dupl, 1, 2)
	dupl = s.Copy()
	assertEqual(t, s, dupl)
	assertSingleSet(t, dupl, 2)

	s.Delete(Int{2})
	assertSetEmpty(t, s)

	assertNotEqual(t, s, dupl)
	assertSingleSet(t, dupl, 2)
	dupl = s.Copy()
	assertEqual(t, s, dupl)
	assertSetEmpty(t, s)
}
