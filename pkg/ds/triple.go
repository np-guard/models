/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ds

func (t Triple[S1, S2, S3]) ID() Triple[S1, S2, S3] {
	return Triple[S1, S2, S3]{S1: t.S1, S2: t.S2, S3: t.S3}
}

func (t Triple[S1, S2, S3]) Swap12() Triple[S2, S1, S3] {
	return Triple[S2, S1, S3]{S1: t.S2, S2: t.S1, S3: t.S3}
}

func (t Triple[S1, S2, S3]) Swap23() Triple[S1, S3, S2] {
	return Triple[S1, S3, S2]{S1: t.S1, S2: t.S3, S3: t.S2}
}

func (t Triple[S1, S2, S3]) Swap13() Triple[S3, S2, S1] {
	return Triple[S3, S2, S1]{S1: t.S3, S2: t.S2, S3: t.S1}
}

func (t Triple[S1, S2, S3]) ShiftLeft() Triple[S2, S3, S1] {
	return Triple[S2, S3, S1]{S1: t.S2, S2: t.S3, S3: t.S1}
}

func (t Triple[S1, S2, S3]) ShiftRight() Triple[S3, S1, S2] {
	return Triple[S3, S1, S2]{S1: t.S3, S2: t.S1, S3: t.S2}
}
