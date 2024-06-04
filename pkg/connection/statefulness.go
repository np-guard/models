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

func PartitionTCPNonTCP(conn *Set) (tcp, nonTCP *Set) {
	tcpFractionOfConn := NewTCPSet().Intersect(conn)
	nonTCPFractionOfConn := conn.Subtract(tcpFractionOfConn)
	return tcpFractionOfConn, nonTCPFractionOfConn
}

// WithStatefulness returns the stateful part of `c`
// `c` represents a src-to-dst connection, and `secondDirectionConn` represents dst-to-src connection.
// This function also returns a connection object with the exact subset of the stateful part (within TCP)
// from the entire connection `c`, and with the original connections on other protocols.
func (c *Set) WithStatefulness(secondDirectionConn *Set) *Set {
	connTCP := c.Intersect(NewTCPSet())
	if connTCP.IsEmpty() {
		return c
	}
	statefulCombinedConnTCP := connTCP.connTCPStatefulness(secondDirectionConn.Intersect(NewTCPSet()))
	return c.Subtract(connTCP).Union(statefulCombinedConnTCP)
}

// connTCPWithStatefulness assumes that both `c` and `secondDirectionConn` are within TCP.
// it returns the subset from `c` which is stateful.
func (c *Set) connTCPStatefulness(secondDirectionConn *Set) *Set {
	// flip src/dst ports before intersection
	statefulCombinedConn := c.Intersect(secondDirectionConn.switchSrcDstPortsOnTCP())
	return statefulCombinedConn
}

// switchSrcDstPortsOnTCP returns a new Set object, built from the input Set object.
// It assumes the input connection object is only within TCP protocol.
// For TCP the src and dst ports on relevant cubes are being switched.
func (c *Set) switchSrcDstPortsOnTCP() *Set {
	if c.IsAll() {
		return c.Copy()
	}
	newConn := c.connectionProperties.SwapDimensions(slices.Index(dimensionsList, srcPort), slices.Index(dimensionsList, dstPort))
	return &Set{
		connectionProperties: newConn,
	}
}
