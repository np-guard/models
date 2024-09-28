/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//nolint:lll //long lines for tests used to document the connection objects
package netset_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

// TestConnectionSetBasicOperations tests basic operations on netset.ConnectionSet objects
func TestConnectionSetBasicOperations(t *testing.T) {
	// relevant src/dst IP objects
	cidr1, _ := netset.IPBlockFromCidr("10.240.10.0/24")
	cidr2, _ := netset.IPBlockFromCidr("10.240.10.0/32")
	cidr1MinusCidr2 := cidr1.Subtract(cidr2) // 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25
	subsetOfCidr1MinusCidr2, _ := netset.IPBlockFromCidr("10.240.10.2/31")
	//  10.240.10.0/25 union 10.240.10.128/25 == 10.240.10.0/24
	leftHalfCidr1, _ := netset.IPBlockFromCidr("10.240.10.0/25")
	rightHalfCidr1, _ := netset.IPBlockFromCidr("10.240.10.128/25")

	// relevant connection set objects
	conn1 := netset.ConnectionSetFrom(cidr1, cidr2, connection.NewTCPSet())           // conns from cidr1 to cidr2 over all TCP
	conn2 := netset.ConnectionSetFrom(cidr1, cidr1MinusCidr2, connection.NewTCPSet()) // conns from cidr1 to cidr1MinusCidr2 over all TCP
	conn3 := netset.ConnectionSetFrom(cidr1, cidr1, connection.NewTCPSet())           // conns from cidr1 to cidr1 over all TCP

	// basic union & Equal test
	unionConn := conn1.Union(conn2)
	require.True(t, unionConn.Equal(conn3)) // union of dest dimension (src and conn dimensions are common)
	require.True(t, conn3.Equal(unionConn))

	// basic subtract & Equal test
	conn4 := netset.ConnectionSetFrom(cidr1, cidr1MinusCidr2, connection.All())
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

	// demonstrate split in allwoed connections for dest dimension, to be reflected in partitions
	conn5 := netset.ConnectionSetFrom(cidr1, subsetOfCidr1MinusCidr2, connection.AllICMP())
	conn5UnionConn2 := conn5.Union(conn2) // for certain dest- icmp+tcp, and for remaining dest- only tcp [common src for both]
	require.Equal(t, 2, len(conn5UnionConn2.Partitions()))

	// other operations on other objects, to get equiv object of conn5UnionConn2:
	tcpAndICMP := connection.NewTCPSet().Union(connection.AllICMP())
	conn6 := netset.ConnectionSetFrom(cidr1, cidr1MinusCidr2, tcpAndICMP)
	deltaCidrs := cidr1MinusCidr2.Subtract(subsetOfCidr1MinusCidr2)
	conn7 := netset.ConnectionSetFrom(cidr1, deltaCidrs, connection.AllICMP())
	conn8 := conn6.Subtract(conn7)
	require.True(t, conn8.Equal(conn5UnionConn2))

	// add udp to tcpAndICMP => check it is All()
	conn9 := netset.ConnectionSetFrom(cidr1, cidr1MinusCidr2, connection.NewUDPSet())
	conn10 := netset.ConnectionSetFrom(cidr1, cidr1MinusCidr2, connection.All())
	conn9UnionConn6 := conn9.Union(conn6)
	require.True(t, conn10.Equal(conn9UnionConn6))

	// demonstrate split in allowed connections for src dimensions, to be reflected in partitions
	// starting from conn8
	udp53 := connection.NewUDP(netp.MinPort, netp.MaxPort, 53, 53)
	conn11 := netset.ConnectionSetFrom(leftHalfCidr1, subsetOfCidr1MinusCidr2, udp53)
	conn12 := conn11.Union(conn8)

	// another way to produce obj equiv to conn12 :
	conn13 := netset.ConnectionSetFrom(leftHalfCidr1, subsetOfCidr1MinusCidr2, tcpAndICMP.Union(udp53))
	conn14 := netset.ConnectionSetFrom(leftHalfCidr1, cidr1MinusCidr2, connection.NewTCPSet())
	conn15 := netset.ConnectionSetFrom(rightHalfCidr1, subsetOfCidr1MinusCidr2, tcpAndICMP)
	conn16 := netset.ConnectionSetFrom(rightHalfCidr1, cidr1MinusCidr2, connection.NewTCPSet())
	conn17 := (conn13.Union(conn14)).Union(conn15.Union(conn16))
	require.True(t, conn12.Equal(conn17))

	// partitions string examples - for the objects used in this test

	// src: 10.240.10.0/24,
	// dst: 10.240.10.0,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	fmt.Printf("conn1 cubes string:\n%s\n", getPartitionsStr(conn1))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27,10.240.10.64/26, 10.240.10.128/25,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	fmt.Printf("conn2 cubes string:\n%s\n", getPartitionsStr(conn2))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.0/24,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	fmt.Printf("conn3 cubes string:\n%s\n", getPartitionsStr(conn3))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: all
	fmt.Printf("conn4 cubes string:\n%s\n", getPartitionsStr(conn4))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: ;all icmp
	fmt.Printf("conn5 cubes string:\n%s\n", getPartitionsStr(conn5))

	// two partitions in the following object:
	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535; all icmp
	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	fmt.Printf("conn5UnionConn2 cubes string:\n%s\n", getPartitionsStr(conn5UnionConn2))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;all icmp
	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	fmt.Printf("conn8 cubes string:\n%s\n", getPartitionsStr(conn8))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: all
	fmt.Printf("conn9UnionConn6 cubes string:\n%s\n", getPartitionsStr(conn9UnionConn6))

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;
	// src: 10.240.10.0/25,
	// dst: 10.240.10.2/31,
	// conns: protocols 1, src-ports 1-65535, dst-ports 53,protocols 0, src-ports 1-65535, dst-ports 1-65535;all icmp
	// src: 10.240.10.128/25,
	// dst: 10.240.10.2/31,
	// conns: protocols 0, src-ports 1-65535, dst-ports 1-65535;all icmp
	fmt.Printf("conn12 cubes string:\n%s\n", getPartitionsStr(conn12))

	fmt.Printf("conn17 cubes string:\n%s\n", getPartitionsStr(conn17))

	fmt.Println("done")
}

