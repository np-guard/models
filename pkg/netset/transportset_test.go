/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

const ICMPValue = netp.DestinationUnreachable

func TestAllConnectionsTransportSet(t *testing.T) {
	c := netset.AllTransportSet()
	// String
	require.Equal(t, netset.AllConnections, c.String())
	// IsAll
	require.True(t, c.IsAll())

	// Partitions
	icmpPartitions := c.ICMPSet().Partitions()
	tcpudpPartitions := c.TCPUDPSet().Partitions()
	require.Equal(t, 1, len(tcpudpPartitions))
	require.Equal(t, 1, len(icmpPartitions))
	// all tcp-udp
	require.True(t, tcpudpPartitions[0].S1.Equal(netset.AllTCPUDPProtocolSet()))
	require.True(t, tcpudpPartitions[0].S2.Equal(netset.AllPorts()))
	require.True(t, tcpudpPartitions[0].S3.Equal(netset.AllPorts()))
	// all icmp
	require.Nil(t, icmpPartitions[0].TypeCode)
}

func TestNoConnectionsTransportSet(t *testing.T) {
	c := netset.AllOrNothingTransport(false, false)
	require.Equal(t, netset.NoConnections, c.String())

	require.True(t, c.IsEmpty())
	icmpPartitions := c.ICMPSet().Partitions()
	tcpudpPartitions := c.TCPUDPSet().Partitions()
	require.Equal(t, 0, len(tcpudpPartitions))
	require.Equal(t, 0, len(icmpPartitions))
}

func TestBasicSetICMPTransportSet(t *testing.T) {
	c, err := netset.NewICMPTransportFromTypeCode(ICMPValue, 5)
	require.Nil(t, err)
	fmt.Println(c) // ICMP icmp-type: 3 icmp-code: 5
	require.Equal(t, "ICMP icmp-type: 3 icmp-code: 5", c.String())
}

func TestBasicSetTCPTransportSet(t *testing.T) {
	e := netset.NewTCPorUDPTransport(netp.ProtocolStringTCP, 1, 65535, 1, 655)
	fmt.Println(e) // TCP dst-ports: 1-655
	require.Equal(t, "TCP dst-ports: 1-655", e.String())

	e = netset.NewTCPorUDPTransport(netp.ProtocolStringTCP, 1, 535, 1, 655)
	fmt.Println(e) // TCP src-ports: 1-535 dst-ports: 1-655
	require.Equal(t, "TCP src-ports: 1-535 dst-ports: 1-655", e.String())

	e = netset.NewTCPorUDPTransport(netp.ProtocolStringTCP, 1, 65535, 1, 65535)
	fmt.Println(e)
	require.Equal(t, "TCP", e.String())

	c := netset.AllTransportSet().Subtract(e)
	fmt.Println(c)
	require.Equal(t, "ICMP,UDP", c.String())

	c = c.Union(e)
	require.Equal(t, netset.AllConnections, c.String())
}

func TestBasicSet2TransportSet(t *testing.T) {
	except1, err := netset.NewICMPTransportFromTypeCode(ICMPValue, 5)
	require.Nil(t, err)
	except2 := netset.NewTCPorUDPTransport(netp.ProtocolStringTCP, 1, 65535, 1, 65535)

	d := netset.AllTransportSet().Subtract(except1).Subtract(except2)
	fmt.Println(d)
	// ICMP icmp-type: 0 icmp-code: 0;icmp-type: 11;icmp-type: 12 icmp-code: 0;icmp-type: 13 icmp-code: 0;
	// icmp-type: 14 icmp-code: 0;icmp-type: 15 icmp-code: 0;icmp-type: 16 icmp-code: 0;icmp-type: 3 icmp-code: 0;
	// icmp-type: 3 icmp-code: 1;icmp-type: 3 icmp-code: 2;icmp-type: 3 icmp-code: 3;icmp-type: 3 icmp-code: 4;
	// icmp-type: 4 icmp-code: 0;icmp-type: 5;icmp-type: 8 icmp-code: 0;UDP

	require.Equal(t, 15, len(d.ICMPSet().Partitions()))
	require.Equal(t, 1, len(d.TCPUDPSet().Partitions()))

	fmt.Println("done")
	/* string from older version:
	"protocol: ICMP icmp-type: 0-2,4-16; "+
	"protocol: ICMP icmp-type: 3 icmp-code: 0-4; "+
	"protocol: UDP", d.String())
	*/
}

func TestBasicSet3TransportSet(t *testing.T) {
	c, err := netset.NewICMPTransportFromTypeCode(ICMPValue, 5)
	c1, _ := netset.NewICMPTransportFromTypeCode(ICMPValue, 5)
	require.Nil(t, err)
	d := netset.AllTransportSet().Subtract(c).Union(c1)
	require.Equal(t, netset.AllConnections, d.String())
}
