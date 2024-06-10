/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
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
// Returns a connection object with the exact subset of the stateful part (within TCP)
// from the entire connection `c` and with the original connections on other protocols.
func (c *Set) WithStatefulness(secondDirectionConn *Set) *Set {
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
