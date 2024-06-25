/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package ds defines and implements a collection of data structures,
// designed to hold cross-products of two Set structures in a succinct and canonical way.
// The most important interface is Product. It is implemented by ProductLeft using injective
// (one-to-one) mapping from sets to sets, where each key-value pair defines a complete cross
// product of the two sets.
package ds

type Sized interface {
	IsEmpty() bool
	// Size returns the actual, full size of the set. For Product, it returns the number of pairs.
	Size() int
}

type Comparable[Self any] interface {
	Equal(Self) bool
	Copy() Self
}

type Hashable[Self any] interface {
	Comparable[Self]
	Hash() int
}

type Set[Self any] interface {
	Hashable[Self]
	Sized
	IsSubset(Self) bool
	Union(Self) Self
	Intersect(Self) Self
	Subtract(Self) Self
}

// Product is a cartesian product of sets S1 x S2
type Product[S1 Set[S1], S2 Set[S2]] interface {
	Set[Product[S1, S2]]
	Partitions() []Pair[S1, S2]
	NumPartitions() int
	Left(empty S1) S1
	Right(empty S2) S2
	Swap() Product[S2, S1]
}

// TripleSet is a 3-product of sets S1 x S2 x S3
type TripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] interface {
	Set[TripleSet[S1, S2, S3]]
	Partitions() []Triple[S1, S2, S3]
}
