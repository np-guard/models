/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ds"
)

// Int is an adapter for ints to be Hashable
type Int struct {
	int
}

func (n Int) Copy() Int {
	return Int{n.int}
}

func (n Int) Equal(other Int) bool {
	return n.int == other.int
}

func (n Int) Hash() int {
	// We want collisions
	return n.int / 2
}

func assertEmpty(t *testing.T, s ds.Sized) {
	t.Helper()
	require.True(t, s.IsEmpty())
	require.True(t, s.Size() == 0)
}

func assertNotEmpty(t *testing.T, s ds.Sized) {
	t.Helper()
	require.False(t, s.IsEmpty())
	require.False(t, s.Size() == 0)
}

func assertNotEqual[T ds.Comparable[T]](t *testing.T, m1, m2 T) {
	t.Helper()
	require.False(t, m1.Equal(m2))
	require.False(t, m2.Equal(m1))
}

func assertEqual[T ds.Comparable[T]](t *testing.T, m1, m2 T) {
	t.Helper()
	require.True(t, m1.Equal(m2))
	require.True(t, m2.Equal(m1))
}
