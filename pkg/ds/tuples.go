// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ds

type Pair[K, V any] struct {
	Key   K
	Value V
}

type Triple[S1 Set[S1], S2 Set[S2], S3 Set[S3]] struct {
	S1 S1
	S2 S2
	S3 S3
}
