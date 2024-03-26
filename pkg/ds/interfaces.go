// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

type Comparable[T any] interface {
	Equal(T) bool
	Copy() T
}

type Hashable[T any] interface {
	Comparable[T]
	Hash() int
}

type Set[Self any] interface {
	Hashable[Self]
	IsEmpty() bool
	ContainedIn(Self) bool
	Intersect(Self) Self
	Union(Self) Self
	Subtract(Self) Self
	Size() int
}
