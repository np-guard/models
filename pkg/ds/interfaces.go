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

type Comparable[Self any] interface {
	Equal(Self) bool
	Copy() Self
}

type Hashable[Self any] interface {
	Comparable[Self]
	Hash() int
}

type Sized interface {
	IsEmpty() bool

	// Size returns the actual, full size of the set.
	// For Product, it returns the number of pairs of concrete elements that belong to the product, not the number of Partitions().
	// In other words, for Product, p.Size() == sum(s1.Size() * s2.Size() for _, (s1, s2) := range p.Partitions())
	Size() int
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

	// Partitions returns a slice of pairs such that, for (p Product):
	// 	  p.Equal(Union(CartesianPairLeft(s1, s2) for _, (s1, s2) := range p.Partitions())
	// (note that the order is arbitrary; we do not return HashSet because Pair is not Hashable)
	// Partitions returns a slice of pairs in the product set  (note that the order is arbitrary)
	Partitions() []Pair[S1, S2]

	// NumPartitions returns len(Partitions()). It is different from Size() which should return the number of concrete pairs of elements.
	NumPartitions() int

	// Left returns the left projection from pairs in the product set on S1, given input of an empty set in S1.
	Left(empty S1) S1

	// Right returns the right projection from pairs in the product set S2, given input of an empty set in S2.
	Right(empty S2) S2

	// Swap returns a new Product object, built from the input object, with left and right swapped.
	Swap() Product[S2, S1]
}

// TripleSet is a 3-product of sets S1 x S2 x S3
type TripleSet[S1 Set[S1], S2 Set[S2], S3 Set[S3]] interface {
	Set[TripleSet[S1, S2, S3]]
	Partitions() []Triple[S1, S2, S3]
}