// simple string functions for testing

func icmpStr(icmpObj netp.ICMP) string {
	if icmpObj.TypeCode.Code != nil {
		return fmt.Sprintf("icmp type: %d, code: %d", icmpObj.TypeCode.Type, *icmpObj.TypeCode.Code)
	}
	return fmt.Sprintf("icmp type: %d", icmpObj.TypeCode.Type)
}

func icmpSetStr(icmpset *netset.ICMPSet) string {
	if icmpset.IsAll() {
		return "all icmp"
	}
	if icmpset.IsEmpty() {
		return ""
	}
	cubes := icmpset.Partitions()
	cubesStrings := make([]string, len(cubes))
	for i := range cubes {
		cubesStrings[i] = icmpStr(cubes[i])
	}
	return strings.Join(cubesStrings, ",")
}

func tcpudpSetStr(tcpupdset *netset.TCPUDPSet) string {
	if tcpupdset.IsAll() {
		return "all tcp,udp"
	}
	if tcpupdset.IsEmpty() {
		return ""
	}
	cubes := tcpupdset.Partitions()
	cubesStrings := make([]string, len(cubes))
	for i := range cubes {
		cube := cubes[i] // each cube is of type ds.Triple[*ProtocolSet, *PortSet, *PortSet]
		cubesStrings[i] = fmt.Sprintf("protocols %s, src-ports %s, dst-ports %s", cube.S1.String(), cube.S2.String(), cube.S3.String())
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

func getPartitionsStr(conn *netset.ConnectionSet) string {
	cubes := conn.Partitions()
	cubesStrings := make([]string, len(cubes))
	for i := range cubes {
		cubesStrings[i] = cubeStr(cubes[i])
	}
	return strings.Join(cubesStrings, "\n")
}
