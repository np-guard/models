// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/connection"
	"github.com/np-guard/models/pkg/netp"
)

const ICMPValue = netp.DestinationUnreachable

func TestAllConnections(t *testing.T) {
	c := connection.All()
	require.Equal(t, "All Connections", c.String())
}

func TestNoConnections(t *testing.T) {
	c := connection.None()
	require.Equal(t, "No Connections", c.String())
}

func TestBasicSetICMP(t *testing.T) {
	c := connection.ICMPConnection(ICMPValue, ICMPValue, 5, 5)
	require.Equal(t, "protocol: ICMP icmp-type: 3 icmp-code: 5", c.String())
	require.Equal(t, "ICMP type: 3 code: 5", c.MinimizeString())
}

func TestBasicSetTCP(t *testing.T) {
	e := connection.TCPorUDPConnection(netp.ProtocolStringTCP, 1, 65535, 1, 655)
	require.Equal(t, "protocol: TCP dst-ports: 1-655", e.String())
	require.Equal(t, "TCP 1-655", e.MinimizeString())

	e = connection.TCPorUDPConnection(netp.ProtocolStringTCP, 1, 535, 1, 655)
	require.Equal(t, "protocol: TCP src-ports: 1-535 dst-ports: 1-655", e.String())
	require.Equal(t, "TCP src-ports: 1-535 dst-ports: 1-655", e.MinimizeString())

	e = connection.TCPorUDPConnection(netp.ProtocolStringTCP, 1, 65535, 1, 65535)
	require.Equal(t, "protocol: TCP", e.String())

	c := connection.All().Subtract(e)
	require.Equal(t, "protocol: ICMP,UDP", c.String())
	require.Equal(t, "ICMP,UDP", c.MinimizeString())

	c = c.Union(e)
	require.Equal(t, "All Connections", c.String())
}

func TestBasicSet2(t *testing.T) {
	except1 := connection.ICMPConnection(ICMPValue, ICMPValue, 5, 5)

	except2 := connection.TCPorUDPConnection(netp.ProtocolStringTCP, 1, 65535, 1, 65535)

	d := connection.All().Subtract(except1).Subtract(except2)
	require.Equal(t, ""+
		"protocol: ICMP icmp-type: 0-2,4-16; "+
		"protocol: ICMP icmp-type: 3 icmp-code: 0-4; "+
		"protocol: UDP", d.String())
	require.Equal(t, ""+
		"ICMP type: 0-2,4-16; "+
		"ICMP type: 3 code: 0-4; "+
		"UDP", d.MinimizeString())
}

func TestBasicSet3(t *testing.T) {
	c := connection.ICMPConnection(ICMPValue, ICMPValue, 5, 5)
	d := connection.All().Subtract(c).Union(connection.ICMPConnection(ICMPValue, ICMPValue, 5, 5))
	require.Equal(t, "All Connections", d.String())
}
