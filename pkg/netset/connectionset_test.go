/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/netset"
)

func TestConnectionSetBasicOperations(t *testing.T) {
	src1, _ := netset.IPBlockFromCidr("10.240.10.0/24")
	dst1, _ := netset.IPBlockFromCidr("10.240.10.0/32")
	dst2 := src1.Subtract(dst1)
	conn1 := netset.ConnectionSetFrom(src1, dst1, connection.NewTCPSet())
	conn2 := netset.ConnectionSetFrom(src1, dst2, connection.NewTCPSet())
	conn3 := netset.ConnectionSetFrom(src1, src1, connection.NewTCPSet())

	// basic union & Equal test
	unionConn := conn1.Union(conn2)
	require.True(t, unionConn.Equal(conn3))
	require.True(t, conn3.Equal(unionConn))

	// basic subtract & Equal test
	conn4 := netset.ConnectionSetFrom(src1, dst2, connection.All())
	subttractionRes := conn3.Subtract(conn4)
	require.True(t, subttractionRes.Equal(conn1))
	require.True(t, conn1.Equal(subttractionRes))

	// basic IsSubset test
	require.True(t, conn1.IsSubset(conn3))
	require.True(t, conn2.IsSubset(conn3))
	require.False(t, conn2.IsSubset(conn1))
	require.False(t, conn1.IsSubset(conn2))

	// basic IsEmpty test
	require.False(t, conn1.IsEmpty())
	require.True(t, netset.NewConnectionSet().IsEmpty())
}
