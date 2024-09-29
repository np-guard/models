/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//nolint:lll //long lines for tests used to document the connection objects
package netset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
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
	subtractionRes := conn3.Subtract(conn4) // removes all connections over (src1, dst2) from conn3
	require.True(t, subtractionRes.Equal(conn1))
	require.True(t, conn1.Equal(subtractionRes))

	// basic IsSubset test
	require.True(t, conn1.IsSubset(conn3))
	require.True(t, conn2.IsSubset(conn3))
	require.False(t, conn2.IsSubset(conn1))
	require.False(t, conn1.IsSubset(conn2))

	// basic IsEmpty test
	require.False(t, conn1.IsEmpty())
	require.True(t, netset.NewConnectionSet().IsEmpty())

	// demonstrate split in allowed connections for dest dimension, to be reflected in partitions
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
	// conns: TCP
	fmt.Printf("conn1 cubes string:\n%s\n", conn1.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27,10.240.10.64/26, 10.240.10.128/25,
	// conns: TCP
	fmt.Printf("conn2 cubes string:\n%s\n", conn2.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.0/24,
	// conns: TCP
	fmt.Printf("conn3 cubes string:\n%s\n", conn3.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: All Connections
	fmt.Printf("conn4 cubes string:\n%s\n", conn4.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: ICMP
	fmt.Printf("conn5 cubes string:\n%s\n", conn5.String())

	// two partitions in the following object:
	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: ICMP,TCP
	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: TCP
	fmt.Printf("conn5UnionConn2 cubes string:\n%s\n", conn5UnionConn2.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.2/31,
	// conns: ICMP,TCP
	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: TCP
	fmt.Printf("conn8 cubes string:\n%s\n", conn8.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.2/31, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns:  All Connections
	fmt.Printf("conn9UnionConn6 cubes string:\n%s\n", conn9UnionConn6.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: TCP
	// src: 10.240.10.0/25,
	// dst: 10.240.10.2/31,
	// conns: ICMP;TCP,UDP dst-ports: 53
	// src: 10.240.10.128/25,
	// dst: 10.240.10.2/31,
	// conns: ICMP,TCP
	fmt.Printf("conn12 cubes string:\n%s\n", conn12.String())

	// src: 10.240.10.0/24,
	// dst: 10.240.10.1/32, 10.240.10.4/30, 10.240.10.8/29, 10.240.10.16/28, 10.240.10.32/27, 10.240.10.64/26, 10.240.10.128/25,
	// conns: TCP,
	// src: 10.240.10.0/25,
	// dst: 10.240.10.2/31,
	// conns: ICMP;TCP,UDP dst-ports: 53,
	// src: 10.240.10.128/25,
	// dst: 10.240.10.2/31,
	// conns: ICMP,TCP
	fmt.Printf("conn17 cubes string:\n%s\n", conn17.String())

	fmt.Println("done")
}
