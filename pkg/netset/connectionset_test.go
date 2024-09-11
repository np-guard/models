/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netset"
)

func icmpSetStr(icmpset *netset.ICMPSet) string {
	if icmpset.IsAll() {
		return "all icmp"
	}
	if icmpset.IsEmpty() {
		return ""
	}
	return ""
}

func tcpudpSetStr(tcpupdset *netset.TCPUDPSet) string {
	if tcpupdset.IsAll() {
		return "all tcp,udp"
	}
	if tcpupdset.IsEmpty() {
		return ""
	}
	cubes := tcpupdset.Partitions()
	cubesStrings := []string{}
	for i := range cubes {
		cube := cubes[i] // ds.Triple[*ProtocolSet, *PortSet, *PortSet]

	}
	return strings.Join(cubesStrings, ",")
}

func transportSetStr(conn *netset.TransportSet) string {
	if conn.IsAll() {
		return "all"
	}
	if conn.IsEmpty() {
		return "empty"
	}
	tcpudpSet := conn.TCPUDPSet()
	icmpSet := conn.ICMPSet()
	resStrList := []string{tcpudpSetStr(tcpudpSet), icmpSetStr(icmpSet)}

	return strings.Join(resStrList, ";")
}

func cubeStr(c ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]) string {
	return fmt.Sprintf("src: %s, dst: %s, conns: %s", c.S1.String(), c.S2.String(), transportSetStr(c.S3))
}

func TestConnectionSetBasicOperations(t *testing.T) {
	// relevant src/dst IP objects
	src1, _ := netset.IPBlockFromCidr("10.240.10.0/24")
	dst1, _ := netset.IPBlockFromCidr("10.240.10.0/32")
	dst2 := src1.Subtract(dst1) // 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25

	// relevant connection set objects
	conn1 := netset.ConnectionSetFrom(src1, dst1, connection.NewTCPSet()) // conns from src1 to dst1 over all TCP
	conn2 := netset.ConnectionSetFrom(src1, dst2, connection.NewTCPSet()) // conns from src1 to dst2 over all TCP
	conn3 := netset.ConnectionSetFrom(src1, src1, connection.NewTCPSet()) // conns from src1 to src1 over all TCP

	// basic union & Equal test
	unionConn := conn1.Union(conn2)
	require.True(t, unionConn.Equal(conn3)) // union of dest dimension (src and conn dimensions are common)
	require.True(t, conn3.Equal(unionConn))

	// basic subtract & Equal test
	conn4 := netset.ConnectionSetFrom(src1, dst2, connection.All())
	subttractionRes := conn3.Subtract(conn4) // removes all connections over (src1, dst2) from conn3
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

	// partitions test
	cubes := conn1.Partitions()

}
