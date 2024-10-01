/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

type Pair[L, R any] struct {
	Left  L
	Right R
}

type Triple[S1, S2, S3 any] struct {
	S1 S1
	S2 S2
	S3 S3
}
