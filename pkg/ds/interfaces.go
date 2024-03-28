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
	Size() int
	IsSubset(Self) bool
	Union(Self) Self
	Intersect(Self) Self
	Subtract(Self) Self
}

// Product is a subset of cartesian product of sets S1 x S2
type Product[S1 Set[S1], S2 Set[S2]] interface {
	Set[Product[S1, S2]]
	Partitions() []Pair[S1, S2]
	Left() []S1
	Right() []S2
	Swap() Product[S2, S1]
}

// TripleSet is a 3-product of sets S1 x S2 x S3
type TripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] interface {
	Set[TripleSet[S1, S2, S3]]
	Partitions() []Triple[S1, S2, S3]
	Swap23() TripleSet[S1, S3, S2]
	Swap12() TripleSet[S2, S1, S3]
	Swap13() TripleSet[S3, S2, S1]
}
