// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection

import (
	"slices"

	"github.com/np-guard/models/pkg/netp"
)

func NewTCPSet() *Set {
	return TCPorUDPConnection(netp.ProtocolStringTCP, MinPort, MaxPort, MinPort, MaxPort)
}

// PartitionTCPNonTCP given a connection returns its TCP and non-TCP sub-connections
func PartitionTCPNonTCP(conn *Set) (tcp, nonTCP *Set) {
	tcpFractionOfConn := NewTCPSet().Intersect(conn)
	nonTCPFractionOfConn := conn.Subtract(tcpFractionOfConn)
	return tcpFractionOfConn, nonTCPFractionOfConn
}

// GetResponsiveConn returns  connection object with the exact the responsive part within TCP
// and with the original connections on other protocols.
// `c` represents a src-to-dst connection, and `secondDirectionConn` represents dst-to-src connection.
// todo: move to analyzer
func (c *Set) GetResponsiveConn(secondDirectionConn *Set) *Set {
	connTCP := c.Intersect(NewTCPSet())
	if connTCP.IsEmpty() {
		return c
	}
	tcpSecondDirection := secondDirectionConn.Intersect(NewTCPSet())
	// flip src/dst ports before intersection
	tcpSecondDirectionFlipped := tcpSecondDirection.SwitchSrcDstPorts()
	// tcp connection stateful subset
	statefulCombinedConnTCP := connTCP.Intersect(tcpSecondDirectionFlipped)
	return c.Subtract(connTCP).Union(statefulCombinedConnTCP)
}

// SwitchSrcDstPorts returns a new Set object, built from the input Set object.
// The src and dst ports on relevant cubes are being switched.
func (c *Set) SwitchSrcDstPorts() *Set {
	if c.IsAll() {
		return c.Copy()
	}
	newConn := c.connectionProperties.SwapDimensions(slices.Index(dimensionsList, srcPort), slices.Index(dimensionsList, dstPort))
	return &Set{
		connectionProperties: newConn,
	}
}
